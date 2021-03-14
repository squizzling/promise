package main

import (
	"sync"
	"syscall/js"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(len(tests))
	for k, v := range tests {
		makeTestRow(k)
		go func(prefix string, fn TestFunc) {
			fn(prefix)
			wg.Done()
		}(k, v)
	}
	wg.Wait()
}

func makeTestRow(k string) {
	document := js.Global().Get("document")
	jsTests := document.Call("getElementById", "tests")
	tr := document.Call("createElement", "tr")
	jsTests.Call("appendChild", tr)
	tdName := document.Call("createElement", "td")
	tdName.Set("innerText", k)
	tr.Call("appendChild", tdName)
	tdTest := document.Call("createElement", "td")
	tdTest.Set("id", k+"Test")
	tr.Call("appendChild", tdTest)
	tdError := document.Call("createElement", "td")
	tdError.Set("id", k+"Error")
	tr.Call("appendChild", tdError)
}
