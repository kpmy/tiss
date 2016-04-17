package gen //import "github.com/kpmy/tiss/gen"

import (
	"io"

	"github.com/kpmy/tiss/ir"
)

type Writer interface {
	WriteExpr(ir.Expression) error
	WriteValue(interface{}) error
}

type Opts struct {
	PrettyPrint bool
}

var NewWriter func(io.Writer, ...Opts) Writer
