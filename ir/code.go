package ir

import (
	"github.com/kpmy/tiss/ir/types"
	"github.com/kpmy/ypk/fn"
	. "github.com/kpmy/ypk/tc"
)

type ReturnExpr struct {
	ns   `sexpr:"return"`
	Expr CodeExpr
}

func (r *ReturnExpr) Validate() error { return nil }
func (r *ReturnExpr) Eval()           {}
func (r *ReturnExpr) Children() (ret []interface{}) {
	return append(ret, r.Expr)
}

type ConstExpr struct {
	Type  types.Type
	Value interface{}
}

func (c *ConstExpr) Name() string {
	return string(c.Type) + ".const"
}

func (c *ConstExpr) Validate() error {
	if c.Type == "" {
		return Error("empty type of const")
	}

	if fn.IsNil(c.Value) {
		return Error("nil const value")
	}

	return nil
}

func (c *ConstExpr) Eval() {}

func (c *ConstExpr) Children() (ret []interface{}) {
	return append(ret, c.Value)
}

type call struct {
	Var    Variable
	Params []CodeExpr
}

func (c *call) Validate() error {
	if c.Var.IsEmpty() {
		return Error("empty call var")
	}
	return nil
}

func (c *call) Eval() {}

func (c *call) Children() (ret []interface{}) {
	ret = append(ret, c.Var)

	for _, p := range c.Params {
		ret = append(ret, p)
	}

	return
}

type CallExpr struct {
	ns `sexpr:"call"`
	call
}

type CallImportExpr struct {
	ns `sexpr:"call_import"`
	call
}

type CallIndirect struct {
	ns `sexpr:"call_indirect"`
	call
	//Var is TypeDef variable
	//Params are params
	Link CodeExpr //Link expr containing adr in table
}

func (c *CallIndirect) Validate() (err error) {
	if err = c.call.Validate(); err == nil {
		if fn.IsNil(c.Link) {
			err = Error("empty link expr of indirect")
		}
	}
	return
}

func (c *CallIndirect) Children() (ret []interface{}) {
	tmp := c.call.Children()
	ret = append(ret, tmp[0])
	ret = append(ret, c.Link)
	for i := 1; i < len(tmp); i++ {
		ret = append(ret, tmp[i])
	}
	return
}

type NopExpr struct {
	ns `sexpr:"nop"`
}

func (n *NopExpr) Validate() error { return nil }

func (n *NopExpr) Eval() {}

type GetLocalExpr struct {
	ns  `sexpr:"get_local"`
	Var Variable
}

func (g *GetLocalExpr) Validate() error {
	if g.Var.IsEmpty() {
		return Error("empty local variable")
	}
	return nil
}

func (g *GetLocalExpr) Children() (ret []interface{}) {
	return append(ret, g.Var)
}

func (*GetLocalExpr) Eval() {}

type SetLocalExpr struct {
	ns   `sexpr:"set_local"`
	Var  Variable
	Expr CodeExpr
}

func (s *SetLocalExpr) Validate() error {
	if s.Var.IsEmpty() {
		return Error("empty local variable")
	}
	if fn.IsNil(s.Expr) {
		return Error("no expr for local varible")
	}
	return nil
}

func (s *SetLocalExpr) Children() (ret []interface{}) {
	return append(ret, s.Var, s.Expr)
}

func (*SetLocalExpr) Eval() {}
