package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type workspaceModule struct {
	engine *Engine
}

func (m *workspaceModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "work init":
		return m.init(inv)
	case "work list":
		return m.list()
	case "work switch":
		return m.switchWorkspace(inv)
	case "work clean":
		return m.clean(inv)
	case "work archive":
		return m.archive(inv)
	case "work pin":
		return m.pin(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *workspaceModule) init(inv Invocation) Response {
	name := ""
	if len(inv.Parsed.Positionals) > 0 {
		name = inv.Parsed.Positionals[0]
	}
	if name == "" {
		name = filepath.Base(inv.Session.cwd)
	}

	manifest := workspaceManifest{
		Name:      name,
		Root:      inv.Session.cwd,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if err := writeJSONFile(workspaceManifestPath(inv.Session.cwd), manifest); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	if err := m.engine.state.Update(func(state *PersistentState) {
		state.ActiveWorkspace = inv.Session.cwd
		state.Workspaces = upsertWorkspace(state.Workspaces, WorkspaceRecord{
			Name:     name,
			Path:     inv.Session.cwd,
			LastUsed: time.Now().Format(time.RFC3339),
		})
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("workspace=%s\nroot=%s", name, inv.Session.cwd)}
}

func (m *workspaceModule) list() Response {
	state := m.engine.StateSnapshot()
	if len(state.Workspaces) == 0 {
		return Response{Output: "workspaces=0"}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)
	for _, workspace := range state.Workspaces {
		flag := ""
		if workspace.Pinned {
			flag = " pinned"
		}
		fmt.Fprintf(buffer, "%s %s%s\n", workspace.Name, workspace.Path, flag)
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *workspaceModule) switchWorkspace(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind work switch <name-or-path>")}
	}

	target := inv.Parsed.Positionals[0]
	state := m.engine.StateSnapshot()

	nextDir := target
	for _, workspace := range state.Workspaces {
		if workspace.Name == target {
			nextDir = workspace.Path
			break
		}
	}

	if _, err := os.Stat(nextDir); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	if err := m.engine.state.Update(func(state *PersistentState) {
		state.ActiveWorkspace = nextDir
		for i := range state.Workspaces {
			if state.Workspaces[i].Path == nextDir {
				state.Workspaces[i].LastUsed = time.Now().Format(time.RFC3339)
			}
		}
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{
		Output: fmt.Sprintf("workspace=%s\npath=%s", target, nextDir),
		Effects: Effects{
			NextDir: nextDir,
		},
	}
}

func (m *workspaceModule) clean(inv Invocation) Response {
	root := inv.Session.cwd
	if len(inv.Parsed.Positionals) > 0 {
		root = inv.Parsed.Positionals[0]
	}

	removed := 0
	for _, path := range []string{filepath.Join(root, "build"), filepath.Join(root, ".indus", "cache")} {
		if err := os.RemoveAll(path); err == nil {
			removed++
		}
	}
	return Response{Output: fmt.Sprintf("workspace=%s\ncleaned=%d", root, removed)}
}

func (m *workspaceModule) archive(inv Invocation) Response {
	root := inv.Session.cwd
	if len(inv.Parsed.Positionals) > 0 {
		root = inv.Parsed.Positionals[0]
	}
	name := filepath.Base(root) + "-" + time.Now().Format("20060102-150405") + ".zip"
	target := filepath.Join(m.engine.paths.ReportsDir, name)

	if err := zipDirectory(root, target); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: fmt.Sprintf("archive=%s", target)}
}

func (m *workspaceModule) pin(inv Invocation) Response {
	name := inv.Parsed.String("name")
	state := m.engine.StateSnapshot()
	target := state.ActiveWorkspace
	if len(inv.Parsed.Positionals) > 0 {
		name = inv.Parsed.Positionals[0]
	}

	if err := m.engine.state.Update(func(state *PersistentState) {
		for i := range state.Workspaces {
			if (name != "" && state.Workspaces[i].Name == name) || state.Workspaces[i].Path == target {
				state.Workspaces[i].Pinned = true
			}
		}
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: "workspace_pinned=true"}
}

func upsertWorkspace(existing []WorkspaceRecord, next WorkspaceRecord) []WorkspaceRecord {
	for i := range existing {
		if existing[i].Path == next.Path || existing[i].Name == next.Name {
			existing[i] = next
			return existing
		}
	}
	return append(existing, next)
}
