package commands

import (
	"context"
	"fmt"
	"io"
	"os"
)

type Version struct {
	version   string
	commit    string
	buildTime string
}

func NewVersion(version, commit, buildTime string) *Version {
	return &Version{
		version:   version,
		commit:    commit,
		buildTime: buildTime,
	}
}

func (c *Version) Name() string        { return "version" }
func (c *Version) Description() string { return "Print version information" }

// Run satisfies cli.Command — delegates to RunStream using real stdio.
func (c *Version) Run(ctx context.Context, args []string) error {
	return c.RunStream(ctx, args, os.Stdin, os.Stdout)
}

// RunStream satisfies cli.StreamCommand — writes to out so it can
// participate in an internal pipeline.
func (c *Version) RunStream(_ context.Context, _ []string, _ io.Reader, out io.Writer) error {
	fmt.Fprintf(out, "version=%s\n", c.version)
	fmt.Fprintf(out, "commit=%s\n", c.commit)
	fmt.Fprintf(out, "build_time=%s\n", c.buildTime)
	return nil
}
