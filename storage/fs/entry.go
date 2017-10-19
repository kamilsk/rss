package fs

import "encoding/json"

// Entry holds information about the source on data access layer.
type Entry struct {
	URL  string           `json:"-"`
	URN  string           `json:"id"`
	Name string           `json:"name"`
	Raw  *json.RawMessage `json:"src"`

	children []*Entry
}

func (e *Entry) AddChild(child *Entry) {
	e.children = append(e.children, child)
}

func (e Entry) Children() []*Entry {
	return e.children
}

func (e Entry) URI() string {
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
