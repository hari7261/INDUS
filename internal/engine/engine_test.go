package engine

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestEngine(t testing.TB) *Engine {
	t.Helper()

	root := t.TempDir()
	t.Setenv("APPDATA", root)
	t.Setenv("LOCALAPPDATA", root)
	t.Setenv("HOME", root)
	t.Setenv("USERPROFILE", root)
	t.Setenv("INDUS_CONFIG", filepath.Join(root, "config.cfg"))

	engine, err := New(Options{
		Version:   "1.5.1-test",
		Commit:    "test",
		BuildTime: "2026-03-08T00:00:00Z",
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
	if engine.RegistryVersion() != "1.5.1" {
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
