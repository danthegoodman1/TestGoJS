package main

import (
	"fmt"
	"github.com/dop251/goja"
	"time"
)

type Thing struct {
	Val string
}

func (t *Thing) SayBlah(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	fmt.Println("SAYING BLAH:", t.Val)
	return goja.Undefined()
}

func gojaInstances() {
	s := time.Now()
	vm := goja.New()
	vm.Set("Thing", func(call goja.ConstructorCall) *goja.Object {
		instance := Thing{
			Val: call.Arguments[0].String(),
		}

		obj := call.This
		obj.Set("SayBlah", func(call goja.FunctionCall) goja.Value {
			return instance.SayBlah(call, vm)
		})

		obj.DefineAccessorProperty("Val", vm.ToValue(func(call goja.FunctionCall) goja.Value {
			// Getter
			fmt.Println("getting")
			return vm.ToValue(instance.Val)
		}), vm.ToValue(func(call goja.FunctionCall) goja.Value {
			// Setter
			fmt.Println("setting")
			instance.Val = call.Argument(0).String()
			return goja.Undefined()
		}), goja.FLAG_TRUE, goja.FLAG_TRUE)

		// return obj
		return nil // will use the obj
	})
	val, err := vm.RunString(`let a = new Thing('blah'); a.SayBlah(); let b = a.Val; a.Val = 'hey'; a.SayBlah()`)
	if err != nil {
		panic(err)
	}
	val.ExportType()
	fmt.Println("Goja multiply in", time.Now().Sub(s), "result:", val.ToInteger())
}
