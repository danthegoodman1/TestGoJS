package main

import (
	"errors"
	"testing"

	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
)

// Object roundtrip measures Go -> JS -> Go conversion: we hand a map to the
// VM, the JS function returns it, and we export the result back to Go. Done in
// a Go-driven loop so the boundary crossing is what we're measuring.

const identityScript = `function r(a) { return a }`

func benchGojaObjectRoundtrip(b *testing.B) {
	vm := goja.New()
	if _, err := vm.RunString(identityScript); err != nil {
		b.Fatal(err)
	}
	r, ok := goja.AssertFunction(vm.Get("r"))
	if !ok {
		b.Fatal(errors.New("r missing"))
	}
	payload := vm.ToValue(map[string]any{"hey": "hi", "num": 1.0})
	var sink any

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := r(goja.Undefined(), payload)
		if err != nil {
			b.Fatal(err)
		}
		sink = res.ToObject(vm).Export()
	}
	_ = sink
}

func benchOttoObjectRoundtrip(b *testing.B) {
	vm := otto.New()
	if _, err := vm.Run(identityScript); err != nil {
		b.Fatal(err)
	}
	payload := map[string]any{"hey": "hi", "num": 1.0}
	var sink any

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := vm.Call("r", nil, payload)
		if err != nil {
			b.Fatal(err)
		}
		sink, err = res.Export()
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = sink
}
