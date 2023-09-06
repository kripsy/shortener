//nolint:testpackage
package utils

import (
	"testing"
)

func BenchmarkCreateShortURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CreateShortURL()
		// fmt.Println(res)
	}
}

func BenchmarkCreateShortURLWithoutFmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CreateShortURLWithoutFmt()
		// fmt.Println(res)
	}
}
