package fs

import "github.com/kamilsk/stream"

func New(path string) *Storage {
	// check path is writable
	return &Storage{path: path}
}

type Storage struct {
	path string
}

func (s *Storage) Store(src stream.Source) error {
	if _, ok := src.(stream.Aggregator); ok {
		// handle aggregator
	}
	// handle entity
	return nil
}
