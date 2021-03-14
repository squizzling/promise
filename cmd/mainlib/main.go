package main

import (
	"syscall/js"

	"github.com/squizzling/promise/pkg/promise"
)

func main() {
	// This is not intended to run, it exists purely to call every function in the library, to measure code size.
	promise.Wrap(js.Undefined()).Then(func(value js.Value) js.Value {
		return js.Undefined()
	}).ThenOrRejected(func(value js.Value) js.Value {
		return js.Undefined()
	}, func(reason js.Value) js.Value {
		return js.Undefined()
	}).Catch(func(reason js.Value) js.Value {
		return js.Undefined()
	}).Finally(func() js.Value {
		return js.Undefined()
	})
}
