package tiss

import (
	"os"
	"testing"

	"github.com/kpmy/tiss/gen"
	"github.com/kpmy/tiss/ir"
)

func TestDump(t *testing.T) {
	if f, err := os.Create("dump.wasm"); err == nil {
		defer f.Close()
		m := &ir.Module{}
		m.Start = &ir.StartExpr{Var: ir.ThisVariable("$start")}
		if err = gen.NewWriter(f).WriteExpr(m); err == nil {

		} else {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}
}
