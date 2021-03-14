package promise

import (
	"syscall/js"
)

type JSPromise struct {
	v js.Value
}

type FulfilledFunc func(value js.Value) js.Value
type RejectedFunc func(reason js.Value) js.Value
type FinallyFunc func() js.Value

// Wrap turns a js.Value in to a JSPromise.  No actual type checking is done.
func Wrap(v js.Value) JSPromise {
	return JSPromise{
		v: v,
	}
}

// Then implements Promise.prototype.then(onFulfilled)
func (jsp JSPromise) Then(fnFulfilled FulfilledFunc) JSPromise {
	return jsp.ThenOrRejected(fnFulfilled, nil)
}

// ThenOrRejected implements Promise.prototype.then(onFulfilled, onRejected)
func (jsp JSPromise) ThenOrRejected(fnFulfilled FulfilledFunc, fnRejected RejectedFunc) JSPromise {
	var jsFulfilled, jsRejected js.Func
	jsFulfilled = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		returnValue := js.Undefined()
		if fnFulfilled != nil {
			returnValue = fnFulfilled(args[0])
		}
		jsFulfilled.Release()
		jsRejected.Release()
		return returnValue
	})
	jsRejected = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		returnValue := js.Undefined()
		if fnRejected != nil {
			returnValue = fnRejected(args[0])
		}
		jsFulfilled.Release()
		jsRejected.Release()
		return returnValue
	})
	return Wrap(jsp.v.Call("then", jsFulfilled, jsRejected))
}

// Catch implements Promise.prototype.catch(onRejected)
func (jsp JSPromise) Catch(fnRejected RejectedFunc) JSPromise {
	return jsp.ThenOrRejected(nil, fnRejected)
}

// Finally implements Promise.prototype.finally(onFinally)
func (jsp JSPromise) Finally(fnFinally FinallyFunc) JSPromise {
	var jsFinally js.Func
	jsFinally = js.FuncOf(func (this js.Value, args[]js.Value) interface{} {
		returnValue := fnFinally()
		jsFinally.Release()
		return returnValue
	})
	return Wrap(jsp.v.Call("finally", jsFinally))
}
