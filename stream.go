// Barnacle streams albums on a local network using html5.
// Run Barnacle in a directory with a media/ directory, inside
// of which should be one or more directories of music albums/playlists.
package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Collection struct houses all the Albums.
type Collection struct {
	// TODO: make a map?
	Albums []*Album
	Owner  string
}

// Album struct keeps track of album title, songs
// and song paths.
type Album struct {
	Title string
	Songs []string
	Paths []string
	Cover string
}

// NewAlbum returns a new album with the title.
func NewAlbum(title string) *Album {
	return &Album{Title: title}
}

func (c *Collection) indexHandler(w http.ResponseWriter,
	r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, c)
}

func (c *Collection) listenHandler(w http.ResponseWriter,
	r *http.Request) {
	album := r.URL.Path[8:]

	for _, a := range c.Albums {
		if a.Title == album {
			t, err := template.ParseFiles("playlist.html")
			if err != nil {
				fmt.Println(err)
			}
			t.Execute(w, a)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	var dir string // to serve
	if len(os.Args) > 1 {
		dir = os.Args[1] // absolute path to media/
	} else {
		dir = "media/" // current directory
	}

	c := new(Collection)
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	c.Owner = u.Username
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
			isCover := strings.HasSuffix(s.Name(), ".jpg") || strings.HasSuffix(s.Name(), ".png") || strings.HasSuffix(s.Name(), ".jpeg")
			if !s.IsDir() && !isCover {
				a.Songs = append(a.Songs, s.Name())
				path := filepath.Join("/media", a.Title,
					s.Name())
				a.Paths = append(a.Paths, path)
			} else if isCover {
				a.Cover = filepath.Join("/media/", a.Title, s.Name())
			}
		}
	}
	fmt.Println(c.Albums[0].Cover)
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.HandleFunc("/", c.indexHandler)
	http.HandleFunc("/listen/", c.listenHandler)

	http.ListenAndServe(":5177", nil)
}
