//go:build debug

package debug

import (
	"fmt"
	"os"
	"runtime"
)

func Dump(a ...interface{}) {
	fmt.Println("\n=== Debug Information ===")
	if _, file, line, ok := runtime.Caller(1); ok {
		fmt.Printf("\nCalled from: %s:%d\n", file, line)
	}
	fmt.Println("\n-------------------------")
	fmt.Println(a...)
	fmt.Println("-------------------------")
	fmt.Println("\n=== End of Debug Information ===")
}

func Dd(a ...interface{}) {
	fmt.Println("\n=== Debug Information ===")
	if _, file, line, ok := runtime.Caller(1); ok {
		fmt.Printf("\nCalled from: %s:%d\n", file, line)
	}
	fmt.Println("\n-------------------------")
	fmt.Println(a...)
	fmt.Println("-------------------------")
	fmt.Println("\n=== Program will exit ===")
	os.Exit(1)
}
