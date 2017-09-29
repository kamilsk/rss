package main

import "github.com/kamilsk/stream"

func main() {
}

// Server is a HTTP server.
type Server struct {
	storage stream.Storage
	view    stream.View
}
