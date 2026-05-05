//go:build js && wasm

package main

import (
	"syscall/js"
)

func main() {
	js.Global().Set("runExtractLexemes", js.FuncOf(jsRunExtractLexemes))
	select {}
}

func jsRunExtractLexemes(this js.Value, args []js.Value) any {
	result := map[string]any{
		"lexemes": js.Null(),
		"error":   js.Null(),
	}

	if len(args) < 1 {
		result["error"] = "expected 1 argument"
		return js.ValueOf(result)
	}

	lexemes, err := safeExtractLexemes(args[0].String())
	if err != nil {
		result["error"] = err.Error()
	} else {
		result["lexemes"] = lexemes
	}
	return js.ValueOf(result)
}
