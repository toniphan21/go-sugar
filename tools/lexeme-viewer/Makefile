GOROOT := $(shell go env GOROOT)
WASM_EXEC := $(GOROOT)/lib/wasm/wasm_exec.js

build-wasm:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o asset/resource/main.wasm ./wasm
	du -sh asset/resource/main.wasm
	cp $(WASM_EXEC) asset/resource

run: build-wasm
	go run ./main.go
