package engine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestEngine(t testing.TB) *Engine {
	t.Helper()

	root := t.TempDir()
	return newTestEngineAt(t, root)
}

func newTestEngineAt(t testing.TB, root string) *Engine {
	t.Helper()

	t.Setenv("APPDATA", root)
	t.Setenv("LOCALAPPDATA", root)
	t.Setenv("HOME", root)
	t.Setenv("USERPROFILE", root)
	t.Setenv("INDUS_CONFIG", filepath.Join(root, "config.cfg"))
	t.Setenv("INDUS_STATE_DIR", filepath.Join(root, "state"))
	t.Setenv("INDUS_CACHE_DIR", filepath.Join(root, "cache"))

	engine, err := New(Options{
		Version:   "1.5.5-test",
		Commit:    "test",
		BuildTime: "2026-04-25T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	return engine
}

func TestRegistryCommandCount(t *testing.T) {
	engine := newTestEngine(t)
	if got := len(engine.registry); got < 55 {
		t.Fatalf("expected at least 55 commands, got %d", got)
	}
	if engine.RegistryVersion() != "1.5.5" {
		t.Fatalf("unexpected registry version: %s", engine.RegistryVersion())
	}
}

func TestParseCommandLinePreservesQuotedArguments(t *testing.T) {
	args := ParseCommandLine(`ind env set INDUS_PROFILE "Production Ready"`)
	if len(args) != 5 {
		t.Fatalf("unexpected arg count: %v", args)
	}
	if args[4] != "Production Ready" {
		t.Fatalf("expected quoted argument to be preserved, got %q", args[4])
	}
}

func TestDeprecatedAliasMapsToCanonical(t *testing.T) {
	engine := newTestEngine(t)
	session := engine.NewSession(t.TempDir())

	response := engine.ExecuteTokens(context.Background(), session, []string{"scan"}, ModeInteractive)
	if response.Err != nil {
		t.Fatalf("expected success, got %v", response.Err)
	}
	if !strings.Contains(response.Warning, "ind scan") {
		t.Fatalf("expected deprecation warning, got %q", response.Warning)
	}
}

func TestInvalidCommandReturnsStructuredError(t *testing.T) {
	engine := newTestEngine(t)
	session := engine.NewSession(t.TempDir())

	response := engine.ExecuteTokens(context.Background(), session, []string{"unknown"}, ModeExecutable)
	if response.Err == nil {
		t.Fatal("expected error")
	}
	if response.Err.Code != "IND_ERR_002" {
		t.Fatalf("expected IND_ERR_002, got %s", response.Err.Code)
	}
}

func TestEnvSetAndList(t *testing.T) {
	engine := newTestEngine(t)
	session := engine.NewSession(t.TempDir())

	setResponse := engine.ExecuteTokens(context.Background(), session, []string{"env", "set", "INDUS_PROFILE", "Production"}, ModeExecutable)
	if setResponse.Err != nil {
		t.Fatalf("env set failed: %v", setResponse.Err)
	}

	listResponse := engine.ExecuteTokens(context.Background(), session, []string{"env", "list"}, ModeExecutable)
	if listResponse.Err != nil {
		t.Fatalf("env list failed: %v", listResponse.Err)
	}
	if !strings.Contains(listResponse.Output, "INDUS_PROFILE=Production") {
		t.Fatalf("expected managed env output, got %q", listResponse.Output)
	}
}

func TestProjectCreateAndBuild(t *testing.T) {
	engine := newTestEngine(t)
	workspace := t.TempDir()
	session := engine.NewSession(workspace)

	createResponse := engine.ExecuteTokens(context.Background(), session, []string{"proj", "create", "orbit-app", "--dir", workspace}, ModeExecutable)
	if createResponse.Err != nil {
		t.Fatalf("proj create failed: %v", createResponse.Err)
	}

	projectRoot := filepath.Join(workspace, "orbit-app")
	buildSession := engine.NewSession(projectRoot)
	buildResponse := engine.ExecuteTokens(context.Background(), buildSession, []string{"proj", "build", projectRoot}, ModeExecutable)
	if buildResponse.Err != nil {
		t.Fatalf("proj build failed: %v", buildResponse.Err)
	}

	if _, err := os.Stat(filepath.Join(projectRoot, "build", "indus-artifact.json")); err != nil {
		t.Fatalf("artifact missing: %v", err)
	}
}

func TestSysStatsUsesCache(t *testing.T) {
	engine := newTestEngine(t)
	session := engine.NewSession(t.TempDir())

	first := engine.ExecuteTokens(context.Background(), session, []string{"sys", "stats"}, ModeExecutable)
	if first.Err != nil {
		t.Fatalf("first sys stats failed: %v", first.Err)
	}

	second := engine.ExecuteTokens(context.Background(), session, []string{"sys", "stats"}, ModeExecutable)
	if second.Err != nil {
		t.Fatalf("second sys stats failed: %v", second.Err)
	}
	if !second.Cached {
		t.Fatal("expected second sys stats call to use cache")
	}
}

func TestDevBenchAcceptsQuotedCommandFlagValue(t *testing.T) {
	engine := newTestEngine(t)
	session := engine.NewSession(t.TempDir())

	response := engine.ExecuteTokens(context.Background(), session, []string{
		"dev", "bench", "--command", "ind sys stats", "--runs", "2",
	}, ModeExecutable)
	if response.Err != nil {
		t.Fatalf("dev bench failed: %v", response.Err)
	}
	if !strings.Contains(response.Output, "command=ind sys stats") {
		t.Fatalf("unexpected output: %q", response.Output)
	}
}

func TestDevBenchAcceptsSplitCommandFlagValue(t *testing.T) {
	engine := newTestEngine(t)
	session := engine.NewSession(t.TempDir())

	response := engine.ExecuteTokens(context.Background(), session, []string{
		"dev", "bench", "--command", "ind", "sys", "stats", "--runs", "2",
	}, ModeExecutable)
	if response.Err != nil {
		t.Fatalf("dev bench failed: %v", response.Err)
	}
	if !strings.Contains(response.Output, "command=ind sys stats") {
		t.Fatalf("unexpected output: %q", response.Output)
	}
}

func TestLegacyHTTPAliasMapsToNetFetch(t *testing.T) {
	engine := newTestEngine(t)

	tokens, warning := engine.normalizeTokens([]string{"ind", "http", "get", "https://api.github.com"}, ModeExecutable)
	if warning != "" {
		t.Fatalf("unexpected warning: %q", warning)
	}

	want := []string{"net", "fetch", "https://api.github.com", "--method", "GET"}
	if strings.Join(tokens, "|") != strings.Join(want, "|") {
		t.Fatalf("unexpected tokens: got %v want %v", tokens, want)
	}
}

func TestLegacyNetHTTPAliasMapsDataToBody(t *testing.T) {
	engine := newTestEngine(t)

	tokens, warning := engine.normalizeTokens([]string{"ind", "net", "http", "post", "https://api.example.com", "--data", "{\"ok\":true}"}, ModeExecutable)
	if warning != "" {
		t.Fatalf("unexpected warning: %q", warning)
	}

	want := []string{"net", "fetch", "https://api.example.com", "--method", "POST", "--body", "{\"ok\":true}"}
	if strings.Join(tokens, "|") != strings.Join(want, "|") {
		t.Fatalf("unexpected tokens: got %v want %v", tokens, want)
	}
}

func TestTermProfilePersistsPromptAndBannerSettings(t *testing.T) {
	root := t.TempDir()
	engine := newTestEngineAt(t, root)
	session := engine.NewSession(t.TempDir())

	setPrompt := engine.ExecuteTokens(context.Background(), session, []string{"term", "profile", "set", "prompt", "INDUS-PRO"}, ModeExecutable)
	if setPrompt.Err != nil {
		t.Fatalf("profile set prompt failed: %v", setPrompt.Err)
	}

	setBanner := engine.ExecuteTokens(context.Background(), session, []string{"term", "profile", "set", "banner", "off"}, ModeExecutable)
	if setBanner.Err != nil {
		t.Fatalf("profile set banner failed: %v", setBanner.Err)
	}

	show := engine.ExecuteTokens(context.Background(), session, []string{"term", "profile"}, ModeExecutable)
	if show.Err != nil {
		t.Fatalf("profile show failed: %v", show.Err)
	}
	if !strings.Contains(show.Output, "prompt_label=INDUS-PRO") || !strings.Contains(show.Output, "show_banner=false") {
		t.Fatalf("unexpected profile output: %q", show.Output)
	}

	reloaded := newTestEngineAt(t, root)
	showReloaded := reloaded.ExecuteTokens(context.Background(), reloaded.NewSession(t.TempDir()), []string{"term", "profile"}, ModeExecutable)
	if showReloaded.Err != nil {
		t.Fatalf("reloaded profile show failed: %v", showReloaded.Err)
	}
	if !strings.Contains(showReloaded.Output, "prompt_label=INDUS-PRO") {
		t.Fatalf("expected prompt label to persist, got %q", showReloaded.Output)
	}
}

func TestTaskCreateRunAndRemove(t *testing.T) {
	engine := newTestEngine(t)
	workspace := t.TempDir()
	session := engine.NewSession(workspace)

	create := engine.ExecuteTokens(context.Background(), session, []string{"task", "create", "smoke", "--commands", "ind env set INDUS_STAGE ready && ind env list"}, ModeExecutable)
	if create.Err != nil {
		t.Fatalf("task create failed: %v", create.Err)
	}

	run := engine.ExecuteTokens(context.Background(), session, []string{"task", "run", "smoke"}, ModeExecutable)
	if run.Err != nil {
		t.Fatalf("task run failed: %v\noutput=%s", run.Err, run.Output)
	}
	if !strings.Contains(run.Output, "INDUS_STAGE=ready") || !strings.Contains(run.Output, "status=ok") {
		t.Fatalf("unexpected task run output: %q", run.Output)
	}

	show := engine.ExecuteTokens(context.Background(), session, []string{"task", "show", "smoke"}, ModeExecutable)
	if show.Err != nil {
		t.Fatalf("task show failed: %v", show.Err)
	}
	if !strings.Contains(show.Output, "last_run_at=") {
		t.Fatalf("expected last_run_at in output: %q", show.Output)
	}

	remove := engine.ExecuteTokens(context.Background(), session, []string{"task", "remove", "smoke"}, ModeExecutable)
	if remove.Err != nil {
		t.Fatalf("task remove failed: %v", remove.Err)
	}
}

func TestUpdateCheckAndDownloadUsesReleaseAPI(t *testing.T) {
	engine := newTestEngine(t)
	session := engine.NewSession(t.TempDir())

	assetPayload := []byte("indus-binary")
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/latest":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"tag_name":     "v1.5.6",
				"name":         "INDUS v1.5.6",
				"published_at": "2026-04-25T00:00:00Z",
				"assets": []map[string]any{
					{
						"name":                 "indus-v1.5.6-windows-amd64.exe",
						"browser_download_url": serverURL + "/asset.exe",
						"size":                 len(assetPayload),
					},
				},
			})
		case "/asset.exe":
			_, _ = w.Write(assetPayload)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	t.Setenv("INDUS_UPDATE_API", server.URL+"/latest")

	check := engine.ExecuteTokens(context.Background(), session, []string{"update"}, ModeExecutable)
	if check.Err != nil {
		t.Fatalf("update check failed: %v", check.Err)
	}
	if !strings.Contains(check.Output, "status=update_available") {
		t.Fatalf("unexpected update check output: %q", check.Output)
	}

	target := filepath.Join(t.TempDir(), "indus-update.exe")
	download := engine.ExecuteTokens(context.Background(), session, []string{"update", "download", "--output", target}, ModeExecutable)
	if download.Err != nil {
		t.Fatalf("update download failed: %v", download.Err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read downloaded asset: %v", err)
	}
	if string(data) != string(assetPayload) {
		t.Fatalf("unexpected downloaded bytes: %q", string(data))
	}
}
