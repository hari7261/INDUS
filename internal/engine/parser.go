package engine

import (
	"strconv"
	"strings"
	"time"
)

type ParsedArgs struct {
	Flags       map[string]string
	BoolFlags   map[string]bool
	Positionals []string
}

func ParseCommandLine(line string) []string {
	var args []string
	var current strings.Builder
	var quote rune
	escaped := false

	for _, char := range line {
		switch {
		case escaped:
			current.WriteRune(char)
			escaped = false
		case char == '\\':
			escaped = true
		case quote != 0:
			if char == quote {
				quote = 0
			} else {
				current.WriteRune(char)
			}
		case char == '"' || char == '\'':
			quote = char
		case char == ' ' || char == '\t':
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func ParseArgs(args []string) ParsedArgs {
	parsed := ParsedArgs{
		Flags:     map[string]string{},
		BoolFlags: map[string]bool{},
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case strings.HasPrefix(arg, "--") && strings.Contains(arg, "="):
			parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
			parsed.Flags[strings.ToLower(parts[0])] = parts[1]
		case strings.HasPrefix(arg, "--"):
			name := strings.ToLower(strings.TrimPrefix(arg, "--"))
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				parsed.Flags[name] = args[i+1]
				i++
				continue
			}
			parsed.BoolFlags[name] = true
		case strings.HasPrefix(arg, "-") && len(arg) > 1:
			name := strings.ToLower(strings.TrimPrefix(arg, "-"))
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				parsed.Flags[name] = args[i+1]
				i++
				continue
			}
			parsed.BoolFlags[name] = true
		default:
			parsed.Positionals = append(parsed.Positionals, arg)
		}
	}

	return parsed
}

func (p ParsedArgs) String(keys ...string) string {
	for _, key := range keys {
		key = strings.ToLower(key)
		if value, ok := p.Flags[key]; ok {
			return value
		}
	}
	return ""
}

func (p ParsedArgs) Bool(keys ...string) bool {
	for _, key := range keys {
		key = strings.ToLower(key)
		if value, ok := p.BoolFlags[key]; ok && value {
			return true
		}
	}
	return false
}

func (p ParsedArgs) Int(defaultValue int, keys ...string) (int, error) {
	value := p.String(keys...)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}

func (p ParsedArgs) Duration(defaultValue time.Duration, keys ...string) (time.Duration, error) {
	value := p.String(keys...)
	if value == "" {
		return defaultValue, nil
	}
	return time.ParseDuration(value)
}
