package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"indus/internal/cli"
	"indus/internal/commands"
	"indus/internal/config"
)

var (
	version   = "1.3.0"
	commit    = "initial"
	buildTime = "2026-02-26T12:00:00Z"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	setConsoleTitleW   = kernel32.NewProc("SetConsoleTitleW")
	getConsoleMode     = kernel32.NewProc("GetConsoleMode")
	setConsoleMode     = kernel32.NewProc("SetConsoleMode")
)

func main() {
	// Enable ANSI color support
	enableVirtualTerminalProcessing()
	
	// Set console title
	setConsoleTitle("INDUS Terminal")
	
	// Get current directory
	currentDir, _ := os.Getwd()
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		cfg = config.Default()
	}

	// Create app and register commands
	app := cli.NewApp()
	app.Register(commands.NewInit(cfg))
	app.Register(commands.NewRun(cfg))
	app.Register(commands.NewVersion(version, commit, buildTime))
	app.Register(commands.NewHTTP(cfg))

	// Start terminal
	terminal := NewTerminal(app, version, commit, buildTime, currentDir)
	if err := terminal.Start(ctx); err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "Terminal error: %v\n", err)
		os.Exit(1)
	}
}

// colorSchemes maps single-letter codes to ANSI escape sequences.
// The color is applied to the prompt accent ("INDUS" label and ">").
var colorSchemes = map[string]struct {
	code  string // ANSI escape
	name  string // human-readable label
}{
	"r": {"\033[31m",           "Red"},
	"g": {"\033[32m",           "Green"},
	"b": {"\033[34m",           "Blue"},
	"y": {"\033[33m",           "Yellow"},
	"c": {"\033[36m",           "Cyan"},
	"m": {"\033[35m",           "Magenta"},
	"w": {"\033[97m",           "White"},
	"o": {"\033[38;5;208m",     "Orange"},
	"p": {"\033[38;5;213m",     "Pink"},
	"d": {"\033[36m",           "Default (Cyan)"},
}

type Terminal struct {
	app        *cli.App
	reader     *bufio.Reader
	version    string
	commit     string
	buildTime  string
	currentDir string
	accentColor string // current prompt accent ANSI code
}

func NewTerminal(app *cli.App, version, commit, buildTime, currentDir string) *Terminal {
	return &Terminal{
		app:         app,
		reader:      bufio.NewReader(os.Stdin),
		version:     version,
		commit:      commit,
		buildTime:   buildTime,
		currentDir:  currentDir,
		accentColor: "\033[36m", // default: cyan
	}
}

func (t *Terminal) Start(ctx context.Context) error {
	clearScreen()
	t.printBanner()
	t.printWelcome()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nGoodbye!")
			return nil
		default:
		}

		// Update current directory
		t.currentDir, _ = os.Getwd()
		
		t.printPrompt()

		line, err := t.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nGoodbye!")
				return nil
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if err := t.handleCommand(ctx, line); err != nil {
			if err.Error() == "exit" {
				fmt.Println("Goodbye!")
				return nil
			}
		}
	}
}

func (t *Terminal) handleCommand(ctx context.Context, line string) error {
	// Detect internal pipeline: split by | before any further parsing.
	segments := splitPipeline(line)
	if len(segments) > 1 {
		return t.app.RunPipeline(ctx, segments)
	}

	args := parseCommandLine(line)
	if len(args) == 0 {
		return nil
	}

	cmd := args[0]

	// Built-in terminal commands
	switch cmd {
	case "exit", "quit":
		return fmt.Errorf("exit")
	case "help", "?":
		t.printHelp()
		return nil
	case "clear", "cls":
		clearScreen()
		t.printBanner()
		return nil
	case "cd":
		return t.changeDirectory(args)
	case "pwd":
		fmt.Println(t.currentDir)
		return nil
	case "color":
		return t.handleColor(args)
	case "indus":
		if len(args) == 1 {
			fmt.Println("You're already in INDUS Terminal. Type 'help' for commands.")
			return nil
		}
		return t.app.Run(ctx, args[1:])
	default:
		// Try INDUS command first.
		err := t.app.Run(ctx, args)
		if errors.Is(err, cli.ErrUnknownCommand) {
			// Fall through to the host shell for unrecognised commands.
			return t.runSystemCommand(ctx, args)
		}
		return err
	}
}

func (t *Terminal) handleColor(args []string) error {
	reset := "\033[0m"

	if len(args) < 2 {
		// Print current color and all available options.
		fmt.Printf("Current accent color: %s██%s\n\n", t.accentColor, reset)
		fmt.Println("Available colors:")
		keys := []string{"r", "g", "b", "y", "c", "m", "w", "o", "p", "d"}
		for _, k := range keys {
			s := colorSchemes[k]
			fmt.Printf("  color %s   %s%-10s%s %s██%s\n", k, s.code, s.name, reset, s.code, reset)
		}
		fmt.Println("\nUsage: color <letter>   e.g. color r")
		return nil
	}

	key := strings.ToLower(args[1])
	s, ok := colorSchemes[key]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown color '%s'. Run 'color' to see all options.\n", args[1])
		return nil
	}

	t.accentColor = s.code
	fmt.Printf("Accent color set to %s%s%s\n", s.code, s.name, reset)
	return nil
}

