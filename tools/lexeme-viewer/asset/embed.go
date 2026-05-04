package asset

import "embed"

//go:embed resource/index.html resource/main.wasm resource/wasm_exec.js
var Content embed.FS
