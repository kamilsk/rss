package html

import (
	"io"

	"github.com/kamilsk/stream"
)

// View is an implementation of stream.View interface.
type View struct {
}

func (View) Render(io.Writer, stream.Source) error {
	panic("implement me")
}
