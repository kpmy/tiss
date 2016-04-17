package ir //import "github.com/kpmy/tiss/ir"

import (
	"reflect"

	"github.com/kpmy/ypk/fn"
	. "github.com/kpmy/ypk/tc"
)

type Expression interface {
	Validate() error
}

type Node interface {
	Children() []interface{}
}

type ns struct{}

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

type Module struct {
	ns    `sexpr:"module"`
	Start *StartExpr
}

func (m *Module) Children() (ret []interface{}) {
	ret = append(ret, m.Start)
	return
}

func (m *Module) Validate() (err error) {
	return
}
