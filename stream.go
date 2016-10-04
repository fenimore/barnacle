// Barnacle streams albums on a local network using html5.
// Run Barnacle in a directory with a Music/ directory, inside
// of which should be one or more directories of music albums/playlists.
package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//go:generate go-bindata -o assets.go templates/

/*
   Structs and Constructors
*/

// Collection struct houses all the Albums.
type Collection struct {
	// TODO: make a map?
	Albums    []*Album
	Owner     string
	Host      string
	Index     string
	Playlist  string
	Directory string
	Genres    []*Genre
	Notes     string
}

type Genre struct {
	Title  string
	Albums []*Album
}

// Album struct keeps track of album title, songs
// and song paths.
type Album struct {
	Title string
	Genre string
	Songs []string
	Paths []string
	Cover string
	Count int
}

// NewGenre returns a new Genre with title.
func NewGenre(title string) *Genre {
	return &Genre{Title: title}
}

// NewAlbum returns a new Album with the title.
func NewAlbum(title, genre string) *Album {
	return &Album{Title: title, Genre: genre}
}

func (a *Album) String() string {
	return a.Title
}

func (g *Genre) String() string {
	return g.Title
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
	isSubDir, _ := regexp.MatchString("/", album)
	if !isSubDir {
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
	} else if isSubDir {
		parts := strings.Split(album, "/")
		if strings.HasSuffix(r.URL.Path, "/undefined") {
			return
		}
		for _, g := range c.Genres {
			if g.Title == parts[0] {
				for _, a := range g.Albums {
					if a.Title == parts[1] {
						t := template.New("playlist")
						t, err := t.Parse(c.Playlist)
						if err != nil {
							fmt.Println(err)
						}
						t.Execute(w, a)
						return
					}
				}
			}
		}
	}
	http.NotFound(w, r)
}

func (c *Collection) refreshHandler(w http.ResponseWriter,
	r *http.Request) {
	// TODO: FIXME
	dir := "Music/" // current directory
	c = InitCollection(dir)
	fmt.Println(dir)
	http.Redirect(w, r, "/", 303)
}

/*
   Main Thread
*/

func main() {
	dirFlag := flag.String("dir", "Music/", "the directory of music, ending in Music/")
	portFlag := flag.String("port", ":5177", "the server port, prefixed by :")
	noteFlag := flag.String("note", "", "notes to display on index")
	flag.Parse()

	dir := *dirFlag

	c := InitCollection(dir)
	if *noteFlag != "" {
		c.Notes = *noteFlag
	}
	// Serve Media
	fs := http.FileServer(http.Dir(dir))
	// Handle Routes
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.HandleFunc("/", Logger(c.indexHandler, "Index"))
	http.HandleFunc("/listen/", Logger(c.listenHandler, "Listen"))
	//http.HandleFunc("/refresh/", c.refreshHandler)
	// Print Connection Information
	fmt.Println("Host:    ", c.Host)
	fmt.Println("Ip Addr: ", GetAddress())
	fmt.Println("Port:    ", *portFlag)
	// Listen and Serve on 5177
	// TODO: Flag for port
	err := http.ListenAndServe(*portFlag, nil)
	if err != nil {
		fmt.Println(err)
	}
}

/*
   Functions and Methods for Barnacle
*/

// GetAddress returns the local ip address.
func GetAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}
	addresses, _ := ifaces[2].Addrs()
	address := addresses[0].String() // Trim the /24?
	return address
}

/*
   InitCollection methods:
       InitOwner
       InitAlbums
       InitSongs
       InitHtml
*/

func InitCollection(dir string) *Collection {
	c := new(Collection)
	c.Directory = dir
	c.InitOwner()
	c.CollectAlbums()
	c.CollectGenres()
	//c.InitSongs()
	for _, a := range c.Albums {
		a.CollectSongs(c.Directory)
	}
	c.InitHtml()
	return c
}

// InitOwner sets hostname and username.
func (c *Collection) InitOwner() {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	h, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	c.Host = h
	c.Owner = strings.Title(u.Username)
}

// CollectAlbums finds sub directories, sets albums.
func (c *Collection) CollectAlbums() {
	// TODO: get subdirectories
	c.Albums = make([]*Album, 0)
	c.Genres = make([]*Genre, 0)
	dirs, err := ioutil.ReadDir(c.Directory)
	if err != nil {
		fmt.Println(err)
	}
CheckDirs:
	for _, d := range dirs {
		if d.IsDir() {
			isGenre := true
			album := NewAlbum(d.Name(), "")
			a, err := ioutil.ReadDir(filepath.Join(c.Directory, d.Name()))
			if err != nil {
				fmt.Println(err)
			}
		CheckForGenre:
			for _, s := range a {
				if !s.IsDir() {
					isGenre = false
					break CheckForGenre
				}
			}
			if len(a) < 1 {
				continue CheckDirs
			}
			if isGenre {
				g := NewGenre(d.Name())
				c.Genres = append(c.Genres, g)
			} else {
				c.Albums = append(c.Albums, album)
			}
		}
	}
}

// CollectGenres sets Genres in collection.
// A Genre is a directory without a file.
func (c *Collection) CollectGenres() {
	for _, g := range c.Genres {
		g.Albums = make([]*Album, 0)
		dirs, err := ioutil.ReadDir(filepath.Join(c.Directory, g.Title))
		if err != nil {
			fmt.Println(err)
		}
		for _, d := range dirs {
			isAlbum := false
			album := NewAlbum(d.Name(), g.Title)
			a, err := ioutil.ReadDir(filepath.Join(c.Directory, g.Title, d.Name()))
			if err != nil {
				fmt.Println(err)
			}
		CheckIfAlbum:
			for _, s := range a {
				if !s.IsDir() {
					isAlbum = true
					break CheckIfAlbum
				}
			}
			if isAlbum {
				g.Albums = append(g.Albums, album)
			}
		}
		for _, a := range g.Albums {
			a.CollectSongs(filepath.Join(c.Directory, g.Title))
		}
	}
}

// CollectSongs sets Album songs.
// Song path is /media/AlbumTitle/SongName.
func (a *Album) CollectSongs(dir string) {
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
			if a.Genre == "" {
				path := filepath.Join("/media", a.Title,
					s.Name())
				a.Paths = append(a.Paths, path)
			} else {
				path := filepath.Join("/media", a.Genre,
					a.Title, s.Name())
				a.Paths = append(a.Paths, path)
			}
		} else if isCover {
			if a.Genre == "" {
				a.Cover = filepath.Join("/media", a.Title,
					s.Name())
			} else {
				a.Cover = filepath.Join("/media", a.Genre,
					a.Title, s.Name())
			}
		}
	}
	a.Count = len(a.Songs)
}

// SetUpHtml collects assets and sets Collection templates.
func (c *Collection) InitHtml() {
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
}

// GetGenre look for directories without music files.
// Should this be it's own struct?
func GetGenre() []string {
	var genres []string
	return genres
}

// GetAlbums looks for directories with music files.
func GetAlbums() []string {
	var albums []string
	return albums

}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