func (t *Terminal) changeDirectory(args []string) error {
	if len(args) < 2 {
		// cd with no args goes to home
		home, _ := os.UserHomeDir()
		return os.Chdir(home)
	}
	
	path := args[1]
	
	// Handle ~ for home directory
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[1:])
	}
	
	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(os.Stderr, "cd: %v\n", err)
		return nil
	}
	
	t.currentDir, _ = os.Getwd()
	return nil
}

func (t *Terminal) runSystemCommand(ctx context.Context, args []string) error {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = t.currentDir
	
	if err := cmd.Run(); err != nil {
		if strings.Contains(err.Error(), "executable file not found") ||
			strings.Contains(err.Error(), "not recognized") {
			fmt.Fprintf(os.Stderr, "'%s' is not recognized as a command.\n", args[0])
			fmt.Fprintln(os.Stderr, "Type 'help' to see INDUS commands.")
			return nil
		}
		return nil // Don't propagate system command errors
	}
	return nil
}

func (t *Terminal) printBanner() {
	saffron := "\033[38;5;208m"
	white := "\033[97m"
	green := "\033[38;5;28m"
	blue := "\033[38;5;33m"
	cyan := "\033[36m"
	yellow := "\033[33m"
	reset := "\033[0m"
	namaste := "\U0001F64F"

	fmt.Println(saffron + "████████████████████████████████████████████████████████████████████████████████████████████████████████████████" + reset)
	fmt.Println(saffron + "████████████████████████████████████████████████████████████████████████████████████████████████████████████████" + reset)
	fmt.Println(white + "████████████████████████████████████████████████████████████████████████████████████████████████████████████████" + reset)
	fmt.Println(white + "██                                                                                                            ██" + reset)
	fmt.Println(white + "██  " + blue + "██╗███╗   ██╗██████╗ ██╗   ██╗███████╗" + white + "    " + cyan + "Terminal v" + t.version + white + "                                            ██" + reset)
	fmt.Println(white + "██  " + blue + "██║████╗  ██║██╔══██╗██║   ██║██╔════╝" + white + "                                                                    ██" + reset)
	fmt.Println(white + "██  " + blue + "██║██╔██╗ ██║██║  ██║██║   ██║███████╗" + white + "    Production-Grade Interactive Terminal                       ██" + reset)
	fmt.Println(white + "██  " + blue + "██║██║╚██╗██║██║  ██║██║   ██║╚════██║" + white + "    for Developers & DevOps                                     ██" + reset)
	fmt.Println(white + "██  " + blue + "██║██║ ╚████║██████╔╝╚██████╔╝███████║" + white + "                                                                    ██" + reset)
	fmt.Println(white + "██  " + blue + "╚═╝╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚══════╝" + white + "                                                                    ██" + reset)
	fmt.Println(white + "██                                                                                                            ██" + reset)
	fmt.Println(white + "████████████████████████████████████████████████████████████████████████████████████████████████████████████████" + reset)
	fmt.Println(green + "████████████████████████████████████████████████████████████████████████████████████████████████████████████████" + reset)
	fmt.Println(green + "████████████████████████████████████████████████████████████████████████████████████████████████████████████████" + reset)
	fmt.Println()
	fmt.Println(yellow + "    " + namaste + "  Namaste! Welcome to INDUS Terminal" + reset)
	fmt.Println()
	fmt.Println(cyan + "    Made with " + saffron + "♥" + reset + cyan + " by hari7261" + reset)
	fmt.Println(blue + "    GitHub: https://github.com/hari7261" + reset)
	fmt.Println()
}

func (t *Terminal) printWelcome() {
	yellow := "\033[33m"
	green := "\033[32m"
	reset := "\033[0m"

	fmt.Printf("%sQuick Start:%s\n", green, reset)
	fmt.Printf("  • Type %shelp%s to see all commands\n", yellow, reset)
	fmt.Printf("  • Use like PowerShell: %scd%s, %sdir%s, %spwd%s, %sipconfig%s\n", yellow, reset, yellow, reset, yellow, reset, yellow, reset)
	fmt.Printf("  • INDUS commands: %sversion%s, %shttp get <url>%s, %sinit%s, %srun%s\n", yellow, reset, yellow, reset, yellow, reset, yellow, reset)
	fmt.Printf("  • Change prompt color: %scolor r%s  %scolor b%s  %scolor g%s  (10 colors)\n", yellow, reset, yellow, reset, yellow, reset)
	fmt.Printf("  • Type %sexit%s to quit\n\n", yellow, reset)
}

