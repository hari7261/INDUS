package commands

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"indus/internal/config"
)

type Run struct {
	cfg *config.Config
}

func NewRun(cfg *config.Config) *Run {
	return &Run{cfg: cfg}
}

func (c *Run) Name() string        { return "run" }
func (c *Run) Description() string { return "Execute a simulated workload with bounded concurrency" }

// Run satisfies cli.Command.
func (c *Run) Run(ctx context.Context, args []string) error {
	return c.RunStream(ctx, args, os.Stdin, os.Stdout)
}

// RunStream satisfies cli.StreamCommand — final results go to out.
func (c *Run) RunStream(ctx context.Context, args []string, _ io.Reader, out io.Writer) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	workers := fs.Int("workers", 4, "Number of concurrent workers")
	tasks := fs.Int("tasks", 20, "Total number of tasks to process")

	if err := fs.Parse(args); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Starting run with %d workers processing %d tasks...\n", *workers, *tasks)

	jobs := make(chan int, *tasks)
	results := make(chan result, *tasks)

	var wg sync.WaitGroup
	for w := 1; w <= *workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c.worker(ctx, id, jobs, results)
		}(w)
	}

	go func() {
		for j := 1; j <= *tasks; j++ {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case jobs <- j:
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	completed := 0
	failed := 0

	for r := range results {
		if r.err != nil {
			failed++
			fmt.Fprintf(os.Stderr, "Task %d failed: %v\n", r.id, r.err)
		} else {
			completed++
			if completed%5 == 0 {
				fmt.Fprintf(os.Stderr, "Progress: %d/%d tasks completed\n", completed, *tasks)
			}
		}
	}

	if ctx.Err() != nil {
		fmt.Fprintf(os.Stderr, "Run canceled: completed=%d failed=%d\n", completed, failed)
		return ctx.Err()
	}

	// Machine-readable output → out (may be a pipe).
	fmt.Fprintf(out, "completed=%d\n", completed)
	fmt.Fprintf(out, "failed=%d\n", failed)
	fmt.Fprintf(out, "total=%d\n", *tasks)

	return nil
}

type result struct {
	id  int
	err error
}

func (c *Run) worker(ctx context.Context, id int, jobs <-chan int, results chan<- result) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				results <- result{id: job, err: ctx.Err()}
				return
			case <-time.After(100 * time.Millisecond):
				results <- result{id: job, err: nil}
			}
		}
	}
}
