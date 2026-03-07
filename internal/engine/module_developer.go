package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type developerModule struct {
	engine *Engine
}

func (m *developerModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "dev bench":
		return m.bench(inv)
	case "dev watch":
		return m.watch(inv)
	case "dev cache":
		return m.cache(inv)
	case "dev reload":
		return m.reload(inv)
	case "dev debug":
		return m.debug(inv)
	case "dev report":
		return m.report(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *developerModule) bench(inv Invocation) Response {
	commandLine := inv.Parsed.String("command")
	if commandLine == "" {
		commandLine = "ind sys stats"
	}
	runs, err := inv.Parsed.Int(5, "runs")
	if err != nil || runs <= 0 {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --runs value")}
	}

	tokens := ParseCommandLine(commandLine)
	if len(tokens) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "empty benchmark command")}
	}
	if strings.Join(lowerSlice(tokens), " ") == "ind dev bench" {
		return Response{Err: invalidArgumentError(inv.Command, "cannot benchmark dev bench recursively")}
	}

	var total time.Duration
	var peak time.Duration
	for i := 0; i < runs; i++ {
		if inv.Parsed.Bool("fresh") {
			m.engine.cache.Clear()
		}
		start := time.Now()
		response := m.engine.ExecuteTokens(context.Background(), inv.Session, tokens, ModeExecutable)
		if response.Err != nil {
			return response
		}
		elapsed := time.Since(start)
		total += elapsed
		if elapsed > peak {
			peak = elapsed
		}
	}
	average := total / time.Duration(runs)
	return Response{Output: fmt.Sprintf("command=%s\nruns=%d\navg=%s\npeak=%s", commandLine, runs, average, peak)}
}

func (m *developerModule) watch(inv Invocation) Response {
	root := inv.Parsed.String("path")
	if root == "" {
		root = inv.Session.cwd
	}
	seconds, err := inv.Parsed.Int(5, "seconds")
	if err != nil || seconds <= 0 {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --seconds value")}
	}

	before, err := fileCountSnapshot(root)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	time.Sleep(time.Duration(seconds) * time.Second)
	after, err := fileCountSnapshot(root)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	changes := 0
	for path, mod := range after {
		if beforeMod, ok := before[path]; !ok || !beforeMod.Equal(mod) {
			changes++
		}
	}
	for path := range before {
		if _, ok := after[path]; !ok {
			changes++
		}
	}
	return Response{Output: fmt.Sprintf("path=%s\nseconds=%d\nchanges=%d", root, seconds, changes)}
}

func (m *developerModule) cache(inv Invocation) Response {
	if inv.Parsed.Bool("clear") {
		count := m.engine.cache.Clear()
		return Response{Output: fmt.Sprintf("cache_cleared=%d", count)}
	}
	return Response{Output: fmt.Sprintf("cache_entries=%d", m.engine.cache.Size())}
}

func (m *developerModule) reload(inv Invocation) Response {
	if err := m.engine.Reload(); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: fmt.Sprintf("registry=%s\ncommands=%d", m.engine.RegistryVersion(), len(m.engine.registry))}
}

func (m *developerModule) debug(inv Invocation) Response {
	paths := m.engine.Paths()
	return Response{Output: fmt.Sprintf("cwd=%s\ntheme=%s\nroot=%s\nstate=%s\ncache=%s\nreports=%s", inv.Session.cwd, inv.Session.theme.Name, paths.RootDir, paths.StateFile, paths.CacheDir, paths.ReportsDir)}
}

func (m *developerModule) report(inv Invocation) Response {
	output := inv.Parsed.String("output")
	if output == "" {
		output = filepath.Join(m.engine.paths.ReportsDir, "indus-report.json")
	}

	report := map[string]any{
		"version":   m.engine.Version(),
		"registry":  m.engine.RegistryVersion(),
		"paths":     m.engine.Paths(),
		"state":     m.engine.StateSnapshot(),
		"metrics":   m.engine.Metrics(),
		"generated": time.Now().Format(time.RFC3339),
	}

	if err := writeJSONFile(output, report); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	if _, err := os.Stat(output); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: fmt.Sprintf("report=%s", output)}
}
