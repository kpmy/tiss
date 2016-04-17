package tiss

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/kpmy/tiss/gen"
	"github.com/kpmy/tiss/ir"
)

func TestDump(t *testing.T) {
	if f, err := os.Create("dump.wasm"); err == nil {
		defer f.Close()

		m := &ir.Module{}
		td := &ir.TypeDef{}
		td.Name("$t0")
		td.Func = &ir.FuncExpr{}
		m.Type = append(m.Type, td)

		f0 := &ir.FuncExpr{}
		f0.Name("$fib")
		par := &ir.Param{}
		par.Name("$x")
		par.Type(ir.Ti64)
		l0 := &ir.Local{}
		l0.Name("$i")
		l0.Type(ir.Ti64)
		f0.Params = append(f0.Params, par)
		f0.Locals = append(f0.Locals, l0)
		f0.Result = &ir.ResultExpr{Result: ir.Ti64}
		ret := &ir.ReturnExpr{}
		ret.Expr = &ir.ConstExpr{Type: ir.Ti64, Value: 0}
		f0.Code = append(f0.Code, ret)

		fn := &ir.FuncExpr{}
		fn.Name("$start")
		fn.Type = &ir.TypeRef{Type: ir.ThisVariable("$t0")}
		call := &ir.CallExpr{Var: ir.ThisVariable("$fib"), Params: []ir.CodeExpr{&ir.ConstExpr{Type: ir.Ti64, Value: 0}}}
		fn.Code = append(fn.Code, call)
		m.Func = append(m.Func, f0, fn)
		m.Start = &ir.StartExpr{Var: ir.ThisVariable("$start")}

		buf := bytes.NewBuffer(nil)
		if err = gen.NewWriter(buf, gen.Opts{PrettyPrint: true}).WriteExpr(m); err == nil {
			t.Log(buf.String())
			io.Copy(f, buf)
		} else {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}
}
