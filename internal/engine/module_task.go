package engine

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

type taskModule struct {
	engine *Engine
}

func (m *taskModule) Execute(ctx context.Context, inv Invocation) Response {
	switch inv.Path {
	case "task create":
		return m.create(inv)
	case "task list":
		return m.list()
	case "task show":
		return m.show(inv)
	case "task run":
		return m.run(ctx, inv)
	case "task remove":
		return m.remove(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *taskModule) create(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind task create <name> --commands \"ind doctor && ind status\"")}
	}

	name := strings.ToLower(strings.TrimSpace(inv.Parsed.Positionals[0]))
	if name == "" {
		return Response{Err: invalidArgumentError(inv.Command, "task name cannot be empty")}
	}

	commandBlob := inv.Parsed.String("commands")
	if commandBlob == "" && len(inv.Parsed.Positionals) > 1 {
		commandBlob = strings.Join(inv.Parsed.Positionals[1:], " ")
	}

	commands := splitTaskCommands(commandBlob)
	if len(commands) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "provide one or more commands with --commands")}
	}

	now := time.Now().Format(time.RFC3339)
	if err := m.engine.state.Update(func(state *PersistentState) {
		record, ok := state.Tasks[name]
		if !ok {
			record = TaskRecord{
				Name:      name,
				CreatedAt: now,
			}
		}
		record.Name = name
		record.Commands = append([]string(nil), commands...)
		record.UpdatedAt = now
		state.Tasks[name] = record
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("task=%s\ncommands=%d\nstatus=saved", name, len(commands))}
}

func (m *taskModule) list() Response {
	state := m.engine.StateSnapshot()
	if len(state.Tasks) == 0 {
		return Response{Output: "tasks=0"}
	}

	names := make([]string, 0, len(state.Tasks))
	for name := range state.Tasks {
		names = append(names, name)
	}
	sort.Strings(names)

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "tasks=%d\n", len(names))
	for _, name := range names {
		task := state.Tasks[name]
		fmt.Fprintf(buffer, "%s commands=%d", task.Name, len(task.Commands))
		if task.LastRunAt != "" {
			fmt.Fprintf(buffer, " last_run=%s", task.LastRunAt)
		}
		fmt.Fprintln(buffer)
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *taskModule) show(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind task show <name>")}
	}

	name := strings.ToLower(inv.Parsed.Positionals[0])
	task, ok := m.engine.StateSnapshot().Tasks[name]
	if !ok {
		return Response{Err: invalidArgumentError(inv.Command, "task not found")}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "task=%s\n", task.Name)
	fmt.Fprintf(buffer, "commands=%d\n", len(task.Commands))
	if task.CreatedAt != "" {
		fmt.Fprintf(buffer, "created_at=%s\n", task.CreatedAt)
	}
	if task.UpdatedAt != "" {
		fmt.Fprintf(buffer, "updated_at=%s\n", task.UpdatedAt)
	}
	if task.LastRunAt != "" {
		fmt.Fprintf(buffer, "last_run_at=%s\n", task.LastRunAt)
	}
	for i, command := range task.Commands {
		fmt.Fprintf(buffer, "%d. %s\n", i+1, command)
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *taskModule) run(ctx context.Context, inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind task run <name>")}
	}

	name := strings.ToLower(inv.Parsed.Positionals[0])
	task, ok := m.engine.StateSnapshot().Tasks[name]
	if !ok {
		return Response{Err: invalidArgumentError(inv.Command, "task not found")}
	}

	continueOnError := inv.Parsed.Bool("continue")
	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	failures := 0
	for i, command := range task.Commands {
		fmt.Fprintf(buffer, "[step %d/%d] %s\n", i+1, len(task.Commands), command)
		response := m.engine.ExecuteTokens(ctx, inv.Session, ParseCommandLine(command), ModeExecutable)
		applyTaskEffects(inv.Session, response.Effects)

		if response.Warning != "" {
			fmt.Fprintf(buffer, "warning=%s\n", strings.ReplaceAll(response.Warning, "\n", " | "))
		}
		if response.Output != "" {
			fmt.Fprintln(buffer, response.Output)
		}
		if response.Err != nil {
			failures++
			fmt.Fprintf(buffer, "error=%s\n", response.Err.Message)
			if !continueOnError {
				return Response{
					Output: strings.TrimSpace(buffer.String()),
					Err: &IndError{
						Code:       response.Err.Code,
						Command:    inv.Command,
						Message:    fmt.Sprintf("task %s failed at step %d", name, i+1),
						Suggestion: response.Err.Suggestion,
					},
				}
			}
		}
	}

	_ = m.engine.state.Update(func(state *PersistentState) {
		record := state.Tasks[name]
		record.LastRunAt = time.Now().Format(time.RFC3339)
		state.Tasks[name] = record
	})

	fmt.Fprintf(buffer, "task=%s\nstatus=%s\nfailures=%d", name, ternaryTaskStatus(failures == 0, "ok", "completed_with_errors"), failures)
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *taskModule) remove(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind task remove <name>")}
	}

	name := strings.ToLower(inv.Parsed.Positionals[0])
	if _, ok := m.engine.StateSnapshot().Tasks[name]; !ok {
		return Response{Err: invalidArgumentError(inv.Command, "task not found")}
	}

	if err := m.engine.state.Update(func(state *PersistentState) {
		delete(state.Tasks, name)
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("task=%s\nstatus=removed", name)}
}

func splitTaskCommands(value string) []string {
	if value == "" {
		return nil
	}
	normalized := strings.NewReplacer("\r\n", "\n", "&&", "\n", ";", "\n").Replace(value)
	parts := strings.Split(normalized, "\n")
	commands := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			commands = append(commands, part)
		}
	}
	return commands
}

func applyTaskEffects(session *Session, effects Effects) {
	if session == nil {
		return
	}
	if effects.NextDir != "" {
		session.SetCWD(effects.NextDir)
	}
	if effects.Theme.Name != "" {
		session.SetTheme(effects.Theme)
	}
}

func ternaryTaskStatus(ok bool, yes, no string) string {
	if ok {
		return yes
	}
	return no
}
