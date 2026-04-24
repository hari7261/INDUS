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
	version   = "1.5.5"
	commit    = "initial"
	buildTime = "2026-04-25T00:00:00Z"
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
		if response.Effects.Exit {
			fmt.Println("Restarting INDUS to apply changes...")
			return nil
		}
	}
}

func (t *Terminal) printBanner() {
	profile := t.engine.StateSnapshot().Profile
	if !profile.ShowBanner {
		return
	}

	const (
		reset   = "\033[0m"
		saffron = "\033[38;5;208m"
		white   = "\033[97m"
		green   = "\033[38;5;46m"
		cyan    = "\033[36m"
	)

	bar := strings.Repeat("=", 68)

	fmt.Printf("%s%s%s\n", saffron, bar, reset)
	fmt.Printf("%s%s%s\n", white, bar, reset)
	fmt.Printf("%s%s%s\n", green, bar, reset)
	fmt.Println("")

	frames := bannerFrames(profile.BannerAnimation)
	if len(frames) == 0 {
		frames = bannerFrames("static")
	}

	frameDuration := 250 * time.Millisecond
	if profile.BannerDurationMS > 0 && len(frames) > 0 {
		total := time.Duration(profile.BannerDurationMS) * time.Millisecond
		frameDuration = total / time.Duration(len(frames)*4)
		if frameDuration < 120*time.Millisecond {
			frameDuration = 120 * time.Millisecond
		}
	}

	if profile.BannerAnimation == "none" || len(frames) == 1 {
		printBannerArt(reset, saffron, white, green, frames[0])
	} else {
		frameCount := int((time.Duration(profile.BannerDurationMS) * time.Millisecond) / frameDuration)
		if frameCount < len(frames) {
			frameCount = len(frames)
		}
		for i := 0; i < frameCount; i++ {
			printBannerArt(reset, saffron, white, green, frames[i%len(frames)])
			if i < frameCount-1 {
				fmt.Print("\033[6A")
				time.Sleep(frameDuration)
			}
		}
	}

	if profile.CompactMode {
		fmt.Printf("\n%s  INDUS v%s%s\n\n", saffron, version, reset)
		return
	}

	fmt.Println("")
	fmt.Printf("%s  Namaste! Welcome to INDUS Terminal v%s%s\n", saffron, version, reset)
	fmt.Printf("%s  Native format: ind <command> [options]%s\n", cyan, reset)
	fmt.Printf("%s  Docs: ind docs | Help: help | Exit: exit%s\n", white, reset)
	fmt.Println("")
}

func bannerFrames(animation string) [][]string {
	static := [][]string{
		{
			"        ▄▄▄      ",
			"       (◉_◉)     ",
			"      /|███|\\    ",
			"       /   \\     ",
		},
	}

	switch animation {
	case "none", "static":
		return static
	case "mascot-wave":
		return [][]string{
			{
				"        ▄▄▄      ",
				"       (◉_◉)ノ   ",
				"      /|███|     ",
				"       /   \\     ",
			},
			{
				"        ▄▄▄      ",
				"       (◉_◉)     ",
				"      /|███|\\    ",
				"       /   \\     ",
			},
		}
	default:
		return static
	}
}

func printBannerArt(reset, saffron, white, green string, bot []string) {
	fmt.Printf("%s  ██╗███╗   ██╗██████╗ ██╗   ██╗███████╗%s%s%s%s\n", saffron, reset, white, bot[0], reset)
	fmt.Printf("%s  ██║████╗  ██║██╔══██╗██║   ██║██╔════╝%s%s%s%s\n", saffron, reset, white, bot[1], reset)
	fmt.Printf("%s  ██║██╔██╗ ██║██║  ██║██║   ██║███████╗%s%s%s%s\n", white, reset, white, bot[2], reset)
	fmt.Printf("%s  ██║██║╚██╗██║██║  ██║██║   ██║╚════██║%s%s%s%s\n", white, reset, white, bot[3], reset)
	fmt.Printf("%s  ██║██║ ╚████║██████╔╝╚██████╔╝███████║%s\n", green, reset)
	fmt.Printf("%s  ╚═╝╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚══════╝%s\n", green, reset)
}

func (t *Terminal) printPrompt() {
	reset := "\033[0m"
	label := t.engine.StateSnapshot().Profile.PromptLabel
	if strings.TrimSpace(label) == "" {
		label = "INDUS"
	}
	fmt.Printf("%s%s%s %s > ", t.session.Theme().Prompt, label, reset, t.session.CWD())
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
