package main

import (
	"syscall/js"
	"time"

	"github.com/squizzling/promise/pkg/promise"
)

type TestFunc func(prefix string)

var tests = map[string]TestFunc{
	"PromiseThen":                 testPromiseThen,
	"PromiseThenAlreadyResolved":  testPromiseThenAlreadyResolved,
	"PromiseCatch":                testPromiseCatch,
	"PromiseCatchAlreadyRejected": testPromiseCatchAlreadyRejected,
	"PromiseFinally":              testPromiseFinally,
}

func setResponse(prefix, message string) {
	js.Global().Get("document").Call("getElementById", prefix+"Test").Set("innerText", message)
}

func setError(prefix, message string) {
	js.Global().Get("document").Call("getElementById", prefix+"Error").Set("innerText", message)
}

func testPromiseThen(prefix string) {
	sigResolveIsDone := newSig()
	sigThenIsAttached := newSig()

	var jsfResolve js.Func
	jsfResolve = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve, _ := args[0], args[1]
		jsfResolve.Release()
		go func() {
			if !sigThenIsAttached.Wait(1 * time.Second) {
				setError(prefix, "timeout waiting for then to be attached")
				return
			}
			resolve.Invoke()
		}()
		return js.Undefined()
	})

	p := promise.Wrap(js.Global().Get("Promise").New(jsfResolve))
	p.Then(func(value js.Value) js.Value {
		sigResolveIsDone.Ready()
		return js.Undefined()
	})

	sigThenIsAttached.Ready()

	if !sigResolveIsDone.Wait(1 * time.Second) {
		setError(prefix, "timeout waiting for resolve to execute")
	} else {
		setResponse(prefix, "Ok")
	}
}

func testPromiseCatch(prefix string) {
	sigResolveIsDone := newSig()
	sigCatchIsAttached := newSig()

	var jsfResolve js.Func
	jsfResolve = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		_, reject := args[0], args[1]
		jsfResolve.Release()
		go func() {
			if !sigCatchIsAttached.Wait(1 * time.Second) {
				setError(prefix, "timeout waiting for catch to be attached")
				return
			}
			reject.Invoke()
		}()
		return js.Undefined()
	})

	p := promise.Wrap(js.Global().Get("Promise").New(jsfResolve))
	p.Catch(func(value js.Value) js.Value {
		sigResolveIsDone.Ready()
		return js.Undefined()
	})

	sigCatchIsAttached.Ready()

	if !sigResolveIsDone.Wait(1 * time.Second) {
		setError(prefix, "timeout waiting for reject to execute")
	} else {
		setResponse(prefix, "Ok")
	}
}

func testPromiseThenAlreadyResolved(prefix string) {
	sigResolveIsDone := newSig()
	sigThenIsInvoked := newSig()

	var jsfResolve js.Func
	jsfResolve = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve, _ := args[0], args[1]
		jsfResolve.Release()
		resolve.Invoke()
		sigResolveIsDone.Ready()
		return js.Undefined()
	})

	p := promise.Wrap(js.Global().Get("Promise").New(jsfResolve))

	if !sigResolveIsDone.Wait(1 * time.Second) {
		setError(prefix, "Timed out")
		return
	}

	p.Then(func(value js.Value) js.Value {
		sigThenIsInvoked.Ready()
		return js.Undefined()
	})

	if !sigThenIsInvoked.Wait(1 * time.Second) {
		setError(prefix, "Timed out")
	} else {
		setResponse(prefix, "Ok")
	}
}

func testPromiseCatchAlreadyRejected(prefix string) {
	sigRejectedIsDone := newSig()
	sigCatchIsInvoked := newSig()

	var jsfReject js.Func
	jsfReject = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		_, reject := args[0], args[1]
		jsfReject.Release()
		reject.Invoke()
		sigRejectedIsDone.Ready()
		return js.Undefined()
	})

	p := promise.Wrap(js.Global().Get("Promise").New(jsfReject))

	if !sigRejectedIsDone.Wait(1 * time.Second) {
		setError(prefix, "Timed out")
		return
	}

	p.Catch(func(value js.Value) js.Value {
		sigCatchIsInvoked.Ready()
		return js.Undefined()
	})

	if !sigCatchIsInvoked.Wait(1 * time.Second) {
		setError(prefix, "Timed out")
	} else {
		setResponse(prefix, "Ok")
	}
}

func testPromiseFinally(prefix string) {
	sigResolveIsDone := newSig()
	sigThenIsAttached := newSig()
	sigFinallyIsDone := newSig()

	var jsfResolve js.Func
	jsfResolve = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve, _ := args[0], args[1]
		jsfResolve.Release()
		go func() {
			if !sigThenIsAttached.Wait(1 * time.Second) {
				setError(prefix, "timeout waiting for then to be attached")
				return
			}
			resolve.Invoke()
		}()
		return js.Undefined()
	})

	p := promise.Wrap(js.Global().Get("Promise").New(jsfResolve))
	p.Then(func(value js.Value) js.Value {
		sigResolveIsDone.Ready()
		return js.Undefined()
	}).Finally(func() js.Value {
		sigFinallyIsDone.Ready()
		return js.Undefined()
	})

	sigThenIsAttached.Ready()

	if !sigResolveIsDone.Wait(1 * time.Second) {
		setError(prefix, "timeout waiting for resolve to execute")
		return
	}

	if !sigFinallyIsDone.Wait(1 * time.Second) {
		setError(prefix, "timeout waiting for finally to execute")
	} else {
		setResponse(prefix, "Ok")
	}
}

func waitForSignal(ch chan struct{}, d time.Duration) bool {
	t := time.NewTimer(d)
	select {
	case <-ch:
		t.Stop()
		return true
	case <-t.C:
		return false
	}
}

type sig chan struct{}

func newSig() sig {
	return make(chan struct{})
}

func (s sig) Wait(d time.Duration) bool {
	t := time.NewTimer(d)
	select {
	case <-s:
		t.Stop()
		return true
	case <-t.C:
		return false
	}
}

func (s sig) Ready() {
	close(s)
}
