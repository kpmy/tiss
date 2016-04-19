package tiss

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/kpmy/tiss/gen"
	"github.com/kpmy/tiss/ir"
	"github.com/kpmy/tiss/ir/ops"
	"github.com/kpmy/tiss/ir/types"
)

func TestDump(t *testing.T) {
	if f, err := os.Create("dump.wast"); err == nil {
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
		par.Type(types.I64)
		l0 := &ir.Local{}
		l0.Name("$i")
		l0.Type(types.I64)
		f0.Params = append(f0.Params, par)
		f0.Locals = append(f0.Locals, l0)
		f0.Result = &ir.ResultExpr{Result: types.I64}
		ret := &ir.ReturnExpr{}
		ret.Expr = &ir.ConstExpr{Type: types.I64, Value: 0}
		f0.Code = append(f0.Code, ret)

		fn := &ir.FuncExpr{}
		fn.Name("$start")
		fn.Type = &ir.TypeRef{Type: ir.ThisVariable("$t0")}
		call := &ir.CallExpr{}
		call.Var = ir.ThisVariable("$fib")
		call.Params = []ir.CodeExpr{&ir.ConstExpr{Type: types.I64, Value: 0}}
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

func TestOp(t *testing.T) {
	t.Log(ops.Monadic(types.I32, ops.Clz))
	t.Log(ops.Monadic(types.F32, ops.Nearest))

	t.Log(ops.Dyadic(types.I64, types.I64, ops.Add))
	t.Log(ops.Dyadic(types.I32, types.I32, ops.Ge, true))

	t.Log(ops.Dyadic(types.F32, types.F32, ops.Min))

	t.Log(ops.Conv(types.I64, types.F32, ops.Convert, true))
	t.Log(ops.Conv(types.I64, types.F64, ops.Reinterpret))
}
