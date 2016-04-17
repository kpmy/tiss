package gen //import "github.com/kpmy/tiss/gen"

import (
	"io"

	"github.com/kpmy/tiss/ir"
)

type Writer interface {
	WriteExpr(ir.Expression) error
	WriteValue(interface{}) error
}

var NewWriter func(io.Writer) Writer
