package main

import (
	"fmt"
)

func LogInfo(f string, args ...interface{}) {
	msg := fmt.Sprintf(f, args...)
	fmt.Println(msg)
	return
}

func LogErr(f string, args ...interface{}) {
	msg := fmt.Sprintf(f, args...)
	fmt.Println(msg)
}
func main() {
	LogInfo("hello")
	tryA()
}
