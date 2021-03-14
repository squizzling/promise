# What

This is a simple wrapper around the JS Promise API in Go via webassembly

It does not currently implement all the Promise API, only the `Promise` objects and not the static functions for now.

# Why

It simplifies managing `js.Func` values. 

# Tests

Run `make buildtest && ./testsvr`, then open a browser on port 9997.

# Implementation note

Golangs webassembly code size is fairly large, and it's easy to make it worse.  This module favors ugly source code over large compiled code, and always will.
