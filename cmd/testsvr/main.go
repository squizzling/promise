package main

import (
	"net/http"

	"github.com/squizzling/promise/internal/static"
)

func main() {
	var mux http.ServeMux

	mux.Handle("/", http.FileServer(http.FS(static.Files)))

	if err := http.ListenAndServe(":9997", http.HandlerFunc(mux.ServeHTTP)); err != nil {
		panic(err)
	}
}
