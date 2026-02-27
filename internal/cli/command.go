package cli

import (
	"context"
	"io"
)

// Command is the base interface every INDUS command must satisfy.
type Command interface {
	Name() string
	Description() string
	Run(ctx context.Context, args []string) error
}

// StreamCommand extends Command so a command can participate in an
// internal pipeline.  RunStream must write its output exclusively to
// out and must not touch os.Stdout.  It may read from in (e.g. for
// POST body or grep-style filtering).  os.Stderr is still available
// for diagnostic / progress messages that must not pollute the stream.
type StreamCommand interface {
	Command
	RunStream(ctx context.Context, args []string, in io.Reader, out io.Writer) error
}

