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

	// author is adding better method for this https://github.com/dop251/goja/issues/584
	getOwnPropertyNames, ok := goja.AssertFunction(vm.Get("Object").ToObject(vm).Get("getOwnPropertyNames"))
	if !ok {
		panic("failed to assert function")
	}
	array, err := getOwnPropertyNames(nil, vm.Get("MyClass").ToObject(vm).Get("prototype"))
	if err != nil {
		panic(err)
	}
	fmt.Println(array.ToObject(vm).Export()) // print class methods

	array, err = getOwnPropertyNames(nil, myClassInstance)
	if err != nil {
		panic(err)
	}
	fmt.Println(array.ToObject(vm).Export()) // print properties (internal variables)

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
