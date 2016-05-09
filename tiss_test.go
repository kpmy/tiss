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
	"github.com/kpmy/tiss/ps"
	"github.com/kpmy/ypk/fn"
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
		fn.Type = &ir.TypeRef{Type: ir.ThisVar("$t0")}
		call := &ir.CallExpr{}
		call.Var = ir.ThisVar("$fib")
		call.Params = []ir.CodeExpr{&ir.ConstExpr{Type: types.I64, Value: 0}}
		fn.Code = append(fn.Code, call)
		m.Func = append(m.Func, f0, fn)
		m.Start = &ir.StartExpr{Var: ir.ThisVar("$start")}

		buf := bytes.NewBuffer(nil)
		if err = gen.NewWriter(buf, gen.Opts{PrettyPrint: false}).WriteExpr(m); err == nil {
			t.Log(buf.String())
			io.Copy(f, buf)
		} else {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}
}

func TestRead(t *testing.T) {
	if f, err := os.Open("dump.wast"); err == nil {
		if m, err := ps.Parse(f); err == nil {
			poo(t, m)
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

func TestExpr(t *testing.T) {
	{
		expr := &ir.LoadExpr{}
		expr.Size = ir.Load32
		expr.Offset = 14
		expr.Align = 32
		expr.Type = types.I64
		expr.Expr = &ir.ConstExpr{Type: types.I64, Value: 34}
		poo(t, expr)
	}
	{
		expr := &ir.StoreExpr{}
		expr.Size = ir.Load32
		expr.Offset = 14
		expr.Align = 32
		expr.Type = types.I64
		expr.Expr = &ir.ConstExpr{Type: types.I64, Value: 34}
		expr.Value = &ir.ConstExpr{Type: types.I64, Value: 124}
		poo(t, expr)
	}
	{
		imp := &ir.Import{}
		imp.Name("$imp")
		imp.Mod = "mod"
		imp.Func = "func"
		poo(t, imp)
	}
	{
		c := &ir.ConstExpr{}
		c.Type = types.F64
		c.Value = 0.05
		poo(t, c)
	}
}

func TestValidation(t *testing.T) {
	b := &ir.Block{}
	poo(t, b)

	l0 := &ir.Local{}
	l0.Name("$i")
	l0.Type(types.I64)
	poo(t, l0)

	s := &ir.Segment{}
	s.Offset = 0
	s.Data = append(s.Data, "f23fcaeaef2f42f34")
	poo(t, s)

	m := &ir.Module{}
	m.Binary = []byte{0, 1, 2, 3}
	poo(t, m)
	/* надо проверить все сущности
	   Br
	   BrIf
	   BrTable
	   CallExpr
	   CallImportExpr
	   CallIndirect
	   CodeExpr
	   ConstExpr
	   ConvertOp
	   CurrentMemoryExpr
	   DyadicOp
	   Export
	   Expression
	   FuncExpr
	   GetLocalExpr
	   GrowMemoryExpr
	   If
	   IfExpr
	   Import
	   LoadExpr
	   Loop
	   Memory
	   Module
	   MonadicOp
	   NopExpr
	   Param
	   ResultExpr
	   ReturnExpr
	   SetLocalExpr
	   StartExpr
	   StoreExpr
	   TableDef
	   TypeDef
	   TypeRef
	   UnreachableExpr
	   Variable
	*/
}

func poo(t *testing.T, e ir.Expression) {
	if fn.IsNil(e) {
		t.Fatal("NIL")
	}
	buf := bytes.NewBufferString("")
	if err := gen.NewWriter(buf).WriteExpr(e); err != nil {
		t.Error(err)
	}
	t.Log(buf.String())
}
