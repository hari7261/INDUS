package engine

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type projectModule struct {
	engine *Engine
}

func (m *projectModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "proj create":
		return m.create(inv)
	case "proj init":
		return m.init(inv)
	case "proj list":
		return m.list(inv)
	case "proj clean":
		return m.clean(inv)
	case "proj build":
		return m.build(inv)
	case "proj run":
		return m.run(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *projectModule) create(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind proj create <name> [--dir .]")}
	}

	name := inv.Parsed.Positionals[0]
	base := inv.Parsed.String("dir")
	if base == "" {
		base = inv.Session.cwd
	}
	root := filepath.Join(base, name)

	if err := ensureProject(root, name, m.engine.Version()); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("project=%s\nroot=%s", name, root)}
}

func (m *projectModule) init(inv Invocation) Response {
	name := inv.Parsed.String("name")
	if name == "" {
		name = filepath.Base(inv.Session.cwd)
	}

	if err := ensureProject(inv.Session.cwd, name, m.engine.Version()); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: fmt.Sprintf("project=%s\nroot=%s", name, inv.Session.cwd)}
}

func (m *projectModule) list(inv Invocation) Response {
	root := inv.Parsed.String("path")
	if root == "" {
		root = inv.Session.cwd
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	found := 0
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && entry.Name() == ".git" {
			return filepath.SkipDir
		}
		if !entry.IsDir() || entry.Name() != "project.json" || filepath.Base(filepath.Dir(path)) != ".indus" {
			return nil
		}
		var manifest projectManifest
		if err := readJSONFile(path, &manifest); err == nil {
			fmt.Fprintf(buffer, "%s %s\n", manifest.Name, filepath.Dir(filepath.Dir(path)))
			found++
		}
		return nil
	})
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	if found == 0 {
		return Response{Output: "projects=0"}
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *projectModule) clean(inv Invocation) Response {
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
	return Response{Output: fmt.Sprintf("project_root=%s\ncleaned=%d", root, removed)}
}

func (m *projectModule) build(inv Invocation) Response {
	root := inv.Session.cwd
	if len(inv.Parsed.Positionals) > 0 {
		root = inv.Parsed.Positionals[0]
	}

	var manifest projectManifest
	if err := readJSONFile(projectManifestPath(root), &manifest); err != nil {
		return Response{Err: invalidArgumentError(inv.Command, "project is not initialized")}
	}

	artifact := map[string]any{
		"project":   manifest.Name,
		"version":   manifest.Version,
		"built_at":  time.Now().Format(time.RFC3339),
		"engine":    m.engine.Version(),
		"workspace": root,
		"artifact":  filepath.Join(root, "build", "indus-artifact.json"),
	}
	path := filepath.Join(root, "build", "indus-artifact.json")
	if err := writeJSONFile(path, artifact); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("build=ok\nartifact=%s", path)}
}

func (m *projectModule) run(inv Invocation) Response {
	root := inv.Session.cwd
	if len(inv.Parsed.Positionals) > 0 {
		root = inv.Parsed.Positionals[0]
	}

	var manifest projectManifest
	if err := readJSONFile(projectManifestPath(root), &manifest); err != nil {
		return Response{Err: invalidArgumentError(inv.Command, "project is not initialized")}
	}

	artifact := filepath.Join(root, "build", "indus-artifact.json")
	if _, err := os.Stat(artifact); err != nil {
		return Response{Output: fmt.Sprintf("project=%s\nmode=simulation\nstatus=build_missing\nnext=ind proj build %s", manifest.Name, root)}
	}

	return Response{Output: fmt.Sprintf("project=%s\nmode=simulation\nstatus=running\nartifact=%s", manifest.Name, artifact)}
}
