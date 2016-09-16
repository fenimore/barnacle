# barnacle

Stream albums on a local network.

## Instructions
Run the stream binary in the parent directory to a `Music/` directory. barnacle will collect the albums, as organized by directories (with no sub directories) inside this `Music/` folder. For example:

    stream
	Music/
	    MyAlbum/
		    Track1.mp3
			Track2.mp3

## Development
I'm using go-bindata to pack the assets (templates) into the binary. Use `go generate` to call `go-bindata -o assets.go templates/` and refresh the html front end. Then `go run` `stream.go` AND `assets.go`.
