package static

import (
	"embed"
)

var (
	//go:embed index.html
	//go:embed tests.wasm
	//go:embed wasm_exec.js
	Files embed.FS
)
