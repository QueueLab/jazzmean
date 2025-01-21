package main

import "syscall/js"

//export add
func add(a, b int) int {
	return a + b
}

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("add", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		return add(p[0].Int(), p[1].Int())
	}))
	<-done
}
