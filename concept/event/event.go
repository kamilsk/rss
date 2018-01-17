package event

import (
	"encoding/json"
	"net/url"
)

type Unpacker interface {
	Unpack() error
}

type Event struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

func (e Event) Unpack() error {
	panic("abstract")
}

type ReadArticleEvent struct {
	Event

	Article *url.URL
}

func (e *ReadArticleEvent) Unpack() error {
	var buf struct {
		RawURL string `json:"article"`
	}
	err := json.Unmarshal(e.Event.Payload, &buf)
	if err != nil {
		return err
	}
	e.Article, err = url.Parse(buf.RawURL)
	return err
}
