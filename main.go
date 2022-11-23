package main

import (
	"encoding/json"
	"syscall/js"
)

func main() {
	done := make(chan struct{})

	js.Global().Set("getDID", js.FuncOf(getDID))

	<-done
}

func getDID(_ js.Value, args []js.Value) interface{} {
	input := args[0].Int()
	result := GetDID(input)
	resultBytes, _ := json.Marshal(result)

	var resultObj map[string]interface{}
	json.Unmarshal(resultBytes, &resultObj)

	return js.ValueOf(resultObj)
}
