package ps

import (
	"io"

	"github.com/kpmy/tiss/ir"
)

var Parse func(io.Reader) (*ir.Module, error)
