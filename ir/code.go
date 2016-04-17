package ir

import (
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
	Type  Type
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

type CallExpr struct {
	ns     `sexpr:"call"`
	Var    Variable
	Params []CodeExpr
}

func (c *CallExpr) Validate() error {
	if c.Var.IsEmpty() {
		return Error("empty call var")
	}
	return nil
}

func (c *CallExpr) Eval() {}

func (c *CallExpr) Children() (ret []interface{}) {
	ret = append(ret, c.Var)

	for _, p := range c.Params {
		ret = append(ret, p)
	}

	return
}
