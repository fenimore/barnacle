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
	// TODO: make a map?
	Albums []*Album
}

type Album struct {
	Title string
	Songs []string
	Paths []string
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

func (c *Collection) listenHandler(w http.ResponseWriter,
	r *http.Request) {
	album := r.URL.Path[1:]
	fmt.Println(album)

	for _, a := range c.Albums {
		if a.Title == album {
			// serve album
			fmt.Println("found it")
			t, _ := template.ParseFiles("playlist.html")
			t.Execute(w, album)
		}
	}
	http.NotFound(w, r)

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

	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	// Get Albums in Collection
	for _, d := range dirs {
		if d.IsDir() {
			album := NewAlbum(d.Name())
			c.Albums = append(c.Albums, album)
		}
	}
	// Get Songs in Albums
	for _, a := range c.Albums {
		songs, err := ioutil.ReadDir(filepath.Join(dir, a.Title))
		if err != nil {
			fmt.Println(err)
		}
		for _, s := range songs {
			if !s.IsDir() {
				a.Songs = append(a.Songs, s.Name())
				path := filepath.Join("media", a.Title,
					s.Name())
				a.Paths = append(a.Paths, path)
			}
		}
	}
	//fmt.Println(c.Albums[0].Paths)
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.HandleFunc("/", c.indexHandler)
	//http.HandleFunc("/play", p.playlistHandler)
	http.HandleFunc("/listen", c.listenHandler)

	http.ListenAndServe(":5177", nil)
}
