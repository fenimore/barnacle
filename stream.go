package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Playlist struct {
	Songs []string
}

func (p *Playlist) playlistHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("playlist.html")
	t.Execute(w, p)
}

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

	p := new(Playlist)
	p.Songs = make([]string, len(files))
	for _, f := range files {
		// fmt.Println(f.Name())
		p.Songs = append(p.Songs, f.Name())
	}

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/media", http.StripPrefix("/media", fs))
	http.HandleFunc("/", p.playlistHandler)
	http.ListenAndServe(":5177", nil)
}
