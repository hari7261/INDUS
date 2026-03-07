package engine

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type environmentModule struct {
	engine *Engine
}

func (m *environmentModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "env list":
		return m.list()
	case "env set":
		return m.set(inv)
	case "env unset":
		return m.unset(inv)
	case "env export":
		return m.export(inv)
	case "env import":
		return m.importEnv(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *environmentModule) list() Response {
	state := m.engine.StateSnapshot()
	if len(state.ManagedEnv) == 0 {
		return Response{Output: "managed_env=0"}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "managed_env=%d\n", len(state.ManagedEnv))
	for _, key := range sortedKeys(state.ManagedEnv) {
		fmt.Fprintf(buffer, "%s=%s\n", key, state.ManagedEnv[key])
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *environmentModule) set(inv Invocation) Response {
	if len(inv.Parsed.Positionals) < 2 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind env set KEY VALUE")}
	}

	key := inv.Parsed.Positionals[0]
	value := inv.Parsed.Positionals[1]

	if err := os.Setenv(key, value); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	if err := m.engine.state.Update(func(state *PersistentState) {
		state.ManagedEnv[key] = value
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("env_set=%s", key)}
}

func (m *environmentModule) unset(inv Invocation) Response {
	if len(inv.Parsed.Positionals) < 1 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind env unset KEY")}
	}

	key := inv.Parsed.Positionals[0]
	_ = os.Unsetenv(key)

	if err := m.engine.state.Update(func(state *PersistentState) {
		delete(state.ManagedEnv, key)
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("env_unset=%s", key)}
}

func (m *environmentModule) export(inv Invocation) Response {
	path := inv.Parsed.String("file")
	if path == "" {
		path = "indus-env.json"
	}

	state := m.engine.StateSnapshot()
	if err := writeJSONFile(path, state.ManagedEnv); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("exported=%s\ncount=%d", path, len(state.ManagedEnv))}
}

func (m *environmentModule) importEnv(inv Invocation) Response {
	path := inv.Parsed.String("file")
	if path == "" {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind env import --file indus-env.json")}
	}

	values := map[string]string{}
	if err := readJSONFile(path, &values); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	for key, value := range values {
		_ = os.Setenv(key, value)
	}

	if err := m.engine.state.Update(func(state *PersistentState) {
		for key, value := range values {
			state.ManagedEnv[key] = value
		}
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("imported=%d\nsource=%s", len(values), path)}
}
