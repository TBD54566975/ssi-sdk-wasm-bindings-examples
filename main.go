package main

import (
	"encoding/json"
	"syscall/js"
)

/*
 * This is the glue to bind the functions into javascript so they can be called
 */
func main() {
	done := make(chan struct{})

	// Bind the functions to javascript
	js.Global().Set("sayHello", js.FuncOf(sayHello))
	js.Global().Set("getDID", js.FuncOf(getDID))

	<-done
}

// 1. Simplest function
func sayHello(_ js.Value, args []js.Value) interface{} {
	return "Hello, Wasm World!"
}

// 2. function with arguments. In this case we will just get one integer argument and use it to generate a DID
// eg "alert(getDID(1));" in the browser/javascript
func getDID(_ js.Value, args []js.Value) interface{} {
	input := args[0].Int()
	result := GetDID(input)
	resultBytes, _ := json.Marshal(result)

	var resultObj map[string]interface{}
	json.Unmarshal(resultBytes, &resultObj)

	return js.ValueOf(resultObj)
}
