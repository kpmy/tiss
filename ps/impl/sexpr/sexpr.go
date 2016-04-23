package sexpr

import (
	"bufio"
	"fmt"
	"io"
	"reflect"

	"github.com/kpmy/tiss/ir"
	"github.com/kpmy/ypk/fn"
	. "github.com/kpmy/ypk/tc"
	"github.com/nsf/sexp"
)

type pr struct {
	root *sexp.Node
	m    *ir.Module
	e    error
}

func (p *pr) err(m ...interface{}) {
	p.e = Error(fmt.Sprint(m...))
	panic(p.e)
}

func node2list(n *sexp.Node) (ret []*sexp.Node) {
	Assert(!fn.IsNil(n) && n.IsList(), 20)
	for x := n.Children; x != nil; x = x.Next {
		ret = append(ret, x)
	}
	return
}

func (p *pr) list2obj(n []*sexp.Node) (ret interface{}) {
	var data []*sexp.Node
	if len(n) > 1 {
		data = n[1:]
	}
	switch typ := n[0].Value; {
	case typ == "module":
		m := &ir.Module{}
		for _, v := range data {
			if v.IsList() {
				_x := p.list2obj(node2list(v))
				switch x := _x.(type) {
				case *ir.TypeDef:
					m.Type = append(m.Type, x)
				case *ir.FuncExpr:
					m.Func = append(m.Func, x)
				case *ir.StartExpr:
					m.Start = x
				default:
					Halt(100, reflect.TypeOf(x))
				}
			} else {
				p.err("unexpected scalar ", v.Value)
			}
		}
		ret = m
	case typ == "type":
		t := &ir.TypeDef{}
		Assert(len(data) > 0, 20)
		var fnode *sexp.Node
		if data[0].IsScalar() {
			t.Name(data[0].Value)
			fnode = data[1]
		} else {
			fnode = data[0]
		}
		f := p.list2obj(node2list(fnode))
		t.Func = f.(*ir.FuncExpr)
		ret = t
	case typ == "func":
		f := &ir.FuncExpr{}

		ret = f
	case typ == "start":
		s := &ir.StartExpr{}
		Assert(len(data) > 0, 20)
		s.Var = ir.ThisVariable(data[0].Value)

		ret = s
	default:
		Halt(100, typ)
	}
	Assert(ret != nil, 60)
	return
}

func (p *pr) mod() {
	m := p.list2obj(node2list(p.root))
	p.m, _ = m.(*ir.Module)
}

func (p *pr) Do() (ret *ir.Module, err error) {
	if p.root.IsList() {
		p.root = p.root.Children
		if p.root.IsList() {
			if mod, _err := p.root.Nth(0); err == nil {
				if mod.IsScalar() && mod.Value == "module" {
					Do(func() {
						p.mod()
					}).Catch(nil, func(e error) {
						p.e = e
					}).Do()
					err = p.e
					ret = p.m
				} else {
					err = Error("not a module")
				}
			} else {
				err = _err
			}
		} else {
			err = Error("not a module")
		}
	} else {
		err = Error("not a module")
	}
	return
}

func Parse(rd io.Reader) (ret *ir.Module, err error) {
	var ctx sexp.SourceContext
	fc := ctx.AddFile("wast", -1)
	p := &pr{}
	if p.root, err = sexp.Parse(bufio.NewReader(rd), fc); err == nil {
		ret, err = p.Do()
	}

	return
}
