// Barnacle streams albums on a local network using html5.
// Run Barnacle in a directory with a Music/ directory, inside
// of which should be one or more directories of music albums/playlists.
package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

//go:generate go-bindata -o assets.go templates/

// Collection struct houses all the Albums.
type Collection struct {
	// TODO: make a map?
	Albums   []*Album
	Owner    string
	Host     string
	Index    string
	Playlist string
}

// Album struct keeps track of album title, songs
// and song paths.
type Album struct {
	Title string
	Songs []string
	Paths []string
	Cover string
	Count int
}

// NewAlbum returns a new album with the title.
func NewAlbum(title string) *Album {
	return &Album{Title: title}
}

func (c *Collection) indexHandler(w http.ResponseWriter,
	r *http.Request) {
	t := template.New("index")
	t, err := t.Parse(c.Index)
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
			t := template.New("playlist")
			t, err := t.Parse(c.Playlist)
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
		dir = os.Args[1] // absolute path to Music/
	} else {
		dir = "Music/" // current directory
	}

	c := new(Collection)
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	h, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	c.Host = h
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
			// The Paths are for the http handlers
			// Not your filesystem
			if !s.IsDir() && !isCover {
				if strings.HasPrefix(s.Name(),
					".") || strings.HasSuffix(s.Name(),
					".aiff") {
					continue
				}
				a.Songs = append(a.Songs, s.Name())
				path := filepath.Join("/media", a.Title,
					s.Name())
				a.Paths = append(a.Paths, path)
			} else if isCover {
				a.Cover = filepath.Join("/media/",
					a.Title, s.Name())
			}
		}
		a.Count = len(a.Songs)
	}

	// Templates from assets
	indexHtml, err := Asset("templates/index.html")
	if err != nil {
		fmt.Println(err)
	}
	playlistHtml, err := Asset("templates/playlist.html")
	if err != nil {
		fmt.Println(err)
	}
	c.Index = string(indexHtml)
	c.Playlist = string(playlistHtml)

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.HandleFunc("/", c.indexHandler)
	http.HandleFunc("/listen/", c.listenHandler)
	fmt.Println("Host:    ", c.Host)
	fmt.Println("Ip Addr: ", getAddress())
	fmt.Println("Port:    ", ":5177")
	err = http.ListenAndServe(":5177", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func getAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}
	addresses, _ := ifaces[2].Addrs()
	address := addresses[0].String() // Trim the /24?
	return address
}
