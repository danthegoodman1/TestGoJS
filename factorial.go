package main

import (
	"errors"
	"testing"

	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
)

// factorialScript exposes a `loop(N)` entry point that recursively computes
// factorial(10) N times. Both engines reach the hot path through the same
// script, so the only thing that varies is the engine.
const factorialScript = `function factorial(n) {
    return n === 1 ? n : n * factorial(--n)
}
function loop(N) {
    for (var i = 0; i < N; i++) {
        factorial(10)
    }
}`

func benchGojaFactorial(b *testing.B) {
	vm := goja.New()
	if _, err := vm.RunString(factorialScript); err != nil {
		b.Fatal(err)
	}
	loop, ok := goja.AssertFunction(vm.Get("loop"))
	if !ok {
		b.Fatal(errors.New("loop missing"))
	}
	n := vm.ToValue(b.N)

	b.ResetTimer()
	if _, err := loop(goja.Undefined(), n); err != nil {
		b.Fatal(err)
	}
}

func benchOttoFactorial(b *testing.B) {
	vm := otto.New()
	if _, err := vm.Run(factorialScript); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	if _, err := vm.Call("loop", nil, b.N); err != nil {
		b.Fatal(err)
	}
}
