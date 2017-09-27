package stream

import (
	"fmt"
	"time"
)

// Based on https://cyber.harvard.edu/rss/rss.html

// RSS represents RSS 2.0.
type RSS interface {
	Channel() Channel
}

// Channel represents RSS 2.0 channel.
type Channel interface {
	// required
	Title() string
	Link() string
	Description() string

	// optional
	Image() string

	Category() string
	Copyright() string
	Docs() string
	Generator() string
	Language() string // RFC 5646

	Rating() fmt.Stringer
	TextInput() string

	PubDate() time.Time       // RFC 822
	LastBuildDate() time.Time // RFC 822
	TTL() int

	Cloud() Cloud
	SkipHours() int
	SkipDays() int

	ManagingEditor() fmt.Stringer
	WebMaster() fmt.Stringer

	Items() []Item
}

// Image specifies a GIF, JPEG or PNG image that can be displayed with the channel.
type Image interface {
	URL() string
	Title() string
	Link() string
	Width() int
	Height() int
}

// Cloud specifies a web service that supports the rssCloud interface
// which can be implemented in HTTP-POST, XML-RPC or SOAP 1.1.
type Cloud interface {
	Domain() string
	Port() int
	Protocol() string
	RegisterProcedure() string
}

// Item represents element in channel in RSS 2.0.
type Item interface {
	Title() string
	Link() string
	Description() string
	PubDate() time.Time
	GUID() string
}
