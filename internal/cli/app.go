package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
)

type App struct {
	commands map[string]Command
	stderr   io.Writer
}

func NewApp() *App {
	return &App{
		commands: make(map[string]Command),
		stderr:   os.Stderr,
	}
}

func (a *App) Register(cmd Command) {
	a.commands[cmd.Name()] = cmd
}

func (a *App) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		a.printHelp()
		return ErrMissingCommand
	}

	cmdName := args[0]

	if cmdName == "help" || cmdName == "--help" || cmdName == "-h" {
		a.printHelp()
		return nil
	}

	cmd, ok := a.commands[cmdName]
	if !ok {
		return ErrUnknownCommand
	}

	return cmd.Run(ctx, args[1:])
}

func (a *App) RunWithOutput(ctx context.Context, args []string) error {
	if len(args) == 0 {
		a.printHelp()
		return ErrMissingCommand
	}

	cmdName := args[0]

	if cmdName == "help" || cmdName == "--help" || cmdName == "-h" {
		a.printHelp()
		return nil
	}

	cmd, ok := a.commands[cmdName]
	if !ok {
		fmt.Fprintf(a.stderr, "Unknown command: %s\n\n", cmdName)
		a.printHelp()
		return ErrUnknownCommand
	}

	return cmd.Run(ctx, args[1:])
}

func (a *App) printHelp() {
	fmt.Fprintln(a.stderr, "indus - Production-grade CLI for API orchestration and developer tooling")
	fmt.Fprintln(a.stderr, "")
	fmt.Fprintln(a.stderr, "Usage:")
	fmt.Fprintln(a.stderr, "  indus <command> [flags]")
	fmt.Fprintln(a.stderr, "")
	fmt.Fprintln(a.stderr, "Available Commands:")

	// Sort names for stable, deterministic output.
	names := make([]string, 0, len(a.commands))
	for name := range a.commands {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		cmd := a.commands[name]
		fmt.Fprintf(a.stderr, "  %-12s %s\n", cmd.Name(), cmd.Description())
	}

	fmt.Fprintln(a.stderr, "")
	fmt.Fprintln(a.stderr, "Use \"indus <command> --help\" for more information about a command.")
}
