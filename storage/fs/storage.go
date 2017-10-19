package fs

import "github.com/kamilsk/stream"

// Storage is an implementation of stream.Storage interface.
// Uses JSON as a supported format to store a source.
type Storage struct {
}

func (Storage) Store(stream.Source) error {
	panic("implement me")
}

func (Storage) StoreDependencies(stream.Aggregator) error {
	panic("implement me")
}

func (Storage) Remove(stream.Source) error {
	panic("implement me")
}

func (Storage) LoadAll() ([]stream.Source, error) {
	panic("implement me")
}

func (Storage) LoadByID(string) (stream.Source, error) {
	panic("implement me")
}

func (Storage) LoadDependencies(stream.Aggregator) error {
	panic("implement me")
}
