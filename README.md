# barnacle

Stream albums on a local network.

## Instructions
Run the stream executable in the parent directory to a `Music/` directory, or define the path with a flag. barnacle will collect the albums, as organized by directories inside this `Music/` folder. For example:

    barnacle # the executable
	Music/
	    MyAlbum/
		    Track1.mp3
			Track2.mp3
		Genre/
			MyAlbum/
				Track1.mp3
        MyCollection/
		    Disc 1/
			    Track1.mp3
			Disc 2/
				Track1.mp3

## Flag Usage

    -dir string
    	the directory of music, ending in Music/ (default "Music/")
    -note string
        notes to display on index
    -port string
    	the server port, prefixed by : (default ":5177")

For example:

	./barnacle -dir

## Development
I'm using go-bindata to pack the assets (templates) into the binary. Use `go generate` to call `go-bindata -o assets.go templates/` and refresh the html front end. Then `go run` `stream.go` AND `assets.go`.


![screenshot](http://polypmer.github.io/img/barnacle.png "Screeshot")


TODO:
Rewrite history without binary
Release binary on github releases
