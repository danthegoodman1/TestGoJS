# TestGoJS

Micro-benchmarks comparing two embedded JavaScript runtimes for Go:
[`goja`](https://github.com/dop251/goja) and
[`otto`](https://github.com/robertkrimen/otto).

Each workload is exercised through `testing.Benchmark`, which auto-tunes the
iteration count until the measurement stabilizes, so the numbers are per-op
costs (ns and ops/sec) rather than single-run wall clocks.

## Running

```bash
go run .
```

All benchmarks run sequentially and print a markdown results table at the end.
Per-benchmark stdout (accessor logging, `console.log` output, etc.) is
redirected to `/dev/null` during timing so the summary stays clean.

## What each benchmark measures

| Name | Direction | Workload |
| --- | --- | --- |
| factorial | JS hot loop | recursive `factorial(10)` in pure JS |
| JS->Go callback | JS hot loop | JS calls a Go `multiply(a, b)` function |
| Go->JS call | Go hot loop | Go calls a JS `multiply(a, b)` function |
| JSON.stringify | JS hot loop | `JSON.stringify` on a small object |
| object roundtrip | Go hot loop | Go hands a map to JS, JS returns it, Go exports it back |
| console.log shim | JS hot loop | JS calls a Go callback via nested property (`console.log`) |
| Go-backed class | JS hot loop | `new Thing('x')` with accessor get/set backed by a Go struct |
| ES6 class | JS hot loop | `new MyClass()` plus method calls defined in JS |

"JS hot loop" means Go calls a single JS function that loops `b.N` times, so the
Go<->JS boundary is crossed once and we measure what the engine does in pure
JS. "Go hot loop" means Go iterates `b.N` times, invoking a JS function each
iteration, so we measure the per-call boundary cost.

## Results

Environment: Apple M3 Max, macOS 24.6.0 (darwin/arm64), Go 1.26.0,
`goja@2024-06-27`, `otto v0.4.0`.

| Benchmark | ns/op | ops/sec | iterations |
| --- | ---: | ---: | ---: |
| goja: factorial(10) in JS loop | 957ns | 1.04M | 1243298 |
| otto: factorial(10) in JS loop | 9.19µs | 108.79K | 131108 |
| goja: JS->Go callback (multiply) | 96ns | 10.42M | 12404194 |
| otto: JS->Go callback (multiply) | 493ns | 2.03M | 2440303 |
| goja: Go->JS call (multiply) | 92ns | 10.87M | 12880580 |
| otto: Go->JS call (multiply) | 1.40µs | 716.33K | 844569 |
| goja: JSON.stringify small object | 639ns | 1.56M | 1886150 |
| otto: JSON.stringify small object | 1.40µs | 716.33K | 853345 |
| goja: Go->JS->Go object roundtrip | 85ns | 11.76M | 13953481 |
| otto: Go->JS->Go object roundtrip | 1.29µs | 775.80K | 948416 |
| goja: console.log shim (nested Go cb) | 87ns | 11.49M | 13569121 |
| otto: console.log shim (nested Go cb) | 463ns | 2.16M | 2594922 |
| goja: Go-backed class (new Thing + access) | 1.24µs | 809.72K | 977298 |
| goja: ES6 class new + method | 515ns | 1.94M | 2285089 |

Headline: goja is roughly 5-15x faster than otto across the workloads that both
engines can run, with the biggest gap on heavy pure-JS work (factorial) and
Go->JS calls. Two benchmarks are goja-only because otto doesn't support ES6
classes or goja's `ConstructorCall`-based Go struct wiring.
