package tests

import (
	"fmt"
	"io/ioutil"
)

func AssertContent(expected string) bool {
	b, _ := ioutil.ReadFile("/tmp/sed-go-output")
	match := string(b) == expected+"\n"
	if !match {
		fmt.Println("expected: \n" + expected)
		fmt.Println("actual: \n" + string(b))
	}
	return match
}
