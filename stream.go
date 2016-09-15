package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var dir string // to serve
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		dir = "media/" // current directory
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/media", http.StripPrefix("/media", fs))

	http.ListenAndServe(":5177", nil)
}
