package main

import (
	"fmt"
	"github.com/dop251/goja"
	"reflect"
	"time"
)

const ObjectJS = `function r(a) {
	return a
}`

func gojaObject() {
	start := time.Now()
	vm := goja.New()
	_, err := vm.RunString(ObjectJS)
	if err != nil {
		panic(err)
	}
	r, ok := goja.AssertFunction(vm.Get("r"))
	if !ok {
		panic("Not a function")
	}

	res, err := r(goja.Undefined(), vm.ToValue(map[string]any{
		"hey": "hi",
		"num": 1.0,
	}))
	if err != nil {
		panic(err)
	}

	out := res.ToObject(vm).Export()
	fmt.Println("Goja stringify in", time.Now().Sub(start), "result:", out, "numtype", reflect.TypeOf(out.(map[string]any)["num"]))
}
