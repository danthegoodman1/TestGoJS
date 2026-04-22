package main

import (
	"errors"
	"testing"

	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
)

// The console.log shim benchmarks exercise nested property lookup into a Go
// callback. Instead of printing, the Go callback just sums the argument into a
// live-referenced counter so the engine can't elide the call.

const consoleScript = `function loop(N) {
    for (var i = 0; i < N; i++) {
        console.log(i)
    }
}`

func benchGojaConsoleShim(b *testing.B) {
	vm := goja.New()
	var total float64
	obj := vm.NewObject()
	if err := obj.Set("log", func(call goja.FunctionCall) goja.Value {
		for _, a := range call.Arguments {
			total += a.ToFloat()
		}
		return goja.Undefined()
	}); err != nil {
		b.Fatal(err)
	}
	if err := vm.Set("console", obj); err != nil {
		b.Fatal(err)
	}
	if _, err := vm.RunString(consoleScript); err != nil {
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
	// Keep the accumulator live so the compiler can't optimize the callback away.
	_ = total
}

func benchOttoConsoleShim(b *testing.B) {
	vm := otto.New()
	var total float64
	obj, err := vm.Object(`({})`)
	if err != nil {
		b.Fatal(err)
	}
	err = obj.Set("log", func(call otto.FunctionCall) otto.Value {
		for _, a := range call.ArgumentList {
			f, _ := a.ToFloat()
			total += f
		}
		return otto.UndefinedValue()
	})
	if err != nil {
		b.Fatal(err)
	}
	if err := vm.Set("console", obj); err != nil {
		b.Fatal(err)
	}
	if _, err := vm.Run(consoleScript); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	if _, err := vm.Call("loop", nil, b.N); err != nil {
		b.Fatal(err)
	}
	_ = total
}
