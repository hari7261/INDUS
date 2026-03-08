package engine

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

type catalogEntry struct {
	Name        string
	Version     string
	Description string
}

var packageCatalog = map[string]catalogEntry{
	"aurora-kit":  {Name: "aurora-kit", Version: "1.4.0", Description: "UI primitives for versioned documentation"},
	"cache-lens":  {Name: "cache-lens", Version: "1.2.1", Description: "runtime cache inspection pack"},
	"env-vault":   {Name: "env-vault", Version: "1.1.3", Description: "managed environment snapshots"},
	"flux-board":  {Name: "flux-board", Version: "2.0.0", Description: "workspace flow helpers"},
	"net-mapper":  {Name: "net-mapper", Version: "1.5.2", Description: "network probing utilities"},
	"orbit-kit":   {Name: "orbit-kit", Version: "0.9.8", Description: "project manifest helpers"},
	"pulse-core":  {Name: "pulse-core", Version: "1.4.4", Description: "runtime telemetry adapters"},
	"trace-slate": {Name: "trace-slate", Version: "1.0.4", Description: "developer trace report pack"},
}

type packageModule struct {
	engine *Engine
}

func (m *packageModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "pkg search":
		return m.search(inv)
	case "pkg install":
		return m.install(inv)
	case "pkg remove":
		return m.remove(inv)
	case "pkg update":
		return m.update(inv)
	case "pkg list":
		return m.list()
	case "pkg audit":
		return m.audit()
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *packageModule) search(inv Invocation) Response {
	query := ""
	if len(inv.Parsed.Positionals) > 0 {
		query = strings.ToLower(inv.Parsed.Positionals[0])
	}

	names := make([]string, 0, len(packageCatalog))
	for name := range packageCatalog {
		names = append(names, name)
	}
	sort.Strings(names)

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	for _, name := range names {
		entry := packageCatalog[name]
		if query == "" || strings.Contains(name, query) || strings.Contains(strings.ToLower(entry.Description), query) {
			fmt.Fprintf(buffer, "%s %s %s\n", entry.Name, entry.Version, entry.Description)
		}
	}

	if buffer.Len() == 0 {
		return Response{Output: "matches=0"}
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *packageModule) install(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind pkg install <name>")}
	}
	name := inv.Parsed.Positionals[0]
	entry, ok := packageCatalog[name]
	if !ok {
		return Response{Err: invalidArgumentError(inv.Command, "package not found in INDUS catalog")}
	}

	if err := m.engine.state.Update(func(state *PersistentState) {
		state.Packages[name] = PackageRecord{
			Name:        entry.Name,
			Version:     entry.Version,
			InstalledAt: time.Now().Format(time.RFC3339),
		}
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("installed=%s\nversion=%s", entry.Name, entry.Version)}
}

func (m *packageModule) remove(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind pkg remove <name>")}
	}
	name := inv.Parsed.Positionals[0]
	if err := m.engine.state.Update(func(state *PersistentState) {
		delete(state.Packages, name)
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("removed=%s", name)}
}

func (m *packageModule) update(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind pkg update <name>")}
	}
	name := inv.Parsed.Positionals[0]
	entry, ok := packageCatalog[name]
	if !ok {
		return Response{Err: invalidArgumentError(inv.Command, "package not found in INDUS catalog")}
	}

	if err := m.engine.state.Update(func(state *PersistentState) {
		if _, installed := state.Packages[name]; installed {
			state.Packages[name] = PackageRecord{
				Name:        entry.Name,
				Version:     entry.Version,
				InstalledAt: time.Now().Format(time.RFC3339),
			}
		}
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	return Response{Output: fmt.Sprintf("updated=%s\nversion=%s", name, entry.Version)}
}

func (m *packageModule) list() Response {
	state := m.engine.StateSnapshot()
	if len(state.Packages) == 0 {
		return Response{Output: "installed_packages=0"}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)
	for _, key := range sortedPackageKeys(state.Packages) {
		entry := state.Packages[key]
		fmt.Fprintf(buffer, "%s %s installed_at=%s\n", entry.Name, entry.Version, entry.InstalledAt)
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *packageModule) audit() Response {
	state := m.engine.StateSnapshot()
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	if len(state.Packages) == 0 {
		return Response{Output: "audit=clean\npackages=0"}
	}

	status := "clean"
	for _, key := range sortedPackageKeys(state.Packages) {
		entry := state.Packages[key]
		catalog, ok := packageCatalog[key]
		if !ok {
			status = "drift"
			fmt.Fprintf(buffer, "%s missing_from_catalog\n", entry.Name)
			continue
		}
		if catalog.Version != entry.Version {
			status = "drift"
			fmt.Fprintf(buffer, "%s installed=%s catalog=%s\n", entry.Name, entry.Version, catalog.Version)
		}
	}

	if buffer.Len() == 0 {
		fmt.Fprintln(buffer, "audit=clean")
		fmt.Fprintf(buffer, "packages=%d", len(state.Packages))
		return Response{Output: strings.TrimSpace(buffer.String())}
	}

	return Response{Output: fmt.Sprintf("audit=%s\n%s", status, strings.TrimSpace(buffer.String()))}
}
