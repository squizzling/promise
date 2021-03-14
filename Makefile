buildtest:
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js internal/static/wasm_exec.js
	GOARCH=wasm GOOS=js go build -o internal/static/tests.wasm ./cmd/tests
	go build ./cmd/testsvr
