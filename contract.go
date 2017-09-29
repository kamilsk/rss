package stream

import "io"

// Aggregator is a one entry point for many different sources.
type Aggregator interface {
	Source
	// Add puts a source into an internal source collection.
	Add(src Source)
	// Remove removes a source from an internal source collection.
	Remove(src Source)
	// Get returns a source with a specified identifier or nil if nothing found.
	Get(id string) Source
	// All returns all sources from an internal source collection.
	All() []Source
}

// Source describes an entry can be presented in Atom 1.0 and RSS 2.0.
type Source interface {
	// ID returns a source identifier.
	ID() string
	// Atom returns the Atom 1.0 representation of a source.
	Atom() Atom
	// RSS returns the RSS 2.0 representation of a source.
	RSS() RSS
}

// Storage describes a data access layer behavior.
type Storage interface {
	// Store stores a source into storage.
	Store(Source) error
	// StoreDependencies stores relations of a provided aggregator.
	StoreDependencies(Aggregator) error
	// Remove removes a source from storage.
	Remove(Source) error
	// LoadAll loads all sources from storage.
	LoadAll() ([]Source, error)
	// LoadByID loads a source with a specified identifier or nil if nothing found.
	LoadByID(string) (Source, error)
	// LoadDependencies loads related sources into a provided aggregator.
	LoadDependencies(Aggregator) error
}

// View describes a view layer behavior.
type View interface {
	// Render renders a provided source into a provided output.
	Render(io.Writer, Source) error
}
