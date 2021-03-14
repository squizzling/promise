package promise

/*
Provides a basic wrapper around the Promise API (not all functionality is implemented)

Implementation detail: Promise.Then and Promise.Catch are implemented in terms of Promise.ThenOrRejected,
in order to ensure that the underlying js.Func objects are properly released.  There is no provision to
cancel a JSPromise which has had Then, ThenOrRejected, Catch, or Finally called, and it will thus leak
memory if never resolved.
*/
