package engine

import (
	"context"
	"fmt"
	"strings"
)

type terminalModule struct {
	engine *Engine
}

func (m *terminalModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "term clearx":
		return Response{Effects: Effects{ClearScreen: true}}
	case "term theme":
		return m.theme(inv)
	case "term history":
		return m.history(inv)
	case "term speed":
		return m.speed()
	case "term reset":
		return m.reset()
	case "term doctor":
		return m.doctor(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *terminalModule) theme(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Output: "themes=" + strings.Join(themeNames(), ",")}
	}

	themeName := inv.Parsed.Positionals[0]
	theme, ok := themeByName(themeName)
	if !ok {
		return Response{Err: invalidArgumentError(inv.Command, "unknown theme")}
	}

	if err := m.engine.state.Update(func(state *PersistentState) {
		state.Theme = theme.Name
	}); err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{
		Output: fmt.Sprintf("theme=%s", theme.Name),
		Effects: Effects{
			Theme: theme,
		},
	}
}

func (m *terminalModule) history(inv Invocation) Response {
	metrics := m.engine.Metrics()
	if len(metrics) == 0 {
		return Response{Output: "history=0"}
	}

	limit, err := inv.Parsed.Int(10, "limit")
	if err != nil || limit <= 0 {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --limit value")}
	}

	if len(metrics) < limit {
		limit = len(metrics)
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	for _, metric := range metrics[len(metrics)-limit:] {
		fmt.Fprintf(buffer, "%s cached=%t duration=%s\n", metric.Command, metric.Cached, metric.Duration)
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *terminalModule) speed() Response {
	metrics := m.engine.Metrics()
	if len(metrics) == 0 {
		return Response{Output: "metrics=0"}
	}

	var total int64
	cacheHits := 0
	for _, metric := range metrics {
		total += metric.Duration.Milliseconds()
		if metric.Cached {
			cacheHits++
		}
	}
	average := float64(total) / float64(len(metrics))
	return Response{Output: fmt.Sprintf("metrics=%d\navg_ms=%.2f\ncache_hits=%d", len(metrics), average, cacheHits)}
}

func (m *terminalModule) reset() Response {
	m.engine.cache.Clear()
	_ = m.engine.state.Update(func(state *PersistentState) {
		state.Theme = defaultTheme().Name
	})
	return Response{
		Output: "theme=saffron\ncache=cleared",
		Effects: Effects{
			Theme: defaultTheme(),
		},
	}
}

func (m *terminalModule) doctor(inv Invocation) Response {
	return Response{Output: fmt.Sprintf("theme=%s\ninteractive=%t\nmetrics=%d", inv.Session.theme.Name, inv.Mode == ModeInteractive, len(m.engine.Metrics()))}
}
