package main

import (
	"fmt"
	"github.com/dop251/goja"
	"time"
)

func gojaJSClass() {
	s := time.Now()
	vm := goja.New()
	_, err := vm.RunString(`
		class MyClass {
			constructor() {
				this.i = 0;
			}
			DoThing() {
				this.i++;
				return this.i;
			}
		}
	`)
	if err != nil {
		panic(err)
	}

	// Retrieve the MyClass constructor from the VM
	myClassConstructor := vm.Get("MyClass").ToObject(vm)

	// Create an instance of MyClass
	myClassInstance, err := vm.New(myClassConstructor)
	if err != nil {
		panic(err)
	}

	fmt.Println(myClassInstance.Keys()) // prints "[i]"

	// Call the DoThing method on the instance
	doThing, ok := goja.AssertFunction(myClassInstance.Get("DoThing"))
	if !ok {
		panic("did not get dothing")
	}

	result, err := doThing(myClassInstance)
	if err != nil {
		panic(err)
	}

	fmt.Println("Result of DoThing:", result.ToInteger())

	result, err = doThing(myClassInstance)
	if err != nil {
		panic(err)
	}

	fmt.Println("Result of DoThing (second):", result.ToInteger())

	fmt.Println("Goja execution time:", time.Now().Sub(s))
}
