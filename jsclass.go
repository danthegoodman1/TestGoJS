package main

import (
	"errors"
	"testing"

	"github.com/dop251/goja"
)

// JS class benchmark: define a MyClass once and then instantiate + call a
// method in a JS loop. Otto has no equivalent benchmark because it doesn't
// support ES6 class syntax.

const jsClassScript = `
class MyClass {
    constructor() {
        this.i = 0
    }
    DoThing() {
        this.i++
        return this.i
    }
}
function loop(N) {
    for (var i = 0; i < N; i++) {
        var m = new MyClass()
        m.DoThing()
        m.DoThing()
    }
}
`

func benchGojaJSClass(b *testing.B) {
	vm := goja.New()
	if _, err := vm.RunString(jsClassScript); err != nil {
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
