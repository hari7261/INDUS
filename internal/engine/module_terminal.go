package engine

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type terminalModule struct {
	engine *Engine
}

func (m *terminalModule) Execute(_ context.Context, inv Invocation) Response {
	switch inv.Path {
	case "term clearx":
		return Response{Effects: Effects{ClearScreen: true}}
	case "term profile":
		return m.profile(inv)
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
		state.Profile = defaultTerminalProfile()
	})
	return Response{
		Output: "theme=saffron\nprofile=default\ncache=cleared",
		Effects: Effects{
			Theme: defaultTheme(),
		},
	}
}

func (m *terminalModule) doctor(inv Invocation) Response {
	profile := m.engine.StateSnapshot().Profile
	return Response{Output: fmt.Sprintf("theme=%s\ninteractive=%t\nmetrics=%d\nbanner=%t\nanimation=%s\nprompt=%s", inv.Session.theme.Name, inv.Mode == ModeInteractive, len(m.engine.Metrics()), profile.ShowBanner, profile.BannerAnimation, profile.PromptLabel)}
}

func (m *terminalModule) profile(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return m.profileShow()
	}

	switch strings.ToLower(inv.Parsed.Positionals[0]) {
	case "show":
		return m.profileShow()
	case "reset":
		if err := m.engine.state.Update(func(state *PersistentState) {
			state.Profile = defaultTerminalProfile()
		}); err != nil {
			return Response{Err: commandFailedError(inv.Command, err)}
		}
		return Response{Output: "profile=default"}
	case "set":
		return m.profileSet(inv)
	default:
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind term profile [show|set|reset]")}
	}
}

func (m *terminalModule) profileShow() Response {
	profile := m.engine.StateSnapshot().Profile
	return Response{Output: fmt.Sprintf(
		"show_banner=%t\nbanner_animation=%s\nbanner_duration_ms=%d\ncompact_mode=%t\nprompt_label=%s",
		profile.ShowBanner,
		profile.BannerAnimation,
		profile.BannerDurationMS,
		profile.CompactMode,
		profile.PromptLabel,
	)}
}

func (m *terminalModule) profileSet(inv Invocation) Response {
	if len(inv.Parsed.Positionals) < 3 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind term profile set <key> <value>")}
	}

	key := strings.ToLower(inv.Parsed.Positionals[1])
	value := strings.TrimSpace(strings.Join(inv.Parsed.Positionals[2:], " "))

	var output string
	err := m.engine.state.Update(func(state *PersistentState) {
		profile := state.Profile
		switch key {
		case "banner":
			parsed, ok := parseProfileBool(value)
			if !ok {
				output = ""
				return
			}
			profile.ShowBanner = parsed
			output = fmt.Sprintf("show_banner=%t", profile.ShowBanner)
		case "animation":
			switch strings.ToLower(value) {
			case "mascot-wave", "static", "none":
				profile.BannerAnimation = strings.ToLower(value)
				output = "banner_animation=" + profile.BannerAnimation
			default:
				output = ""
			}
		case "duration":
			duration, parseErr := time.ParseDuration(value)
			if parseErr == nil {
				profile.BannerDurationMS = int(duration / time.Millisecond)
				output = fmt.Sprintf("banner_duration_ms=%d", profile.BannerDurationMS)
				break
			}
			if milliseconds, intErr := strconv.Atoi(value); intErr == nil && milliseconds > 0 {
				profile.BannerDurationMS = milliseconds
				output = fmt.Sprintf("banner_duration_ms=%d", profile.BannerDurationMS)
			}
		case "compact":
			parsed, ok := parseProfileBool(value)
			if !ok {
				output = ""
				return
			}
			profile.CompactMode = parsed
			output = fmt.Sprintf("compact_mode=%t", profile.CompactMode)
		case "prompt":
			profile.PromptLabel = value
			output = "prompt_label=" + profile.PromptLabel
		default:
			output = ""
		}

		state.Profile = normalizeTerminalProfile(profile)
	})
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	if output == "" {
		return Response{Err: invalidArgumentError(inv.Command, "unsupported profile value")}
	}
	return Response{Output: output}
}

func parseProfileBool(value string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on", "enable", "enabled":
		return true, true
	case "0", "false", "no", "off", "disable", "disabled":
		return false, true
	default:
		return false, false
	}
}
