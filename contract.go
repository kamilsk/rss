package stream

// Aggregator is one entry point for many different sources.
type Aggregator interface {
	Source

	Add(src Source)
	Remove(src Source)
	Get(id string)
	All() []Source
}

// Source describes a source can be presented in Atom 1.0 and RSS 2.0.
type Source interface {
	ID() string
	Atom() Atom
	RSS() RSS
}

// Storage describes a data access object behavior.
type Storage interface {
	Store(src Aggregator)
	Remove(src Aggregator)
	Load() []Aggregator
	LoadByID(id string) Aggregator
}