func (t *Terminal) printHelp() {
	cyan := "\033[36m"
	yellow := "\033[33m"
	green := "\033[32m"
	reset := "\033[0m"

	fmt.Printf("\n%s╔══════════════════════════════════════════════════════════════════════╗%s\n", cyan, reset)
	fmt.Printf("%s║                      INDUS TERMINAL COMMANDS                         ║%s\n", cyan, reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════════╝%s\n\n", cyan, reset)

	fmt.Printf("%sBUILT-IN SHELL COMMANDS:%s\n", yellow, reset)
	fmt.Printf("  %scd%s <path>          Change directory\n", green, reset)
	fmt.Printf("  %spwd%s                Print working directory\n", green, reset)
	fmt.Printf("  %sclear%s, %scls%s        Clear screen\n", green, reset, green, reset)
	fmt.Printf("  %sexit%s, %squit%s        Exit terminal\n", green, reset, green, reset)
	fmt.Printf("  %shelp%s               Show this help\n", green, reset)
	fmt.Printf("  %scolor%s <letter>     Change prompt accent color\n", green, reset)

	fmt.Printf("\n%sPIPELINES (internal, no OS shell):%s\n", yellow, reset)
	fmt.Printf("  Pipe INDUS commands with %s|%s:\n", green, reset)
	fmt.Printf("  %sversion | http post https://example.com/log%s\n", green, reset)
	fmt.Printf("  %srun --tasks 5 | http post https://example.com/results%s\n", green, reset)
	fmt.Printf("\n%sCOLOR CODES:%s\n", yellow, reset)
	keys := []string{"r", "g", "b", "y", "c", "m", "w", "o", "p", "d"}
	for i, k := range keys {
		s := colorSchemes[k]
		fmt.Printf("  %s%-2s%s= %s%-14s%s", green, k, reset, s.code, s.name, reset)
		if (i+1)%2 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
	
	fmt.Printf("\n%sINDUS COMMANDS:%s\n", yellow, reset)
	fmt.Printf("  %sinit%s --name <project> [--dir <path>]\n", green, reset)
	fmt.Printf("  %srun%s [--workers <n>] [--tasks <n>]\n", green, reset)
	fmt.Printf("  %sversion%s            Show version info\n", green, reset)
	fmt.Printf("  %shttp get%s <url>     Make HTTP GET request\n", green, reset)
	fmt.Printf("  %shttp post%s <url> <data>\n", green, reset)
	
	fmt.Printf("\n%sSYSTEM COMMANDS:%s\n", yellow, reset)
	fmt.Printf("  All Windows commands work: %sipconfig%s, %sping%s, %sdir%s, %sgit%s, %sdocker%s, etc.\n\n", 
		green, reset, green, reset, green, reset, green, reset, green, reset)
}

func (t *Terminal) printPrompt() {
	accent := t.accentColor
	green  := "\033[32m"
	reset  := "\033[0m"

	// Shorten path: replace home dir prefix with ~
	shortPath := t.currentDir
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(shortPath, home) {
		shortPath = "~" + strings.TrimPrefix(shortPath, home)
	}

	fmt.Printf("%sINDUS%s %s%s%s %s>%s ", accent, reset, green, shortPath, reset, accent, reset)
}

// splitPipeline splits a raw command line by | (pipe) characters that
// are not inside single or double quotes.  Each segment is then parsed
// by parseCommandLine.  Returns a single segment when no unquoted | is
// found so callers can treat the single-command case cheaply.
func splitPipeline(line string) [][]string {
	var segments [][]string
	var current strings.Builder
	inQuotes := false
	var quoteChar byte

	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case (c == '"' || c == '\'') && !inQuotes:
			inQuotes = true
			quoteChar = c
			current.WriteByte(c)
		case inQuotes && c == quoteChar:
			inQuotes = false
			current.WriteByte(c)
		case c == '|' && !inQuotes:
			seg := strings.TrimSpace(current.String())
			if seg != "" {
				segments = append(segments, parseCommandLine(seg))
			}
			current.Reset()
		default:
			current.WriteByte(c)
		}
	}
	if seg := strings.TrimSpace(current.String()); seg != "" {
		segments = append(segments, parseCommandLine(seg))
	}
	return segments
}

func parseCommandLine(line string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	
	for i := 0; i < len(line); i++ {
		c := line[i]
		
		switch c {
		case '"', '\'':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				current.WriteByte(c)
			} else if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	
	return args
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func enableVirtualTerminalProcessing() {
	var mode uint32
	handle := syscall.Handle(os.Stdout.Fd())
	
	getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	mode |= 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING
	setConsoleMode.Call(uintptr(handle), uintptr(mode))
}

func setConsoleTitle(title string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	setConsoleTitleW.Call(uintptr(unsafe.Pointer(titlePtr)))
}
