//go:build debug

package utils

import "log"

func Dump(a ...interface{}) {
	log.Println(a...)
}

func Dd(a ...interface{}) {
	log.Fatalln(a...)
}
