package main

import (
	"io"
	"log"
)

func func1(s string) (n int, err error) {
	defer func() {
		log.Printf("func1(%q) = %d, %v", s, n, err)
		// 1. notice the last "()"
	}()
	// 2. compare the difference
	//defer log.Printf("func1(%q) = %d, %v", s, n, err)
	return 7, io.EOF
}

func main() {
	_, _ = func1("Go")
}