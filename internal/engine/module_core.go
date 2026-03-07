package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type coreModule struct {
	engine *Engine
}

func (m *coreModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "about":
		return m.about()
	case "doctor":
		return m.doctor(inv)
	case "docs":
		return m.docs()
	case "scan":
		return m.scan(inv)
	case "status":
		return m.status(inv)
	case "version":
		return m.version()
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *coreModule) about() Response {
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "product=INDUS Terminal\n")
	fmt.Fprintf(buffer, "release=%s\n", m.engine.Version())
	fmt.Fprintf(buffer, "registry=%s\n", m.engine.RegistryVersion())
	fmt.Fprintf(buffer, "commands=%d\n", len(m.engine.registry))
	fmt.Fprintf(buffer, "engine=registry-backed native command bus\n")
	fmt.Fprintf(buffer, "features=lazy loading,caching,async safety,versioned docs\n")
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *coreModule) version() Response {
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "version=%s\n", m.engine.Version())
	fmt.Fprintf(buffer, "commit=%s\n", m.engine.Commit())
	fmt.Fprintf(buffer, "build_time=%s\n", m.engine.BuildTime())
	fmt.Fprintf(buffer, "registry=%s\n", m.engine.RegistryVersion())
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *coreModule) docs() Response {
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "index=%s\n", filepath.Join(m.engine.paths.DocsDir, "index.html"))
	fmt.Fprintf(buffer, "commands=%s\n", filepath.Join(m.engine.paths.DocsDir, "commands.html"))
	fmt.Fprintf(buffer, "versions=%s\n", filepath.Join(m.engine.paths.DocsDir, "versions.html"))
	fmt.Fprintf(buffer, "v1_3_0=%s\n", filepath.Join(m.engine.paths.DocsDir, "version", "v1.3.0.html"))
	fmt.Fprintf(buffer, "v1_4_0=%s\n", filepath.Join(m.engine.paths.DocsDir, "version", "v1.4.0.html"))
	fmt.Fprintf(buffer, "command_index=%s\n", filepath.Join(m.engine.paths.DocsDir, "commands.json"))
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *coreModule) status(inv Invocation) Response {
	state := m.engine.StateSnapshot()
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "cwd=%s\n", inv.Session.cwd)
	fmt.Fprintf(buffer, "theme=%s\n", inv.Session.theme.Name)
	fmt.Fprintf(buffer, "active_workspace=%s\n", state.ActiveWorkspace)
	fmt.Fprintf(buffer, "managed_env=%d\n", len(state.ManagedEnv))
	fmt.Fprintf(buffer, "installed_packages=%d\n", len(state.Packages))
	fmt.Fprintf(buffer, "cache_entries=%d\n", m.engine.cache.Size())
	fmt.Fprintf(buffer, "recent_metrics=%d\n", len(m.engine.Metrics()))
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *coreModule) scan(inv Invocation) Response {
	state := m.engine.StateSnapshot()
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "scan=ok\n")
	fmt.Fprintf(buffer, "cwd=%s\n", inv.Session.cwd)
	fmt.Fprintf(buffer, "registry_commands=%d\n", len(m.engine.registry))
	fmt.Fprintf(buffer, "docs_present=%t\n", exists(m.engine.paths.DocsDir))
	fmt.Fprintf(buffer, "state_present=%t\n", exists(m.engine.paths.StateFile))
	fmt.Fprintf(buffer, "managed_env=%d\n", len(state.ManagedEnv))
	fmt.Fprintf(buffer, "workspaces=%d\n", len(state.Workspaces))
	fmt.Fprintf(buffer, "packages=%d\n", len(state.Packages))
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *coreModule) doctor(inv Invocation) Response {
	state := m.engine.StateSnapshot()
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "registry_ok=%t\n", len(m.engine.registry) >= 50)
	fmt.Fprintf(buffer, "docs_ok=%t\n", exists(m.engine.paths.DocsDir))
	fmt.Fprintf(buffer, "state_ok=%t\n", exists(m.engine.paths.StateFile))
	fmt.Fprintf(buffer, "theme_ok=%t\n", inv.Session.theme.Name != "")
	fmt.Fprintf(buffer, "workspace_ok=%t\n", state.ActiveWorkspace == "" || exists(state.ActiveWorkspace))
	fmt.Fprintf(buffer, "doctor_hint=run ind sys doctor for file-level diagnostics\n")
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
