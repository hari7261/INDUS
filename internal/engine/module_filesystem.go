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

type filesystemModule struct {
	engine *Engine
}

func (m *filesystemModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "fs tree":
		return m.tree(inv)
	case "fs find":
		return m.find(inv)
	case "fs inspect":
		return m.inspect(inv)
	case "fs size":
		return m.size(inv)
	case "fs sync":
		return m.sync(inv)
	case "fs digest":
		return m.digest(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *filesystemModule) tree(inv Invocation) Response {
	root := inv.Session.cwd
	if len(inv.Parsed.Positionals) > 0 {
		root = inv.Parsed.Positionals[0]
	}

	depth, err := inv.Parsed.Int(2, "depth")
	if err != nil || depth < 0 {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --depth value")}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "root=%s\n", root)
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil || rel == "." {
			return nil
		}
		level := strings.Count(rel, string(os.PathSeparator))
		if level >= depth {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		prefix := strings.Repeat("  ", level)
		fmt.Fprintf(buffer, "%s%s\n", prefix, entry.Name())
		return nil
	})
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *filesystemModule) find(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind fs find <pattern> [--path .]")}
	}

	pattern := strings.ToLower(inv.Parsed.Positionals[0])
	root := inv.Parsed.String("path")
	if root == "" {
		root = inv.Session.cwd
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	matches := 0
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(entry.Name()), pattern) {
			fmt.Fprintln(buffer, path)
			matches++
		}
		return nil
	})
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	if matches == 0 {
		return Response{Output: "matches=0"}
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *filesystemModule) inspect(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind fs inspect <path>")}
	}
	path := inv.Parsed.Positionals[0]
	info, err := os.Stat(path)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "path=%s\n", path)
	fmt.Fprintf(buffer, "type=%s\n", fileType(info))
	fmt.Fprintf(buffer, "size=%s\n", humanBytes(info.Size()))
	fmt.Fprintf(buffer, "modified=%s\n", info.ModTime().Format(time.RFC3339))
	fmt.Fprintf(buffer, "mode=%s\n", info.Mode().String())
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *filesystemModule) size(inv Invocation) Response {
	path := inv.Session.cwd
	if len(inv.Parsed.Positionals) > 0 {
		path = inv.Parsed.Positionals[0]
	}

	info, err := os.Stat(path)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	if !info.IsDir() {
		return Response{Output: fmt.Sprintf("path=%s\nsize=%s\nfiles=1", path, humanBytes(info.Size()))}
	}

	total, count, err := directorySize(path)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: fmt.Sprintf("path=%s\nsize=%s\nfiles=%d", path, humanBytes(total), count)}
}

func (m *filesystemModule) sync(inv Invocation) Response {
	if len(inv.Parsed.Positionals) < 2 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind fs sync <source> <target>")}
	}

	source := inv.Parsed.Positionals[0]
	target := inv.Parsed.Positionals[1]
	if err := copyTree(source, target); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: fmt.Sprintf("synced_from=%s\nsynced_to=%s", source, target)}
}

func (m *filesystemModule) digest(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind fs digest <path>")}
	}

	digest, err := sha256Digest(inv.Parsed.Positionals[0])
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: fmt.Sprintf("sha256=%s", digest)}
}

func fileType(info os.FileInfo) string {
	if info.IsDir() {
		return "directory"
	}
	return "file"
}
