package engine

import (
	"context"
	"testing"
	"time"
)

func TestSysStatsLatencyBudget(t *testing.T) {
	engine := newTestEngine(t)
	session := engine.NewSession(t.TempDir())

	start := time.Now()
	response := engine.ExecuteTokens(context.Background(), session, []string{"sys", "stats"}, ModeExecutable)
	if response.Err != nil {
		t.Fatalf("sys stats failed: %v", response.Err)
	}

	if elapsed := time.Since(start); elapsed > 150*time.Millisecond {
		t.Fatalf("sys stats exceeded latency budget: %s", elapsed)
	}
}

func TestSysStatsAllocations(t *testing.T) {
	engine := newTestEngine(t)
	session := engine.NewSession(t.TempDir())

	allocs := testing.AllocsPerRun(50, func() {
		response := engine.ExecuteTokens(context.Background(), session, []string{"sys", "stats"}, ModeExecutable)
		if response.Err != nil {
			t.Fatalf("sys stats failed: %v", response.Err)
		}
	})

	if allocs > 250 {
		t.Fatalf("expected bounded allocations, got %.2f", allocs)
	}
}

func BenchmarkSysStats(b *testing.B) {
	engine := newTestEngine(b)
	session := engine.NewSession(b.TempDir())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		response := engine.ExecuteTokens(context.Background(), session, []string{"sys", "stats"}, ModeExecutable)
		if response.Err != nil {
			b.Fatalf("sys stats failed: %v", response.Err)
		}
	}
}
