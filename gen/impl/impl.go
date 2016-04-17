package impl

import (
	"io"

	"github.com/kpmy/tiss/gen"
	"github.com/kpmy/tiss/gen/impl/sexpr"
)

func init() {
	gen.NewWriter = func(w io.Writer, o ...gen.Opts) gen.Writer {
		return sexpr.New(w, o...)
	}
}
