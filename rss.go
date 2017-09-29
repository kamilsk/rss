package stream

import (
	"fmt"
	"time"
)

// Based on https://cyber.harvard.edu/rss/rss.html

// RSS represents the RSS 2.0.
type RSS interface {
	// Channel returns RSS channel.
	Channel() Channel
}

// Channel contains information about the channel and its contents.
type Channel interface {
	/* metadata */
	// required
	// Title returns a name of the channel. It's how people refer to your service.
	// If you have an HTML website that contains the same information as your RSS file,
	// the title of your channel should be the same as the title of your website.
	Title() string
	// Link returns a URL to the HTML website corresponding to the channel.
	Link() string
	// Description returns a phrase or sentence describing the channel.
	Description() string

	// optional
	Image() string

	Category() string
	// Copyright returns copyright notice for content in the channel.
	Copyright() string
	Docs() string
	Generator() string
	// Language returns a language the channel is written in.
	// This allows aggregators to group all Italian language sites,
	// for example, on a single page. Should be conform to RFC 5646.
	Language() string

	Rating() fmt.Stringer
	TextInput() string

	PubDate() time.Time       // RFC 822
	LastBuildDate() time.Time // RFC 822
	TTL() int

	Cloud() Cloud
	SkipHours() int
	SkipDays() int

	// ManagingEditor returns an email address for person responsible for editorial content.
	ManagingEditor() fmt.Stringer
	// WebMaster returns an email address for person responsible for technical issues
	// relating to channel.
	WebMaster() fmt.Stringer

	/* content */
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
