package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"indus/internal/engine"
)

var (
	version   = "1.5.3"
	commit    = "initial"
	buildTime = "2026-03-08T00:00:00Z"
)

func main() {
	enableConsoleFeatures()
	setConsoleTitle("INDUS")

	runtime, err := engine.New(engine.Options{
		Version:   version,
		Commit:    commit,
		BuildTime: buildTime,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, (&engine.IndError{
			Code:       "IND_ERR_005",
			Command:    "ind",
			Message:    err.Error(),
			Suggestion: "verify the registry and rerun \"ind doctor\"",
		}).Render())
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startDir := currentDirectory()

	if len(os.Args) > 1 {
		bootstrap, bootstrapErr := parseBootstrapArgs(os.Args[1:])
		if bootstrapErr != nil {
			fmt.Fprintln(os.Stderr, (&engine.IndError{
				Code:       "IND_ERR_003",
				Command:    "ind",
				Message:    bootstrapErr.Error(),
				Suggestion: "use --cwd <directory> before the command or run \"ind doctor\"",
			}).Render())
			os.Exit(2)
		}

		if bootstrap.CWD != "" {
			absoluteCWD, cwdErr := resolveWorkingDirectory(bootstrap.CWD)
			if cwdErr != nil {
				fmt.Fprintln(os.Stderr, (&engine.IndError{
					Code:       "IND_ERR_003",
					Command:    "ind --cwd",
					Message:    cwdErr.Error(),
					Suggestion: "pass a valid directory path",
				}).Render())
				os.Exit(2)
			}
			if chdirErr := os.Chdir(absoluteCWD); chdirErr != nil {
				fmt.Fprintln(os.Stderr, (&engine.IndError{
					Code:       "IND_ERR_004",
					Command:    "ind --cwd",
					Message:    chdirErr.Error(),
					Suggestion: "verify directory permissions and retry",
				}).Render())
				os.Exit(1)
			}
			startDir = absoluteCWD
		}

		if len(bootstrap.Tokens) > 0 {
			session := runtime.NewSession(startDir)
			response := runtime.ExecuteTokens(ctx, session, bootstrap.Tokens, engine.ModeExecutable)
			renderResponse(response)
			os.Exit(exitCode(response))
		}
	}

	terminal := &Terminal{
		engine:  runtime,
		reader:  bufio.NewReader(os.Stdin),
		session: runtime.NewSession(startDir),
	}
	if err := terminal.Start(ctx); err != nil && err != context.Canceled {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Terminal struct {
	engine  *engine.Engine
	reader  *bufio.Reader
	session *engine.Session
}

func (t *Terminal) Start(ctx context.Context) error {
	clearScreen()
	t.printBanner()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nGoodbye.")
			return nil
		default:
		}

		t.printPrompt()
		line, err := t.reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nGoodbye.")
			return nil
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch strings.ToLower(line) {
		case "exit", "quit":
			fmt.Println("Goodbye.")
			return nil
		case "help", "?":
			fmt.Println(t.engine.HelpText())
			continue
		}

		response := t.engine.ExecuteLine(ctx, t.session, line)
		t.applyEffects(response)
		renderResponse(response)
	}
}

func (t *Terminal) printBanner() {
	const (
		reset   = "\033[0m"
		saffron = "\033[38;5;208m"
		white   = "\033[97m"
		green   = "\033[38;5;46m"
		cyan    = "\033[36m"
	)

	bar := strings.Repeat("=", 68)

	waves := []string{
		"~    ~      ~    ~      ~",
		"  ~    ~      ~    ~     ",
		"    ~    ~      ~    ~   ",
	}

	for i := 0; i < len(waves); i++ {
		frame := waves[i]
		switch i {
		case 0:
			fmt.Printf("%s  %s%s\n", saffron, frame, reset)
		case 1:
			fmt.Printf("%s  %s%s\n", white, frame, reset)
		default:
			fmt.Printf("%s  %s%s\n", green, frame, reset)
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("")
	fmt.Printf("%s%s%s\n", saffron, bar, reset)
	fmt.Printf("%s%s%s\n", white, bar, reset)
	fmt.Printf("%s%s%s\n", green, bar, reset)
	fmt.Println("")

	fmt.Printf("%s  ██╗███╗   ██╗██████╗ ██╗   ██╗███████╗%s%s        ▄▄▄%s\n", saffron, reset, white, reset)
	fmt.Printf("%s  ██║████╗  ██║██╔══██╗██║   ██║██╔════╝%s%s       (◉_◉)%s\n", saffron, reset, white, reset)
	fmt.Printf("%s  ██║██╔██╗ ██║██║  ██║██║   ██║███████╗%s%s      /|███|\\%s\n", white, reset, white, reset)
	fmt.Printf("%s  ██║██║╚██╗██║██║  ██║██║   ██║╚════██║%s%s       /   \\%s\n", white, reset, white, reset)
	fmt.Printf("%s  ██║██║ ╚████║██████╔╝╚██████╔╝███████║%s\n", green, reset)
	fmt.Printf("%s  ╚═╝╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚══════╝%s\n", green, reset)

	fmt.Println("")
	fmt.Printf("%s  Namaste! Welcome to INDUS Terminal v%s%s\n", saffron, version, reset)
	fmt.Printf("%s  Native format: ind <command> [options]%s\n", cyan, reset)
	fmt.Printf("%s  Docs: ind docs | Help: help | Exit: exit%s\n", white, reset)
	fmt.Println("")
}

func (t *Terminal) printPrompt() {
	reset := "\033[0m"
	fmt.Printf("%sINDUS%s %s > ", t.session.Theme().Prompt, reset, t.session.CWD())
}

func (t *Terminal) applyEffects(response engine.Response) {
	if response.Effects.Theme.Name != "" {
		t.session.SetTheme(response.Effects.Theme)
	}
	if response.Effects.NextDir != "" {
		if err := os.Chdir(response.Effects.NextDir); err == nil {
			t.session.SetCWD(response.Effects.NextDir)
		}
	}
	if response.Effects.ClearScreen {
		clearScreen()
		t.printBanner()
	}
}

func renderResponse(response engine.Response) {
	if response.Warning != "" {
		fmt.Fprintln(os.Stderr, response.Warning)
	}
	if response.Err != nil {
		fmt.Fprintln(os.Stderr, response.Err.Render())
		return
	}
	if response.Output != "" {
		fmt.Println(response.Output)
	}
}

func exitCode(response engine.Response) int {
	if response.Err == nil {
		return 0
	}
	switch response.Err.Code {
	case "IND_ERR_001", "IND_ERR_002", "IND_ERR_003":
		return 2
	default:
		return 1
	}
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func currentDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

type bootstrapConfig struct {
	CWD    string
	Tokens []string
}

func parseBootstrapArgs(args []string) (bootstrapConfig, error) {
	cfg := bootstrapConfig{}
	remaining := append([]string(nil), args...)

	for len(remaining) > 0 {
		switch remaining[0] {
		case "--cwd":
			if len(remaining) < 2 {
				return cfg, fmt.Errorf("missing required directory after --cwd")
			}
			cfg.CWD = remaining[1]
			remaining = remaining[2:]
		case "--":
			cfg.Tokens = append([]string(nil), remaining[1:]...)
			return cfg, nil
		default:
			cfg.Tokens = append([]string(nil), remaining...)
			return cfg, nil
		}
	}

	return cfg, nil
}

func resolveWorkingDirectory(path string) (string, error) {
	resolved, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", resolved)
	}
	return resolved, nil
}
