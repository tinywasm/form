//go:build wasm

package main

import (
	"syscall/js"
)

func main() {
	// Your WebAssembly code here

	// Create h1 element
	dom := js.Global().Get("document").Call("createElement", "h1")
	dom.Set("innerHTML", "Hello from WebAssembly!")

	// Get body and append element
	body := js.Global().Get("document").Get("body")
	body.Call("appendChild", dom)

	select {}
}