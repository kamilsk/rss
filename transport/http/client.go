package http

import "github.com/kamilsk/stream"

// Client is an implementation of stream.Client interface.
type Client struct {
}

func (Client) Get(URL string) (stream.Source, error) {
	panic("implement me")
}
