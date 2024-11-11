//go:build !debug

package utils

func Dump(a ...interface{}) {
	// Ignored in production
}

func Dd(a ...interface{}) {
	// Ignored in production
}
