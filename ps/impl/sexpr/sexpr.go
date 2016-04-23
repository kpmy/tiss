package sexpr

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/kpmy/tiss/ir"
	"github.com/kpmy/tiss/ir/types"
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

var typList = []types.Type{types.I32, types.I64, types.F32, types.F64}

func isOp(name string) bool {
	for _, t := range typList {
		if strings.HasPrefix(name, string(t)+".") {
			return true
		}
	}
	return false
}

func splitOp(name string) (t types.Type, n string, s *bool, c types.Type) {
	m := make(map[rune][]rune)
	var rn []rune = []rune(name)
	key := 't'
	for _, r := range rn {
		switch {
		case r == '.':
			key = 'n'
		case r == '_':
			key = 's'
		case r == '/':
			key = 'c'
		default:
			m[key] = append(m[key], r)
		}
	}
	if v, ok := m['t']; ok {
		t = types.Type(v)
	}
	if v, ok := m['s']; ok {
		si := false
		if v[0] == 's' {
			si = true
		} else if v[0] == 'u' {
			si = false
		}
		s = &si
	}
	if v, ok := m['n']; ok {
		n = string(v)
	}
	if v, ok := m['c']; ok {
		c = types.Type(v)
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
		Assert(len(data) > 0, 20)
		if len(data) > 1 || (len(data) == 1 && data[0].IsList()) {
			t := &ir.TypeDef{}
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
		} else if len(data) == 1 && data[0].IsScalar() {
			t := &ir.TypeRef{}
			t.Type = ir.ThisVariable(data[0].Value)
			ret = t
		}
	case typ == "func":
		f := &ir.FuncExpr{}
		if len(data) > 0 {
			start := 0
			if data[0].IsScalar() {
				f.Name(data[0].Value)
				start = 1
			}
			for i := start; i < len(data); i++ {
				_x := p.list2obj(node2list(data[i]))
				switch x := _x.(type) {
				case *ir.Param:
					f.Params = append(f.Params, x)
				case *ir.ResultExpr:
					f.Result = x
				case *ir.Local:
					f.Locals = append(f.Locals, x)
				case ir.CodeExpr:
					f.Code = append(f.Code, x)
				case *ir.TypeRef:
					f.Type = x
				default:
					Halt(100, reflect.TypeOf(x))
				}
			}
		}
		ret = f
	case typ == "start":
		s := &ir.StartExpr{}
		Assert(len(data) > 0, 20)
		s.Var = ir.ThisVariable(data[0].Value)

		ret = s
	case typ == "param":
		p := &ir.Param{}
		if len(data) == 2 {
			Assert(data[0].IsScalar(), 20)
			Assert(data[1].IsScalar(), 21)
			p.Name(data[0].Value)
			p.Type(types.Type(data[1].Value))
		} else if len(data) == 1 {
			Assert(data[0].IsScalar(), 20)
			p.Type(types.Type(data[0].Value))
		}
		ret = p
	case typ == "result":
		r := &ir.ResultExpr{}
		if len(data) > 0 {
			r.Result = types.Type(data[0].Value)
		}

		ret = r
	case typ == "local":
		l := &ir.Local{}
		if len(data) == 2 {
			Assert(data[0].IsScalar(), 20)
			Assert(data[1].IsScalar(), 21)
			l.Name(data[0].Value)
			l.Type(types.Type(data[1].Value))
		} else if len(data) == 1 {
			Assert(data[0].IsScalar(), 20)
			l.Type(types.Type(data[0].Value))
		}
		ret = l
	case typ == "return":
		r := &ir.ReturnExpr{}
		if len(data) > 0 {
			x := p.list2obj(node2list(data[0]))
			r.Expr = x.(ir.CodeExpr)
		}
		ret = r
	case typ == "call":
		c := &ir.CallExpr{}
		start := 0
		if len(data) > 0 && data[0].IsScalar() {
			c.Var = ir.ThisVariable(data[0].Value)
			start = 1
		}
		for i := start; i < len(data); i++ {
			x := p.list2obj(node2list(data[i]))
			c.Params = append(c.Params, x.(ir.CodeExpr))
		}

		ret = c
	case isOp(typ):
		t, n, s, c := splitOp(typ)
		switch n {
		case "const":
			op := &ir.ConstExpr{}
			Assert(len(data) == 1, 20)
			Assert(data[0].IsScalar(), 21)
			op.Type = t
			vn := sexp.Help(data[0])
			if f, err := vn.Float64(); err == nil {
				if i, err := vn.Int(); err == nil {
					op.Value = i
				} else {
					op.Value = f
				}
			}
			ret = op
		default:
			Halt(100, t, n, s, c)
		}
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
					}).Do() //.Catch(nil, func(e error) {
					//p.e = e
					//}).Do()
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
