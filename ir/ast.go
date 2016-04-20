package ir //import "github.com/kpmy/tiss/ir"

import (
	"reflect"

	"github.com/kpmy/tiss/ir/types"
	"github.com/kpmy/ypk/fn"
	. "github.com/kpmy/ypk/tc"
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

type snamed struct {
	name string
}

func (n *snamed) Name(s ...string) string {
	if len(s) > 0 {
		Assert(s[0][0] == '$' || s[0] == "", 20)
		n.name = s[0]
	}
	return n.name
}

type named struct {
	name string
}

func (n *named) Name(s ...string) string {
	if len(s) > 0 {
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
	snamed
	Type   *TypeRef
	Params []*Param
	Locals []*Local
	Result *ResultExpr
	Code   []CodeExpr
}

func (f *FuncExpr) Validate() (err error) {
	return
}

type megaParam struct {
	ns    `sexpr:"param"`
	Types []types.Type
}

func (*megaParam) Validate() error { return nil }
func (m *megaParam) Children() (ret []interface{}) {
	for _, t := range m.Types {
		ret = append(ret, t)
	}
	return
}

type megaLocal struct {
	ns    `sexpr:"local"`
	Types []types.Type
}

func (*megaLocal) Validate() error { return nil }
func (m *megaLocal) Children() (ret []interface{}) {
	for _, t := range m.Types {
		ret = append(ret, t)
	}
	return
}

func (f *FuncExpr) Children() (ret []interface{}) {
	if f.name != "" {
		ret = append(ret, f.name)
	}

	if !fn.IsNil(f.Type) {
		ret = append(ret, f.Type)
	}

	{
		var tmp []*Param
		merge := true
		mp := &megaParam{}
		for _, p := range f.Params {
			tmp = append(tmp, p)
			mp.Types = append(mp.Types, p.typ)
			if p.name != "" {
				merge = false
			}
		}
		if merge {
			ret = append(ret, mp)
		} else {
			for _, p := range tmp {
				ret = append(ret, p)
			}
		}
	}

	if !fn.IsNil(f.Result) {
		ret = append(ret, f.Result)
	}
	{
		var tmp []*Local
		merge := true
		ml := &megaLocal{}

		for _, l := range f.Locals {
			tmp = append(tmp, l)
			ml.Types = append(ml.Types, l.typ)
			if l.name != "" {
				merge = false
			}
		}
		if merge {
			ret = append(ret, ml)
		} else {
			for _, l := range tmp {
				ret = append(ret, l)
			}
		}
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
	Result types.Type
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
	snamed
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
	snamed
	typ types.Type
}

func (o *object) Type(t ...types.Type) types.Type {
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
	Table *TableDef
	Imp   []*Import
	Exp   []*Export
	Mem   *Memory
}

func (m *Module) Children() (ret []interface{}) {
	if m.Mem != nil {
		ret = append(ret, m.Mem)
	}

	for _, t := range m.Type {
		ret = append(ret, t)
	}

	if m.Table != nil {
		ret = append(ret, m.Table)
	}

	for _, i := range m.Imp {
		ret = append(ret, i)
	}

	for _, f := range m.Func {
		ret = append(ret, f)
	}

	for _, e := range m.Exp {
		ret = append(ret, e)
	}

	ret = append(ret, m.Start)
	return
}

func (m *Module) Validate() (err error) {
	return
}

type TableDef struct {
	ns    `sexpr:"module"`
	Index []Variable
}

func (t *TableDef) Validate() error { return nil }

func (t *TableDef) Children() (ret []interface{}) {
	for _, v := range t.Index {
		ret = append(ret, v)
	}
	return
}

type Import struct {
	ns `sexpr:"import"`
	snamed
	Mod    string
	Func   string
	Params []*Param
	Result *ResultExpr
}

func (i *Import) Validate() error { return nil }

func (i *Import) Children() (ret []interface{}) {
	if i.name != "" {
		ret = append(ret, i.name)
	}

	ret = append(ret, []rune(i.Mod))
	ret = append(ret, []rune(i.Func))

	for _, p := range i.Params {
		ret = append(ret, p)
	}

	if !fn.IsNil(i.Result) {
		ret = append(ret, i.Result)
	}

	return
}

type Export struct {
	ns `sexpr:"export"`
	named
	Var Variable
	Mem bool
}

func (e *Export) Validate() error {
	if e.Mem != e.Var.IsEmpty() {
		return Error("empty export")
	}
	return nil
}

func (e *Export) Children() (ret []interface{}) {

	ret = append(ret, e.name)

	if !e.Var.IsEmpty() {
		ret = append(ret, e.Var)
	} else {
		ret = append(ret, "memory")
	}
	return
}
