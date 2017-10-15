package html

import (
	"io"

	"github.com/kamilsk/stream"
)

// HTMLView is an implementation of stream.View interface.
type HTMLView struct {
}

func (HTMLView) Render(io.Writer, stream.Source) error {
	panic("implement me")
}
