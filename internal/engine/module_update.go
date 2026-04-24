package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const defaultReleaseAPI = "https://api.github.com/repos/hari7261/INDUS/releases/latest"

type updateModule struct {
	engine *Engine
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type latestRelease struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	PublishedAt string         `json:"published_at"`
	Assets      []releaseAsset `json:"assets"`
}

func (m *updateModule) Execute(ctx context.Context, inv Invocation) Response {
	action := "check"
	if len(inv.Parsed.Positionals) > 0 {
		switch strings.ToLower(inv.Parsed.Positionals[0]) {
		case "check", "download", "install":
			action = strings.ToLower(inv.Parsed.Positionals[0])
		}
	}
	if inv.Parsed.Bool("download") {
		action = "download"
	}
	if inv.Parsed.Bool("install") {
		action = "install"
	}

	release, err := m.fetchLatestRelease(ctx)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	current := normalizeVersion(m.engine.Version())
	latest := normalizeVersion(release.TagName)
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "current=%s\nlatest=%s\n", current, latest)
	if release.PublishedAt != "" {
		fmt.Fprintf(buffer, "published_at=%s\n", release.PublishedAt)
	}

	if compareVersions(latest, current) <= 0 && !inv.Parsed.Bool("force") {
		fmt.Fprint(buffer, "status=up_to_date")
		return Response{Output: strings.TrimSpace(buffer.String())}
	}

	asset, err := selectReleaseAsset(release.Assets)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	fmt.Fprintf(buffer, "asset=%s\n", asset.Name)

	switch action {
	case "check":
		fmt.Fprint(buffer, "status=update_available")
		return Response{Output: strings.TrimSpace(buffer.String())}
	case "download":
		target := inv.Parsed.String("output")
		if target == "" {
			target = filepath.Join(m.engine.paths.UpdatesDir, asset.Name)
		}
		target, err = m.downloadAsset(ctx, asset, target)
		if err != nil {
			return Response{Err: commandFailedError(inv.Command, err)}
		}
		fmt.Fprintf(buffer, "status=downloaded\npath=%s", target)
		return Response{Output: strings.TrimSpace(buffer.String())}
	case "install":
		target := filepath.Join(m.engine.paths.UpdatesDir, asset.Name)
		target, err = m.downloadAsset(ctx, asset, target)
		if err != nil {
			return Response{Err: commandFailedError(inv.Command, err)}
		}
		exePath, err := os.Executable()
		if err != nil {
			return Response{Err: commandFailedError(inv.Command, err)}
		}
		if err := stageUpdaterScript(target, exePath); err != nil {
			return Response{Err: commandFailedError(inv.Command, err)}
		}
		fmt.Fprintf(buffer, "status=install_scheduled\npath=%s\nrestart_required=true", target)
		return Response{
			Output: strings.TrimSpace(buffer.String()),
			Effects: Effects{
				Exit: true,
			},
		}
	default:
		return Response{Err: invalidArgumentError(inv.Command, "unsupported update action")}
	}
}

func (m *updateModule) fetchLatestRelease(ctx context.Context) (latestRelease, error) {
	url := os.Getenv("INDUS_UPDATE_API")
	if url == "" {
		url = defaultReleaseAPI
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return latestRelease{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "INDUS-Terminal/"+m.engine.Version())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return latestRelease{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return latestRelease{}, fmt.Errorf("release lookup failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var release latestRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return latestRelease{}, err
	}
	return release, nil
}

func (m *updateModule) downloadAsset(ctx context.Context, asset releaseAsset, target string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.BrowserDownloadURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "INDUS-Terminal/"+m.engine.Version())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("asset download failed: %s", resp.Status)
	}

	file, err := os.Create(target)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", err
	}
	return target, nil
}

func selectReleaseAsset(assets []releaseAsset) (releaseAsset, error) {
	if runtime.GOOS != "windows" {
		return releaseAsset{}, fmt.Errorf("self-update is currently supported on Windows builds only")
	}

	want := fmt.Sprintf("windows-%s.exe", runtime.GOARCH)
	for _, asset := range assets {
		if strings.Contains(strings.ToLower(asset.Name), want) {
			return asset, nil
		}
	}
	for _, asset := range assets {
		if strings.HasSuffix(strings.ToLower(asset.Name), ".exe") {
			return asset, nil
		}
	}
	return releaseAsset{}, fmt.Errorf("no compatible release asset found")
}

func stageUpdaterScript(source, target string) error {
	scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("indus-update-%d.ps1", time.Now().UnixNano()))
	script := fmt.Sprintf(`$source = %q
$target = %q
$targetDir = Split-Path -Parent $target
Start-Sleep -Seconds 2
$updated = $false
for ($i = 0; $i -lt 12; $i++) {
  try {
    Copy-Item -LiteralPath $source -Destination $target -Force
    $updated = $true
    break
  } catch {
    Start-Sleep -Milliseconds 500
  }
}
if (-not $updated) { exit 1 }
Start-Process -FilePath $target -WorkingDirectory $targetDir -WindowStyle Normal
`, source, target)

	if err := os.WriteFile(scriptPath, []byte(script), 0o600); err != nil {
		return err
	}

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-WindowStyle", "Hidden", "-File", scriptPath)
	return cmd.Start()
}

func normalizeVersion(value string) string {
	value = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(value), "v"))
	if idx := strings.IndexAny(value, "-+"); idx >= 0 {
		value = value[:idx]
	}
	return value
}

func compareVersions(left, right string) int {
	leftParts := versionParts(left)
	rightParts := versionParts(right)
	size := len(leftParts)
	if len(rightParts) > size {
		size = len(rightParts)
	}
	for i := 0; i < size; i++ {
		lv := 0
		rv := 0
		if i < len(leftParts) {
			lv = leftParts[i]
		}
		if i < len(rightParts) {
			rv = rightParts[i]
		}
		switch {
		case lv > rv:
			return 1
		case lv < rv:
			return -1
		}
	}
	return 0
}

func versionParts(value string) []int {
	if value == "" {
		return nil
	}
	raw := strings.Split(value, ".")
	parts := make([]int, 0, len(raw))
	for _, item := range raw {
		n, err := strconv.Atoi(item)
		if err != nil {
			parts = append(parts, 0)
			continue
		}
		parts = append(parts, n)
	}
	return parts
}
