package main

import (
	"fmt"
	"os"
)

func main() {
	benchmarks := []benchmark{
		{name: "goja: factorial(10) in JS loop", fn: benchGojaFactorial},
		{name: "otto: factorial(10) in JS loop", fn: benchOttoFactorial},

		{name: "goja: JS->Go callback (multiply)", fn: benchGojaMultiplyJSToGo},
		{name: "otto: JS->Go callback (multiply)", fn: benchOttoMultiplyJSToGo},

		{name: "goja: Go->JS call (multiply)", fn: benchGojaMultiplyGoToJS},
		{name: "otto: Go->JS call (multiply)", fn: benchOttoMultiplyGoToJS},

		{name: "goja: JSON.stringify small object", fn: benchGojaJSONStringify},
		{name: "otto: JSON.stringify small object", fn: benchOttoJSONStringify},

		{name: "goja: Go->JS->Go object roundtrip", fn: benchGojaObjectRoundtrip},
		{name: "otto: Go->JS->Go object roundtrip", fn: benchOttoObjectRoundtrip},

		{name: "goja: console.log shim (nested Go cb)", fn: benchGojaConsoleShim},
		{name: "otto: console.log shim (nested Go cb)", fn: benchOttoConsoleShim},

		{name: "goja: Go-backed class (new Thing + access)", fn: benchGojaGoBackedClass},
		{name: "goja: ES6 class new + method", fn: benchGojaJSClass},
	}

	results := runBenchmarks(benchmarks)

	fmt.Println()
	printMarkdownTable(os.Stdout, results)

	for _, r := range results {
		if r.summary.N == 0 {
			os.Exit(1)
		}
	}
}
