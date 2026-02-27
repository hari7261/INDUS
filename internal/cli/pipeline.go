package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// RunPipeline executes a chain of INDUS commands connected by internal
// io.Pipe streams — no OS shell or subprocesses are involved.
//
// Layout for N segments:
//
//	os.Stdin → [seg0] → pipe(0) → [seg1] → pipe(1) → … → [segN-1] → os.Stdout
//
// Rules:
//   - Every command in the pipeline must implement StreamCommand.
//   - All commands share the same derived context; the first failure
//     cancels the whole chain.
//   - Each command runs in its own goroutine; the function blocks until
//     all goroutines finish.
func (a *App) RunPipeline(ctx context.Context, segments [][]string) error {
	n := len(segments)
	if n == 0 {
		return nil
	}
	// Single segment — no pipe needed, reuse the normal path.
	if n == 1 {
		return a.Run(ctx, segments[0])
	}

	// Validate every command before allocating goroutines.
	cmds := make([]StreamCommand, n)
	for i, seg := range segments {
		if len(seg) == 0 {
			return fmt.Errorf("pipeline segment %d is empty", i)
		}
		raw, ok := a.commands[seg[0]]
		if !ok {
			return fmt.Errorf("unknown command in pipeline: %q", seg[0])
		}
		sc, ok := raw.(StreamCommand)
		if !ok {
			return fmt.Errorf("command %q does not support streaming pipes", seg[0])
		}
		cmds[i] = sc
	}

	// Build N-1 pipes connecting adjacent stages.
	readers := make([]*io.PipeReader, n-1)
	writers := make([]*io.PipeWriter, n-1)
	for i := 0; i < n-1; i++ {
		readers[i], writers[i] = io.Pipe()
	}

	// Shared cancellable context — any stage error cancels the rest.
	pCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		i := i // capture loop variable
		seg := segments[i]

		// Determine this stage's reader and writer.
		var stageIn io.Reader
		var stageOut io.Writer

		if i == 0 {
			stageIn = os.Stdin
		} else {
			stageIn = readers[i-1]
		}
		if i == n-1 {
			stageOut = os.Stdout
		} else {
			stageOut = writers[i]
		}

		wg.Add(1)
		go func(cmd StreamCommand, args []string, in io.Reader, out io.Writer, pipeOut *io.PipeWriter) {
			defer wg.Done()
			// Always close the write end of our pipe when we finish so
			// the next stage receives EOF and doesn't block forever.
			if pipeOut != nil {
				defer pipeOut.Close()
			}

			if err := cmd.RunStream(pCtx, args, in, out); err != nil {
				// Propagate a non-context error.
				if pCtx.Err() == nil {
					errCh <- err
					cancel()
				}
			}
		}(cmds[i], seg[1:], stageIn, stageOut, writerOrNil(writers, i, n))
	}

	wg.Wait()
	close(errCh)

	// Return the first real error (ignore context.Canceled from cancel()).
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

// writerOrNil returns writers[i] when i < n-1 (i.e. we own a write end),
// otherwise nil so the goroutine skips the deferred close on os.Stdout.
func writerOrNil(writers []*io.PipeWriter, i, n int) *io.PipeWriter {
	if i < n-1 {
		return writers[i]
	}
	return nil
}
