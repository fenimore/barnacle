package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Playlist struct {
	Songs []string
}

type Collection struct {
	Albums []*Album
}

type Album struct {
	Title string
	Songs []string
}

func NewAlbum(title string) *Album {
	return &Album{Title: title}
}

func (p *Playlist) playlistHandler(w http.ResponseWriter,
	r *http.Request) {
	t, _ := template.ParseFiles("playlist.html")
	t.Execute(w, p)
}

func (c *Collection) indexHandler(w http.ResponseWriter,
	r *http.Request) {
	// load list of albums
	var albumList string
	for _, a := range c.Albums {
		albumList += "\n" + a.Title
	}
	fmt.Fprintf(w, "Albums:\n %s", albumList)
}

func main() {
	var dir string // to serve
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		dir = "media/" // current directory
	}

	// Get files in /media directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	p := new(Playlist)
	p.Songs = make([]string, 0)
	for _, f := range files {
		if f.Name() != "" {
			p.Songs = append(p.Songs, f.Name())
		}
	}

	c := new(Collection)
	c.Albums = make([]*Album, 0)
	// Get albums as directories in /media/
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	for _, d := range dirs {
		if d.IsDir() {
			album := NewAlbum(d.Name())
			c.Albums = append(c.Albums, album)
			//albums = append(albums, d.Name())
		}
	}

	for _, a := range c.Albums {
		songs, err := ioutil.ReadDir(filepath.Join(dir, a.Title))
		if err != nil {
			fmt.Println(err)
		}
		for _, s := range songs {
			if !s.IsDir() {
				a.Songs = append(a.Songs, s.Name())

			}
		}
	}

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.HandleFunc("/play", p.playlistHandler)
	//http.handleFunc("/listen", listenHandler)
	http.HandleFunc("/", c.indexHandler)
	http.ListenAndServe(":5177", nil)
}
