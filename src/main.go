package main

import (
	"fmt"
	"strings"
)

func main() {
	a := makeAddSuffix(".jpeg")
	b := makeAddSuffix(".png")
	fmt.Println(a("pic1"))
	fmt.Println(b("pic2"))
}

func makeAddSuffix(suffix string) func(fileName string) string {
	return func(fileName string) string {
		if !strings.HasSuffix(fileName, suffix) {
			return fileName + suffix
		}
		return fileName
	}
}
