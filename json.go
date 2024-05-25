package main

import (
	"fmt"
	"github.com/dop251/goja"
	"time"
)

const JSONJS = `function stringify(a) {
	return JSON.stringify(a)
}`

func gojaJSON() {
	start := time.Now()
	vm := goja.New()
	_, err := vm.RunString(JSONJS)
	if err != nil {
		panic(err)
	}
	stringify, ok := goja.AssertFunction(vm.Get("stringify"))
	if !ok {
		panic("Not a function")
	}

	res, err := stringify(goja.Undefined(), vm.ToValue(map[string]any{
		"hey": "hi",
		"num": 1,
	}))
	if err != nil {
		panic(err)
	}

	fmt.Println("Goja stringify in", time.Now().Sub(start), "result:", res.ToString())
}
