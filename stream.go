package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var dir string // to serve
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		dir = "media/" // current directory
	}
	fmt.Printf("Path: %s\n", dir)
	err := filepath.Walk(dir, scanAlbums)
	if err != nil {
		fmt.Println(err)
	}
}

func scanAlbums(path string, f os.FileInfo, err error) error {
	fmt.Printf("Scanned: %s\n", path)
	return nil
}
