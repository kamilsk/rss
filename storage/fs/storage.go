package fs

import "github.com/kamilsk/stream"

// FileStorage is an implementation of stream.Storage interface.
type FileStorage struct {
}

func (FileStorage) Store(stream.Source) error {
	panic("implement me")
}

func (FileStorage) StoreDependencies(stream.Aggregator) error {
	panic("implement me")
}

func (FileStorage) Remove(stream.Source) error {
	panic("implement me")
}

func (FileStorage) LoadAll() ([]stream.Source, error) {
	panic("implement me")
}

func (FileStorage) LoadByID(string) (stream.Source, error) {
	panic("implement me")
}

func (FileStorage) LoadDependencies(stream.Aggregator) error {
	panic("implement me")
}
