package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

type benchmark struct {
	name string
	fn   func(*testing.B)
}

type result struct {
	name    string
	summary testing.BenchmarkResult
}

// runBenchmarks drives each benchmark through testing.Benchmark, which auto
// tunes b.N until the measurement stabilizes. Stdout is swapped for /dev/null
// while each benchmark runs so chatty workloads (accessor logging, console
// shims, etc.) don't pollute the final summary.
func runBenchmarks(benchmarks []benchmark) []result {
	results := make([]result, 0, len(benchmarks))
	for _, bench := range benchmarks {
		fmt.Fprintf(os.Stderr, "running %s...\n", bench.name)
		summary := withSilencedStdout(func() testing.BenchmarkResult {
			return testing.Benchmark(bench.fn)
		})
		results = append(results, result{name: bench.name, summary: summary})
	}
	return results
}

func withSilencedStdout(fn func() testing.BenchmarkResult) testing.BenchmarkResult {
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return fn()
	}
	defer devnull.Close()
	orig := os.Stdout
	os.Stdout = devnull
	defer func() { os.Stdout = orig }()
	return fn()
}

func printMarkdownTable(w io.Writer, results []result) {
	fmt.Fprintln(w, "| Benchmark | ns/op | ops/sec | iterations |")
	fmt.Fprintln(w, "| --- | ---: | ---: | ---: |")
	for _, r := range results {
		// N == 0 means testing.Benchmark bailed before it could measure
		// anything (typically because the benchmark called b.Fatal during
		// setup).
		if r.summary.N == 0 {
			fmt.Fprintf(w, "| %s | FAILED | | |\n", r.name)
			continue
		}
		ns := r.summary.NsPerOp()
		opsPerSec := 1e9 / float64(ns)
		fmt.Fprintf(w, "| %s | %s | %s | %d |\n",
			r.name, formatNs(ns), humanCount(opsPerSec), r.summary.N)
	}
}

func formatNs(ns int64) string {
	switch {
	case ns >= 1_000_000_000:
		return fmt.Sprintf("%.2fs", float64(ns)/1e9)
	case ns >= 1_000_000:
		return fmt.Sprintf("%.2fms", float64(ns)/1e6)
	case ns >= 1_000:
		return fmt.Sprintf("%.2fµs", float64(ns)/1e3)
	default:
		return fmt.Sprintf("%dns", ns)
	}
}

func humanCount(n float64) string {
	switch {
	case n >= 1e9:
		return fmt.Sprintf("%.2fG", n/1e9)
	case n >= 1e6:
		return fmt.Sprintf("%.2fM", n/1e6)
	case n >= 1e3:
		return fmt.Sprintf("%.2fK", n/1e3)
	default:
		return fmt.Sprintf("%.2f", n)
	}
}
