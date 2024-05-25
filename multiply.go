package main

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
	"time"
)

const MultiplyJS = `for (let i = 0; i < 1000; i++) {
	multiply(5, i)
}`

func gojaMultiply() {
	s := time.Now()
	vm := goja.New()
	vm.Set("multiply", func(call goja.FunctionCall) goja.Value {
		l := call.Arguments[0]
		r := call.Arguments[1]
		return vm.ToValue(l.ToFloat() * r.ToFloat())
	})
	val, err := vm.RunString(MultiplyJS)
	if err != nil {
		panic(err)
	}
	val.ExportType()
	fmt.Println("Goja multiply in", time.Now().Sub(s), "result:", val.ToInteger())
}

func ottoMultiply() {
	s := time.Now()
	vm := otto.New()
	vm.Set("multiply", func(call otto.FunctionCall) otto.Value {
		l, _ := call.Argument(0).ToFloat()
		r, _ := call.Argument(1).ToFloat()
		val, _ := vm.ToValue(l * r)
		return val
	})
	val, err := vm.Run(MultiplyJS)
	if err != nil {
		panic(err)
	}
	i, _ := val.ToInteger()
	fmt.Println("Otto multiply in", time.Now().Sub(s), "result:", i)
}
