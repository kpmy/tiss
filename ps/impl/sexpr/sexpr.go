package sexpr

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/kpmy/tiss/ir"
	"github.com/kpmy/tiss/ir/ops"
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
				case *ir.Memory:
					m.Mem = x
				case *ir.Import:
					m.Imp = append(m.Imp, x)
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
			t.Type = ir.ThisVar(data[0].Value)
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
		s.Var = ir.ThisVar(data[0].Value)

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
			c.Var = ir.ThisVar(data[0].Value)
			start = 1
		}
		for i := start; i < len(data); i++ {
			x := p.list2obj(node2list(data[i]))
			c.Params = append(c.Params, x.(ir.CodeExpr))
		}

		ret = c
	case typ == "call_import":
		c := &ir.CallImportExpr{}
		start := 0
		if len(data) > 0 && data[0].IsScalar() {
			c.Var = ir.ThisVar(data[0].Value)
			start = 1
		}
		for i := start; i < len(data); i++ {
			x := p.list2obj(node2list(data[i]))
			c.Params = append(c.Params, x.(ir.CodeExpr))
		}

		ret = c
	case typ == "memory":
		m := &ir.Memory{}
		start := 0
		if data[0].IsScalar() {
			m.Initial = uint(sexp.Help(data[0]).MustInt())
			start = 1
		} else {
			Halt(100)
		}

		if data[1].IsScalar() {
			m.Max = uint(sexp.Help(data[1]).MustInt())
			start = 2
		} else {
			Halt(100)
		}

		if len(data) > start && data[start].IsList() {
			for i := start; i < len(data); i++ {
				x := p.list2obj(node2list(data[i]))
				m.Segments = append(m.Segments, x.(*ir.Segment))
			}
		}
		ret = m
	case typ == "import":
		i := &ir.Import{}
		start := 0

		if strings.HasPrefix(data[0].Value, "$") { //imp name
			i.Name(data[0].Value)
			start = 1
		}

		Assert(data[1].IsScalar(), 20)
		i.Mod = strings.Trim(data[1].Value, `"`)
		Assert(data[2].IsScalar(), 21)
		i.Func = strings.Trim(data[2].Value, `"`)
		start += 2

		for j := start; j < len(data); j++ {
			_x := p.list2obj(node2list(data[j]))
			switch x := _x.(type) {
			case *ir.Param:
				i.Params = append(i.Params, x)
			default:
				Halt(100, reflect.TypeOf(x))
			}
		}

		ret = i
	case typ == "set_local":
		s := &ir.SetLocalExpr{}
		s.Var = ir.ThisVar(data[0].Value)
		x := p.list2obj(node2list(data[1]))
		s.Expr = x.(ir.CodeExpr)
		ret = s
	case typ == "get_local":
		g := &ir.GetLocalExpr{}
		g.Var = ir.ThisVar(data[0].Value)

		ret = g
	case typ == "loop":
		l := &ir.Loop{}
		start := 0
		if data[start].IsScalar() {
			l.End.Name(data[0].Value)
			start++
		}

		if data[start].IsScalar() {
			l.Start.Name(data[start].Value)
			start++
		}

		for i := start; i < len(data); i++ {
			x := p.list2obj(node2list(data[i]))
			l.Expr = append(l.Expr, x.(ir.CodeExpr))
		}
		ret = l
	case typ == "if":
		i := &ir.If{}
		i.CondExpr = p.list2obj(node2list(data[0])).(ir.CodeExpr)

		t := p.list2obj(node2list(data[1]))
		i.Expr = t.([]ir.CodeExpr)

		if len(data) > 2 {
			e := p.list2obj(node2list(data[2]))
			i.ElseExpr = e.([]ir.CodeExpr)
		}
		ret = i
	case typ == "then" || typ == "else":
		var el []ir.CodeExpr
		for _, x := range data {
			e := p.list2obj(node2list(x)).(ir.CodeExpr)
			el = append(el, e)
		}
		ret = el
	case typ == "br":
		b := &ir.Br{}
		b.Var = ir.ThisVar(data[0].Value)

		if len(data) > 1 {
			x := p.list2obj(node2list(data[1]))
			b.Expr = x.(ir.CodeExpr)
		}

		ret = b
	case typ == "nop":
		ret = &ir.NopExpr{}
	case isOp(typ):
		t, n, s, c := splitOp(typ)
		switch {
		case n == "const":
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
		case strings.HasPrefix(n, "load"):
			l := &ir.LoadExpr{}
			l.Type = t
			size, _ := strconv.Atoi(strings.TrimPrefix(n, "load"))
			l.Size = ir.LoadSize(size)
			if s != nil {
				l.Signed = *s
			}

			x := p.list2obj(node2list(data[0]))
			l.Expr = x.(ir.CodeExpr)
			ret = l
		case strings.HasPrefix(n, "store"):
			s := &ir.StoreExpr{}
			s.Type = t
			size, _ := strconv.Atoi(strings.TrimPrefix(n, "store"))
			s.Size = ir.LoadSize(size)

			x := p.list2obj(node2list(data[0]))
			s.Expr = x.(ir.CodeExpr)
			v := p.list2obj(node2list(data[1]))
			s.Value = v.(ir.CodeExpr)
			ret = s
		case strings.Contains(fmt.Sprint(ops.Add, ops.Eq), n):
			d := &ir.DyadicOp{}
			if s != nil {
				Halt(100)
			} else {
				op := ops.Dyadic(t, t, ops.Op(n))
				d.Op = &op
			}
			l := p.list2obj(node2list(data[0]))
			d.Left = l.(ir.CodeExpr)
			r := p.list2obj(node2list(data[1]))
			d.Right = r.(ir.CodeExpr)

			ret = d
		default:
			Halt(100, t, " ", n, " ", s, " ", c)
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
