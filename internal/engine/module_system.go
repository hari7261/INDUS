package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type systemModule struct {
	engine *Engine
}

func (m *systemModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "sys stats":
		return m.stats(inv)
	case "sys info":
		return m.info(inv)
	case "sys clean":
		return m.clean(inv)
	case "sys doctor":
		return m.doctor(inv)
	case "sys watch":
		return m.watch(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *systemModule) stats(inv Invocation) Response {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "runtime_go=%s\n", runtime.Version())
	fmt.Fprintf(buffer, "platform=%s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(buffer, "cpu=%d\n", runtime.NumCPU())
	fmt.Fprintf(buffer, "goroutines=%d\n", runtime.NumGoroutine())
	fmt.Fprintf(buffer, "memory_alloc=%s\n", humanBytes(int64(mem.Alloc)))
	fmt.Fprintf(buffer, "memory_sys=%s\n", humanBytes(int64(mem.Sys)))
	fmt.Fprintf(buffer, "cache_entries=%d\n", m.engine.cache.Size())
	fmt.Fprintf(buffer, "cwd=%s\n", inv.Session.cwd)
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *systemModule) info(inv Invocation) Response {
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	exe, _ := os.Executable()
	home, _ := os.UserHomeDir()

	fmt.Fprintf(buffer, "version=%s\n", m.engine.Version())
	fmt.Fprintf(buffer, "registry=%s\n", m.engine.RegistryVersion())
	fmt.Fprintf(buffer, "commit=%s\n", m.engine.Commit())
	fmt.Fprintf(buffer, "build_time=%s\n", m.engine.BuildTime())
	fmt.Fprintf(buffer, "exe=%s\n", exe)
	fmt.Fprintf(buffer, "root=%s\n", m.engine.paths.RootDir)
	fmt.Fprintf(buffer, "docs=%s\n", m.engine.paths.DocsDir)
	fmt.Fprintf(buffer, "home=%s\n", home)
	fmt.Fprintf(buffer, "cwd=%s\n", inv.Session.cwd)
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *systemModule) clean(inv Invocation) Response {
	removed := 0
	if count := m.engine.cache.Clear(); count > 0 {
		removed += count
	}

	entries, err := os.ReadDir(m.engine.paths.ReportsDir)
	if err == nil {
		for _, entry := range entries {
			if removeErr := os.RemoveAll(filepath.Join(m.engine.paths.ReportsDir, entry.Name())); removeErr == nil {
				removed++
			}
		}
	}

	cacheEntries, err := os.ReadDir(m.engine.paths.CacheDir)
	if err == nil {
		for _, entry := range cacheEntries {
			if removeErr := os.RemoveAll(filepath.Join(m.engine.paths.CacheDir, entry.Name())); removeErr == nil {
				removed++
			}
		}
	}

	return Response{Output: fmt.Sprintf("cleaned_items=%d\nstate=ok", removed)}
}

func (m *systemModule) doctor(inv Invocation) Response {
	results := []string{
		fmt.Sprintf("registry_loaded=%t", len(m.engine.registry) >= 50),
		fmt.Sprintf("command_count=%d", len(m.engine.registry)),
	}

	if _, err := os.Stat(m.engine.paths.DocsDir); err == nil {
		results = append(results, "docs=ok")
	} else {
		results = append(results, "docs=missing")
	}

	if _, err := os.Stat(m.engine.paths.StateFile); err == nil {
		results = append(results, "state=ok")
	} else {
		results = append(results, "state=missing")
	}

	if _, err := os.Stat(m.engine.paths.RegistryPath); err == nil {
		results = append(results, "registry_file=ok")
	} else {
		results = append(results, "registry_file=embedded")
	}

	if len(results) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "no diagnostic results")}
	}

	return Response{Output: strings.Join(results, "\n")}
}

func (m *systemModule) watch(inv Invocation) Response {
	interval, err := inv.Parsed.Duration(1*time.Second, "interval")
	if err != nil {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --interval value")}
	}

	count, err := inv.Parsed.Int(5, "count")
	if err != nil || count <= 0 {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --count value")}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	for i := 0; i < count; i++ {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		fmt.Fprintf(buffer, "sample=%d alloc=%s goroutines=%d\n", i+1, humanBytes(int64(mem.Alloc)), runtime.NumGoroutine())
		if i != count-1 {
			time.Sleep(interval)
		}
	}

	return Response{Output: strings.TrimSpace(buffer.String())}
}
