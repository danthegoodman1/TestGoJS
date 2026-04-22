package main

import (
	"errors"
	"testing"

	"github.com/dop251/goja"
)

// Thing is exposed to JS via goja's ConstructorCall API so each `new Thing(x)`
// in JS creates a fresh Go-backed instance with an accessor-backed `Val` and a
// `SayBlah` method. The benchmark then exercises constructor + method dispatch
// + accessor get/set in a tight loop.

type Thing struct {
	Val string
}

func (t *Thing) SayBlah() {
	// The original sample printed here; we keep the method body minimal so we
	// measure call dispatch rather than fmt overhead.
	_ = t.Val
}

const thingScript = `function loop(N) {
    for (var i = 0; i < N; i++) {
        var a = new Thing('blah')
        a.SayBlah()
        var b = a.Val
        a.Val = 'hey'
        a.SayBlah()
    }
}`

func benchGojaGoBackedClass(b *testing.B) {
	vm := goja.New()
	err := vm.Set("Thing", func(call goja.ConstructorCall) *goja.Object {
		instance := &Thing{Val: call.Arguments[0].String()}
		obj := call.This
		_ = obj.Set("SayBlah", func(goja.FunctionCall) goja.Value {
			instance.SayBlah()
			return goja.Undefined()
		})
		obj.DefineAccessorProperty("Val",
			vm.ToValue(func(goja.FunctionCall) goja.Value {
				return vm.ToValue(instance.Val)
			}),
			vm.ToValue(func(call goja.FunctionCall) goja.Value {
				instance.Val = call.Argument(0).String()
				return goja.Undefined()
			}),
			goja.FLAG_TRUE, goja.FLAG_TRUE)
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
	if _, err := vm.RunString(thingScript); err != nil {
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
