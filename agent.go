package main

import "syscall/js"

//export add
func add(a, b int) int {
	return a + b
}

//export subtract
func subtract(a, b int) int {
	return a - b
}

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("add", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		if len(p) < 2 || p[0].Type() != js.TypeNumber || p[1].Type() != js.TypeNumber {
			return js.ValueOf("Invalid parameters")
		}
		return add(p[0].Int(), p[1].Int())
	}))
	js.Global().Set("subtract", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		if len(p) < 2 || p[0].Type() != js.TypeNumber || p[1].Type() != js.TypeNumber {
			return js.ValueOf("Invalid parameters")
		}
		return subtract(p[0].Int(), p[1].Int())
	}))
	<-done
}
