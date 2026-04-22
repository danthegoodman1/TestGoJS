package main

import (
	"errors"
	"testing"

	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
)

// The multiply benchmarks exercise three call directions:
//   1. JS -> Go: the JS loop calls a Go callback named `multiply`
//   2. Go -> JS: Go drives a loop that invokes a JS-defined `multiply`
// Each direction gets a goja and otto variant for side-by-side comparison.

const multiplyFromJS = `function loop(N) {
    for (var i = 0; i < N; i++) {
        multiply(5, i)
    }
}`

const multiplyInJS = `function multiply(a, b) {
    return a * b
}`

func benchGojaMultiplyJSToGo(b *testing.B) {
	vm := goja.New()
	if err := vm.Set("multiply", func(call goja.FunctionCall) goja.Value {
		l := call.Arguments[0].ToFloat()
		r := call.Arguments[1].ToFloat()
		return vm.ToValue(l * r)
	}); err != nil {
		b.Fatal(err)
	}
	if _, err := vm.RunString(multiplyFromJS); err != nil {
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

func benchOttoMultiplyJSToGo(b *testing.B) {
	vm := otto.New()
	err := vm.Set("multiply", func(call otto.FunctionCall) otto.Value {
		l, _ := call.Argument(0).ToFloat()
		r, _ := call.Argument(1).ToFloat()
		v, _ := vm.ToValue(l * r)
		return v
	})
	if err != nil {
		b.Fatal(err)
	}
	if _, err := vm.Run(multiplyFromJS); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	if _, err := vm.Call("loop", nil, b.N); err != nil {
		b.Fatal(err)
	}
}

func benchGojaMultiplyGoToJS(b *testing.B) {
	vm := goja.New()
	if _, err := vm.RunString(multiplyInJS); err != nil {
		b.Fatal(err)
	}
	multiply, ok := goja.AssertFunction(vm.Get("multiply"))
	if !ok {
		b.Fatal(errors.New("multiply missing"))
	}
	a := vm.ToValue(40)
	c := vm.ToValue(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := multiply(goja.Undefined(), a, c); err != nil {
			b.Fatal(err)
		}
	}
}

func benchOttoMultiplyGoToJS(b *testing.B) {
	vm := otto.New()
	if _, err := vm.Run(multiplyInJS); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := vm.Call("multiply", nil, 40, 2); err != nil {
			b.Fatal(err)
		}
	}
}
