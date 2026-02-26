package main

import (
	"bufio"
	"context"
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
	version   = "1.0.0"
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

type Terminal struct {
	app       *cli.App
	reader    *bufio.Reader
	version   string
	commit    string
	buildTime string
	currentDir string
}

func NewTerminal(app *cli.App, version, commit, buildTime, currentDir string) *Terminal {
	return &Terminal{
		app:       app,
		reader:    bufio.NewReader(os.Stdin),
		version:   version,
		commit:    commit,
		buildTime: buildTime,
		currentDir: currentDir,
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
	case "indus":
		if len(args) == 1 {
			fmt.Println("You're already in INDUS Terminal. Type 'help' for commands.")
			return nil
		}
		// Run indus subcommand
		return t.app.Run(ctx, args[1:])
	default:
		// Try INDUS command first
		err := t.app.Run(ctx, args)
		if err != nil && strings.Contains(err.Error(), "nknown command") {
			// Run as system command
			return t.runSystemCommand(ctx, args)
		}
		return err
	}
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
	cyan := "\033[36m"
	green := "\033[32m"
	yellow := "\033[33m"
	reset := "\033[0m"
	
	// Get short path
	shortPath := t.currentDir
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(shortPath, home) {
		shortPath = "~" + strings.TrimPrefix(shortPath, home)
	}
	
	fmt.Printf("%sINDUS%s %s%s%s %s>%s ", cyan, reset, green, shortPath, reset, yellow, reset)
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
