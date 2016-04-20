package ir

import (
	. "github.com/kpmy/ypk/tc"
)

type Block struct {
	ns   `sexpr:"block"`
	End  snamed //label in the end
	Expr []CodeExpr
}

func (b *Block) Validate() error { return nil }
func (b *Block) Eval()           {}
func (b *Block) Children() (ret []interface{}) {
	if b.End.name != "" {
		ret = append(ret, b.End.name)
	}
	for _, e := range b.Expr {
		ret = append(ret, e)
	}
	return
}

type Loop struct {
	ns    `sexpr:"loop"`
	Start snamed
	End   snamed
	Expr  []CodeExpr
}

func (l *Loop) Validate() error { return nil }
func (l *Loop) Eval()           {}
func (l *Loop) Children() (ret []interface{}) {
	if l.End.name != "" {
		ret = append(ret, l.End.name)
	}

	if l.Start.name != "" {
		ret = append(ret, l.Start.name)
	}

	for _, e := range l.Expr {
		ret = append(ret, e)
	}
	return
}

type SelectExpr struct {
	ns             `sexpr:"select"`
	Expr, ElseExpr CodeExpr
	CondExpr       CodeExpr
}

func (s *SelectExpr) Validate() error {
	if s.Expr == nil {
		return Error("select expr is nil")
	}

	if s.ElseExpr == nil {
		return Error("select else expr is nil")
	}

	if s.CondExpr == nil {
		return Error("select cond expr is nil")
	}
	return nil
}

func (*SelectExpr) Eval() {}

func (s *SelectExpr) Children() (ret []interface{}) {
	return append(ret, s.Expr, s.ElseExpr, s.CondExpr)
}

type If struct {
	ns             `sexpr:"if"`
	Expr, ElseExpr []CodeExpr
	Name, ElseName snamed
	CondExpr       CodeExpr
}

type thenExpr struct {
	ns `sexpr:"then"`
	snamed
	Expr []CodeExpr
}

type elseExpr struct {
	ns `sexpr:"else"`
	snamed
	Expr []CodeExpr
}

func (t *thenExpr) Validate() error { return nil }
func (t *thenExpr) Children() (ret []interface{}) {
	if t.name != "" {
		ret = append(ret, t.name)
	}
	for _, e := range t.Expr {
		ret = append(ret, e)
	}
	return
}

func (e *elseExpr) Validate() error { return nil }
func (e *elseExpr) Children() (ret []interface{}) {
	if e.name != "" {
		ret = append(ret, e.name)
	}
	for _, e := range e.Expr {
		ret = append(ret, e)
	}
	return
}

func (i *If) Validate() error {
	if i.Expr == nil {
		return Error("if expr is nil")
	}

	if i.ElseExpr == nil {
		return Error("if else expr is nil")
	}

	if i.CondExpr == nil {
		return Error("if cond expr is nil")
	}

	return nil
}

func (*If) Eval() {}

func (s *If) Children() (ret []interface{}) {
	ret = append(ret, s.CondExpr)

	then := &thenExpr{}
	then.snamed = s.Name
	then.Expr = s.Expr

	ret = append(ret, then)

	els := &elseExpr{}
	els.snamed = s.ElseName
	els.Expr = s.ElseExpr

	ret = append(ret, els)
	return
}

type IfExpr struct {
	ns             `sexpr:"if"`
	Expr, ElseExpr CodeExpr
	CondExpr       CodeExpr
}

func (i *IfExpr) Validate() error {
	if i.Expr == nil {
		return Error("if expr is nil")
	}

	if i.CondExpr == nil {
		return Error("if cond expr is nil")
	}
	return nil
}

func (*IfExpr) Eval() {}

func (i *IfExpr) Children() (ret []interface{}) {
	ret = append(ret, i.CondExpr, i.Expr)
	if i.ElseExpr != nil {
		ret = append(ret, i.ElseExpr)
	}
	return
}

type Br struct {
	ns   `sexpr:"br"`
	Var  Variable
	Expr CodeExpr
}

func (b *Br) Validate() error {
	if b.Var.IsEmpty() {
		return Error("br variable empty")
	}
	return nil
}

func (b *Br) Eval() {}

func (b *Br) Children() (ret []interface{}) {
	ret = append(ret, b.Var)

	if b.Expr != nil {
		ret = append(ret, b.Expr)
	}

	return
}

type BrIf struct {
	ns         `sexpr:"br_if"`
	Var        Variable
	Cond, Expr CodeExpr
}

func (b *BrIf) Validate() error {
	if b.Var.IsEmpty() {
		return Error("br_if variable empty")
	}
	return nil
}

func (b *BrIf) Eval() {}

func (b *BrIf) Children() (ret []interface{}) {
	ret = append(ret, b.Var)

	if b.Cond != nil {
		ret = append(ret, b.Cond)
	}
	if b.Expr != nil {
		ret = append(ret, b.Expr)
	}

	return
}

type BrTable struct {
	ns         `sexpr:"br_table"`
	Vars       []Variable
	Default    Variable
	Cond, Expr CodeExpr
}

func (b *BrTable) Validate() error {
	if b.Default.IsEmpty() {
		return Error("br_table variable empty")
	}
	if b.Cond == nil {
		return Error("br_table condition is empty")
	}
	return nil
}

func (b *BrTable) Eval() {}

func (b *BrTable) Children() (ret []interface{}) {
	for _, v := range b.Vars {
		ret = append(ret, v)
	}

	ret = append(ret, b.Default, b.Cond)

	if b.Expr != nil {
		ret = append(ret, b.Expr)
	}
	return
}
