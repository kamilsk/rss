package stream

// Based on https://tools.ietf.org/html/rfc4287

// Atom represents Atom 1.0
type Atom interface {
	Author() interface{} // *
	Category() interface{}
	Feed() interface{}
	Rights() interface{}
	Subtitle() interface{}
	Summary() interface{}
	Content() interface{}
	Generator() interface{}
	ID() interface{} // *
	Logo() interface{}
	Entry() interface{}
	Updated() interface{} // *
	Link() interface{}    // *
	Contributor() interface{}
	Published() interface{}
	Title() interface{} // *
}
