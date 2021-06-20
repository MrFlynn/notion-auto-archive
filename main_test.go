package main

import "testing"

// Prevents flag library from interrupting tests.
var _ = func() bool {
	testing.Init()
	return true
}()
