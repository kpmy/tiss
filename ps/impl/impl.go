package impl

import (
	"io"

	"github.com/kpmy/tiss/ir"
	"github.com/kpmy/tiss/ps"
	"github.com/kpmy/tiss/ps/impl/sexpr"
)

func init() {
	ps.Parse = func(r io.Reader) (*ir.Module, error) {
		return sexpr.Parse(r)
	}
}
