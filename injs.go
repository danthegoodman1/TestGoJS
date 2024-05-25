package main

import (
	"fmt"
	"github.com/dop251/goja"
	"time"
)

const MultiplyInJS = `function multiply(a, b) {
	return a*b
}`

func gojaMultiplyJS() {
	start := time.Now()
	vm := goja.New()
	_, err := vm.RunString(MultiplyInJS)
	if err != nil {
		panic(err)
	}
	multiply, ok := goja.AssertFunction(vm.Get("multiply"))
	if !ok {
		panic("Not a function")
	}

	res, err := multiply(goja.Undefined(), vm.ToValue(40), vm.ToValue(2))
	if err != nil {
		panic(err)
	}

	fmt.Println("Goja multiply in", time.Now().Sub(start), "result:", res.ToInteger())
}
