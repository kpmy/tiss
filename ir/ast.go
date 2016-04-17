package ir //import "github.com/kpmy/tiss/ir"

import (
	"reflect"

	"github.com/kpmy/ypk/fn"
	. "github.com/kpmy/ypk/tc"
)

type Type string

const (
	Ti32 Type = "i32"
	Ti64 Type = "i64"
	Tf32 Type = "f32"
	Tf64 Type = "f64"
)

type Expression interface {
	Validate() error
}

type CodeExpr interface {
	Eval()
}

type Node interface {
	Children() []interface{}
}

type ns struct{}

type named struct {
	name string
}

func (n *named) Name(s ...string) string {
	if len(s) > 0 {
		Assert(s[0][0] == '$' || s[0] == "", 20)
		n.name = s[0]
	}
	return n.name
}

type Variable struct {
	sv *string
	iv *int
}

func ThisVariable(_x interface{}) (ret Variable) {
	switch x := _x.(type) {
	case string:
		Assert(x[0] == '$', 20)
		ret.sv = &x
	case int:
		ret.iv = &x
	default:
		Halt(100, reflect.TypeOf(x))
	}
	return
}

func (v Variable) IsEmpty() bool {
	return fn.IsNil(v.iv) && fn.IsNil(v.sv)
}

func (v Variable) ValueOf() (ret interface{}) {
	if !fn.IsNil(v.sv) {
		ret = *v.sv
	} else if !fn.IsNil(v.iv) {
		ret = *v.iv
	}
	return
}

type StartExpr struct {
	ns  `sexpr:"start"`
	Var Variable
}

func (s *StartExpr) Validate() (err error) {
	if s.Var.IsEmpty() {
		err = Error("empty start variable")
	}
	return
}

func (s *StartExpr) Children() (ret []interface{}) {
	return []interface{}{s.Var}
}

type FuncExpr struct {
	ns `sexpr:"func"`
	named
	Type   *TypeRef
	Params []*Param
	Locals []*Local
	Result *ResultExpr
	Code   []CodeExpr
}

func (f *FuncExpr) Validate() (err error) {
	return
}

func (f *FuncExpr) Children() (ret []interface{}) {
	if f.name != "" {
		ret = append(ret, f.name)
	}

	if !fn.IsNil(f.Type) {
		ret = append(ret, f.Type)
	}

	for _, p := range f.Params {
		ret = append(ret, p)
	}

	if !fn.IsNil(f.Result) {
		ret = append(ret, f.Result)
	}

	for _, l := range f.Locals {
		ret = append(ret, l)
	}

	for _, c := range f.Code {
		ret = append(ret, c)
	}
	return
}

type TypeRef struct {
	ns   `sexpr:"type"`
	Type Variable
}

func (t *TypeRef) Validate() error {
	if t.Type.IsEmpty() {
		return Error("invalid type ref")
	}
	return nil
}

func (t *TypeRef) Children() (ret []interface{}) {
	return append(ret, t.Type)
}

type ResultExpr struct {
	ns     `sexpr:"result"`
	Result Type
}

func (r *ResultExpr) Validate() error {
	if r.Result == "" {
		return Error("empty result type")
	}
	return nil
}

func (r *ResultExpr) Children() (ret []interface{}) {
	return append(ret, string(r.Result))
}

type TypeDef struct {
	ns `sexpr:"type"`
	named
	Func *FuncExpr
}

func (t *TypeDef) Validate() error {
	if fn.IsNil(t.Func) {
		return Error("typedef func is null")
	}
	if t.Func.Name() != "" {
		return Error("typedef func cannot have name")
	}
	return nil
}

func (t *TypeDef) Children() (ret []interface{}) {
	if t.name != "" {
		ret = append(ret, t.name)
	}

	ret = append(ret, t.Func)
	return
}

type Param struct {
	ns `sexpr:"param"`
	object
}

type Local struct {
	ns `sexpr:"local"`
	object
}

type object struct {
	named
	typ Type
}

func (o *object) Type(t ...Type) Type {
	if len(t) > 0 {
		o.typ = t[0]
	}
	return o.typ
}

func (o *object) Validate() error {
	if o.name == "" {
		return Error("empty object name")
	}

	if o.typ == "" {
		return Error("empty object type")
	}

	return nil
}

func (o *object) Children() (ret []interface{}) {
	return append(ret, o.name, string(o.typ))
}

type Module struct {
	ns    `sexpr:"module"`
	Start *StartExpr
	Func  []*FuncExpr
	Type  []*TypeDef
}

func (m *Module) Children() (ret []interface{}) {
	for _, t := range m.Type {
		ret = append(ret, t)
	}

	for _, f := range m.Func {
		ret = append(ret, f)
	}

	ret = append(ret, m.Start)
	return
}

func (m *Module) Validate() (err error) {
	return
}
