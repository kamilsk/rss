package fs

import (
	"encoding/json"
	"errors"
)

// Entity holds information about a source on data access layer.
type Entity struct {
	URL  string          `json:"-"`
	URN  string          `json:"id,omitempty"`
	Name string          `json:"name,omitempty"`
	Raw  json.RawMessage `json:"src"`

	parent   *Entity
	children []*Entity
}

// Feeds returns all feed URLs including feeds of internal entities.
func (e *Entity) Feeds() []string {
	feeds := make([]string, 0, 1+len(e.children))
	if e.URL != "" {
		feeds = append(feeds, e.URL)
	}
	for _, child := range e.children {
		feeds = append(feeds, child.Feeds()...)
	}
	return feeds
}

// URI returns Uniform Resource Identifier of the entity.
// It takes into account its nesting level.
func (e *Entity) URI() string {
	if e.URL != "" {
		return e.URL
	}
	if e.parent != nil {
		return e.parent.URI() + separator + e.URN
	}
	return e.URN
}

// IsValid returns true if internal state of the entity is consistent.
func (e *Entity) IsValid() bool {
	return e.URI() != "" && (e.URL != "" || len(e.children) > 0)
}

// Children returns all internal entities.
func (e *Entity) Children() []*Entity {
	return e.children
}

// AddChild adds a related entity.
func (e *Entity) AddChild(child *Entity) bool {
	if child == nil || !child.IsValid() {
		return false
	}
	for i := range e.children {
		if e.children[i].URI() == child.URI() {
			return false
		}
	}
	child.parent = e
	e.children = append(e.children, child)
	return true
}

// RemoveChild removes entity from internal list.
func (e *Entity) RemoveChild(child *Entity) bool {
	if child == nil {
		return false
	}
	for i := range e.children {
		if e.children[i] == child {
			e.children = append(e.children[:i], e.children[i+1:]...)
			e.children = append(e.children, nil)
			e.children = e.children[:len(e.children)-1]
			child.parent = nil
			return true
		}
	}
	return false
}

// RemoveChildByURI removes entity with the specified URI from internal list.
func (e *Entity) RemoveChildByURI(URI string) bool {
	if URI == "" {
		return false
	}
	for _, child := range e.children {
		if child.URI() == URI {
			return e.RemoveChild(child)
		}
	}
	return false
}

const (
	quote     = 34
	separator = "/"
)

type entity Entity // to prevent recursion

// MarshalJSON implements the json.Marshaler interface.
func (e *Entity) MarshalJSON() ([]byte, error) {
	if !e.IsValid() {
		return nil, errors.New("cannot marshal invalid entity")
	}
	var (
		buf = entity(*e)
		raw []byte
		err error
	)
	if buf.URL != "" {
		raw, err = json.Marshal(buf.URL)
	} else {
		raw, err = json.Marshal(buf.children)
	}
	if err != nil {
		return nil, err
	}
	buf.Raw = raw
	return json.Marshal(buf)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *Entity) UnmarshalJSON(data []byte) error {
	var (
		buf entity
		err error
	)
	err = json.Unmarshal(data, &buf)
	if len(buf.Raw) > 0 {
		if buf.Raw[0] == quote {
			err = json.Unmarshal(buf.Raw, &buf.URL)
		} else {
			err = json.Unmarshal(buf.Raw, &buf.children)
		}
	}
	if err != nil {
		return err
	}
	*e = Entity(buf)
	for _, child := range e.children {
		child.parent = e
	}
	return nil
}
