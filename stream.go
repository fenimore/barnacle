package main

import (
	"fmt"
	"os"
)

func main() {
	var dir string // to serve
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		dir = "." // current directory
	}
	fmt.Printf("Path: %s\n", dir)
}
