package commands

import (
	"context"
	"flag"
	"fmt"
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

func (c *Run) Name() string {
	return "run"
}

func (c *Run) Description() string {
	return "Execute a simulated workload with bounded concurrency"
}

func (c *Run) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	
	workers := fs.Int("workers", 4, "Number of concurrent workers")
	tasks := fs.Int("tasks", 20, "Total number of tasks to process")
	
	if err := fs.Parse(args); err != nil {
		return err
	}
	
	fmt.Fprintf(os.Stderr, "Starting run with %d workers processing %d tasks...\n", *workers, *tasks)
	
	// Create job channel
	jobs := make(chan int, *tasks)
	results := make(chan result, *tasks)
	
	// Start worker pool
	var wg sync.WaitGroup
	for w := 1; w <= *workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c.worker(ctx, id, jobs, results)
		}(w)
	}
	
	// Send jobs
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
	
	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Process results
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
	
	// Check if canceled
	if ctx.Err() != nil {
		fmt.Fprintf(os.Stderr, "Run canceled: completed=%d failed=%d\n", completed, failed)
		return ctx.Err()
	}
	
	// Machine-readable output to stdout
	fmt.Printf("completed=%d\n", completed)
	fmt.Printf("failed=%d\n", failed)
	fmt.Printf("total=%d\n", *tasks)
	
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
			
			// Simulate work
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
