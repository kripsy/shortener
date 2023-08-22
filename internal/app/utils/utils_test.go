package utils

import (
	"testing"
)

func BenchmarkCreateShortURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CreateShortURL()
	}
}

func BenchmarkCreateShortURLWithoutFmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CreateShortURLWithoutFmt()
	}
}
