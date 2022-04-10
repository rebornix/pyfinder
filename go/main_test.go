package main

import (
	"testing"
)

func BenchmarkGetPython(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getpython()
	}
}
