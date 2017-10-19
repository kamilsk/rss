package fs

import (
	"encoding/json"
	"errors"
)

// Entry holds information about the source on data access layer.
type Entry struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Raw *json.RawMessage `json:"src"`
	sub []Entry
	url string
}

type entry Entry // to prevent recursion

func (e *Entry) URL() string {
	return e.url
}

func (e *Entry) SubEntries() []Entry {
	return e.sub
}

// MarshalJSON implements json.Marshaler interface.
func (e *Entry) MarshalJSON() ([]byte, error) {
	var (
		v   = entry(*e)
		raw []byte
		err error
	)
	if len(v.sub) == 0 {
		raw, err = json.Marshal(v.url)
	} else {
		raw, err = json.Marshal(v.sub)
	}
	if err != nil {
		return nil, err
	}
	if v.Raw == nil {
		v.Raw = &json.RawMessage{}
	}
	*v.Raw = json.RawMessage(raw)
	return json.Marshal(v)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (e *Entry) UnmarshalJSON(data []byte) error {
	var v entry
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v.Raw == nil {
		return errors.New(`json: "src" node is nil`)
	}
	err := json.Unmarshal(*v.Raw, &v.sub)
	if err != nil {
		err = json.Unmarshal(*v.Raw, &v.url)
	}
	if err != nil {
		return err
	}
	*e = Entry(v)
	return nil
}
