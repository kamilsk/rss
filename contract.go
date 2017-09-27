package stream

// Aggregator is one entry point for many different sources.
type Aggregator interface {
	Source
	Add(Source)
}

// Source describes a source can be presented in Atom 1.0 and RSS 2.0.
type Source interface {
	Atom() Atom
	RSS() RSS
}
