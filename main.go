package main

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
	"time"
)

func main() {
	gojaInstances()
}

func creation() {
	js := `function factorial(n) {
    return n === 1 ? n : n * factorial(--n);
}

var i = 0;

while (i++ < 1e6) {
    factorial(10);
}`

	vm := otto.New()
	s := time.Now()
	_, err := vm.Run(js)
	if err != nil {
		panic(err)
	}
	fmt.Println("Otto in", time.Now().Sub(s))

	gm := goja.New()
	s = time.Now()
	_, err = gm.RunString(js)
	if err != nil {
		panic(err)
	}
	fmt.Println("Goja in", time.Now().Sub(s))

}

func factorial() {
	js := `function factorial(n) {
    return n === 1 ? n : n * factorial(--n);
}

var i = 0;

while (i++ < 1e6) {
    factorial(10);
}`

	vm := otto.New()
	s := time.Now()
	_, err := vm.Run(js)
	if err != nil {
		panic(err)
	}
	fmt.Println("Otto in", time.Now().Sub(s))

	gm := goja.New()
	s = time.Now()
	_, err = gm.RunString(js)
	if err != nil {
		panic(err)
	}
	fmt.Println("Goja in", time.Now().Sub(s))

	vm = otto.New()
	s = time.Now()
	_, err = vm.Run(js)
	if err != nil {
		panic(err)
	}
	fmt.Println("Otto in", time.Now().Sub(s))

	gm = goja.New()
	s = time.Now()
	_, err = gm.RunString(js)
	if err != nil {
		panic(err)
	}
	fmt.Println("Goja in", time.Now().Sub(s))
}
