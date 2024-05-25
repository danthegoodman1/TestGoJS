package main

import (
	"fmt"
	"github.com/dop251/goja"
	"time"
)

func gojaObjectFunction() {
	s := time.Now()
	vm := goja.New()
	obj := vm.NewObject()
	obj.Set("log", func(call goja.FunctionCall) goja.Value {
		fmt.Println([]any{"LOG:", func() []any {
			var items []any
			for _, item := range call.Arguments {
				items = append(items, fmt.Sprint(item.Export()))
			}
			return items
		}()}...)
		return goja.Undefined()
	})
	vm.Set("console", obj)
	val, err := vm.RunString(`for (let i = 0; i < 10; i++) {
	console.log(i)
}`)
	if err != nil {
		panic(err)
	}
	val.ExportType()
	fmt.Println("Goja multiply in", time.Now().Sub(s), "result:", val.ToInteger())
}
