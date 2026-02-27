package commands

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"indus/internal/cli"
	"indus/internal/config"
)

type Init struct {
	cfg *config.Config
}

func NewInit(cfg *config.Config) *Init {
	return &Init{cfg: cfg}
}

func (c *Init) Name() string        { return "init" }
func (c *Init) Description() string { return "Initialize a new project structure" }

// Run satisfies cli.Command.
func (c *Init) Run(ctx context.Context, args []string) error {
	return c.RunStream(ctx, args, os.Stdin, os.Stdout)
}

// RunStream satisfies cli.StreamCommand — path output goes to out.
func (c *Init) RunStream(ctx context.Context, args []string, _ io.Reader, out io.Writer) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	name := fs.String("name", "", "Project name (required)")
	dir := fs.String("dir", ".", "Target directory")

	if err := fs.Parse(args); err != nil {
		return &cli.UserError{Msg: "failed to parse flags"}
	}

	if *name == "" {
		fmt.Fprintln(os.Stderr, "Error: --name is required")
		fs.Usage()
		return &cli.UserError{Msg: "missing required flag: --name"}
	}

	projectDir := filepath.Join(*dir, *name)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	dirs := []string{
		projectDir,
		filepath.Join(projectDir, "cmd"),
		filepath.Join(projectDir, "internal"),
		filepath.Join(projectDir, "pkg"),
		filepath.Join(projectDir, "config"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return &cli.InternalError{Msg: "failed to create directory", Err: err}
		}
	}

	readmePath := filepath.Join(projectDir, "README.md")
	readmeContent := fmt.Sprintf("# %s\n\nProject initialized by indus.\n", *name)
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return &cli.InternalError{Msg: "failed to write README", Err: err}
	}

	// Machine-readable output → out (may be a pipe).
	fmt.Fprintf(out, "project_dir=%s\n", projectDir)
	fmt.Fprintf(out, "project_name=%s\n", *name)

	return nil
}
