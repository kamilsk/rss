package fs

import "encoding/json"

// Entry holds information about the source on data access layer.
type Entry struct {
	URL  string           `json:"-"`
	URN  string           `json:"id,omitempty"`
	Name string           `json:"name,omitempty"`
	Raw  *json.RawMessage `json:"src"`

	parent   *Entry
	children []*Entry
}

func (e *Entry) AddChild(child *Entry) bool {
	if child == nil || child.URI() == "" {
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

func (e *Entry) RemoveChild(child *Entry) bool {
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

func (e *Entry) RemoveChildByURI(URI string) bool {
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

func (e *Entry) Children() []*Entry {
	return e.children
}

func (e *Entry) URI() string {
	if e.URL == "" {
		return e.URN
	}
	return e.URL
}

type entry Entry // to prevent recursion

// MarshalJSON implements json.Marshaler interface.
func (e *Entry) MarshalJSON() ([]byte, error) {
	var (
		v   = entry(*e)
		raw []byte
		err error
	)
	if v.URL != "" {
		raw, err = json.Marshal(v.URL)
	} else {
		raw, err = json.Marshal(v.children)
	}
	if err != nil {
		return nil, err
	}
	v.Raw = &json.RawMessage{}
	*v.Raw = json.RawMessage(raw)
	return json.Marshal(v)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (e *Entry) UnmarshalJSON(data []byte) error {
	var (
		err error
		v   entry
	)
	err = json.Unmarshal(data, &v)
	if v.Raw != nil && len(*v.Raw) > 0 {
		raw := *v.Raw
		// 34 - "
		if raw[0] == 34 {
			err = json.Unmarshal(raw, &v.URL)
		} else {
			err = json.Unmarshal(raw, &v.children)
		}
	}
	if err != nil {
		return err
	}
	*e = Entry(v)
	return nil
}
