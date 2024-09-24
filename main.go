package main

import (
	"bytes"
	"fmt"
	"github.com/dop251/goja"
	"github.com/google/pprof/profile"
	"github.com/robertkrimen/otto"
	"time"
)

func main() {
	// gojaJSClass()
	profiled()
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

func profiled() {
	js := `function factorial(n) {
    return n === 1 ? n : n * factorial(--n);
}

var i = 0;

while (i++ < 100000) {
    factorial(10);
}`

	b := &bytes.Buffer{}

	err := goja.StartProfile(b)
	if err != nil {
		panic(err)
	}

	gm := goja.New()
	s := time.Now()
	_, err = gm.RunString(js)
	if err != nil {
		panic(err)
	}
	fmt.Println("Goja in", time.Now().Sub(s))

	goja.StopProfile()

	// read the profile
	p, err := profile.Parse(bytes.NewReader(b.Bytes()))
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
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
