package main

import (
	"bufio"
	"bytes"
	"container/list"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/kpmy/tiss/gen"
	_ "github.com/kpmy/tiss/gen/impl"
	"github.com/kpmy/tiss/ir"
	"github.com/kpmy/tiss/ir/ops"
	"github.com/kpmy/tiss/ir/types"
	"github.com/kpmy/tiss/ps"
	_ "github.com/kpmy/tiss/ps/impl"
	"github.com/kpmy/ypk/fn"
	. "github.com/kpmy/ypk/tc"
)

type ConsFunc func(ir.CodeExpr)

var filename string
var rd io.RuneReader
var mod *ir.Module
var consumer_stack *list.List = list.New()

func init() {
	flag.StringVar(&filename, "i", "test0.bf", "-i <filename>")
}

func emit(e ir.CodeExpr) {
	Assert(consumer_stack.Len() > 0, 20)
	consumer := consumer_stack.Front().Value.(ConsFunc)
	consumer(e)
}

func down(f ConsFunc) {
	consumer_stack.PushFront(f)
}

func up() {
	Assert(consumer_stack.Len() > 1, 20)
	consumer_stack.Remove(consumer_stack.Front())
}

func depth() int {
	return consumer_stack.Len()
}

func module() {
	p := &ir.Local{}
	p.Name("$p")
	p.Type(types.I32)

	v := &ir.Local{}
	v.Name("$v")
	v.Type(types.I32)

	f := &ir.FuncExpr{}
	f.Name("$start")
	f.Locals = append(f.Locals, p, v)

	down(func(e ir.CodeExpr) {
		f.Code = append(f.Code, e)
	})

	ip := &ir.Param{}
	ip.Type(types.I32)

	imp := &ir.Import{}
	imp.Name("$print")
	imp.Mod = "spectest"
	imp.Func = "print"
	imp.Params = append(imp.Params, ip)

	s := &ir.StartExpr{}
	s.Var = ir.ThisVar("$start")

	m := &ir.Memory{}
	m.Initial = 1
	m.Max = 1

	mod.Mem = m
	mod.Func = append(mod.Func, f)
	mod.Imp = append(mod.Imp, imp)
	mod.Start = s
}

var Vv = ir.ThisVar("$v")
var Pv = ir.ThisVar("$p")
var Pr = ir.ThisVar("$print")

func do(cmd string) {
	switch cmd {
	case "+":
		{
			get := &ir.GetLocalExpr{}
			get.Var = Pv

			ld := &ir.LoadExpr{}
			ld.Size = ir.Load8
			ld.Type = types.I32
			ld.Signed = false
			ld.Expr = get

			set := &ir.SetLocalExpr{}
			set.Var = Vv
			set.Expr = ld

			emit(set)
		}
		{
			get := &ir.GetLocalExpr{Var: Vv}

			inc := &ir.DyadicOp{}
			inc.Left = get
			inc.Right = &ir.ConstExpr{Type: types.I32, Value: 1}
			op := ops.Dyadic(types.I32, types.I32, ops.Add)
			inc.Op = &op

			set := &ir.SetLocalExpr{}
			set.Var = Vv
			set.Expr = inc

			emit(set)
		}
		{
			set := &ir.StoreExpr{}
			set.Type = types.I32
			set.Size = ir.Load8
			set.Expr = &ir.GetLocalExpr{Var: Pv}
			set.Value = &ir.GetLocalExpr{Var: Vv}
			emit(set)
		}
	case "-":
		{
			get := &ir.GetLocalExpr{}
			get.Var = Pv

			ld := &ir.LoadExpr{}
			ld.Size = ir.Load8
			ld.Type = types.I32
			ld.Signed = false
			ld.Expr = get

			inc := &ir.DyadicOp{}
			inc.Left = ld
			inc.Right = &ir.ConstExpr{Type: types.I32, Value: -1}
			op := ops.Dyadic(types.I32, types.I32, ops.Add)
			inc.Op = &op

			set := &ir.StoreExpr{}
			set.Type = types.I32
			set.Size = ir.Load8
			set.Expr = &ir.GetLocalExpr{Var: Pv}
			set.Value = inc
			emit(set)
		}
	case ".":
		{
			get := &ir.LoadExpr{}
			get.Type = types.I32
			get.Signed = false
			get.Size = ir.Load8
			get.Expr = &ir.GetLocalExpr{Var: Pv}

			call := &ir.CallImportExpr{}
			call.Var = Pr
			call.Params = append(call.Params, get)

			emit(call)
		}
	case ">":
		cmp := &ir.DyadicOp{}
		cmp.Left = &ir.GetLocalExpr{Var: Pv}
		cmp.Right = &ir.ConstExpr{Type: types.I32, Value: 30000}
		op := ops.Dyadic(types.I32, types.I32, ops.Eq)
		cmp.Op = &op

		inc := &ir.DyadicOp{}
		inc.Left = &ir.GetLocalExpr{Var: Pv}
		inc.Right = &ir.ConstExpr{Type: types.I32, Value: 1}
		aop := ops.Dyadic(types.I32, types.I32, ops.Add)
		inc.Op = &aop

		set := &ir.SetLocalExpr{}
		set.Var = Pv
		set.Expr = inc

		cond := &ir.If{}
		cond.CondExpr = cmp
		cond.Expr = append(cond.Expr, &ir.SetLocalExpr{Var: Pv, Expr: &ir.ConstExpr{Type: types.I32, Value: 0}})
		cond.ElseExpr = append(cond.ElseExpr, set)
		emit(cond)
	case "<":
		cmp := &ir.DyadicOp{}
		cmp.Left = &ir.GetLocalExpr{Var: Pv}
		cmp.Right = &ir.ConstExpr{Type: types.I32, Value: 0}
		op := ops.Dyadic(types.I32, types.I32, ops.Eq)
		cmp.Op = &op

		inc := &ir.DyadicOp{}
		inc.Left = &ir.GetLocalExpr{Var: Pv}
		inc.Right = &ir.ConstExpr{Type: types.I32, Value: -1}
		aop := ops.Dyadic(types.I32, types.I32, ops.Add)
		inc.Op = &aop

		set := &ir.SetLocalExpr{}
		set.Var = Pv
		set.Expr = inc

		cond := &ir.If{}
		cond.CondExpr = cmp
		cond.Expr = append(cond.Expr, &ir.SetLocalExpr{Var: Pv, Expr: &ir.ConstExpr{Type: types.I32, Value: 30000}})
		cond.ElseExpr = append(cond.ElseExpr, set)
		emit(cond)

	case "[":
		loop := &ir.Loop{}
		loop.Start.Name(fmt.Sprint("$start", strconv.Itoa(depth())))
		loop.End.Name(fmt.Sprint("$end", strconv.Itoa(depth())))
		emit(loop)

		down(func(e ir.CodeExpr) {
			loop.Expr = append(loop.Expr, e)
		})

		end := &ir.Br{}
		end.Var = ir.ThisVar(loop.End.Name())

		get := &ir.GetLocalExpr{}
		get.Var = Pv

		ld := &ir.LoadExpr{}
		ld.Size = ir.Load8
		ld.Type = types.I32
		ld.Signed = false
		ld.Expr = get

		cmp := &ir.DyadicOp{}
		cmp.Left = ld
		cmp.Right = &ir.ConstExpr{Type: types.I32, Value: 0}
		op := ops.Dyadic(types.I32, types.I32, ops.Eq)
		cmp.Op = &op

		cond := &ir.If{}
		cond.CondExpr = cmp
		cond.Expr = append(cond.Expr, end)
		cond.ElseExpr = append(cond.ElseExpr, &ir.NopExpr{})
		emit(cond)
	case "]":
		br := &ir.Br{}
		br.Var = ir.ThisVar(fmt.Sprint("$start", strconv.Itoa(depth()-1)))
		emit(br)
		up()
	case ",":
		log.Println("no input in this implementation :(")
	default:
	}
}

func compile() {
	Assert(!fn.IsNil(rd), 20)
	module()
	var err error
	for err == nil {
		var r rune
		if r, _, err = rd.ReadRune(); err == nil {
			if !unicode.IsSpace(r) {
				do(string([]rune{r}))
			}
		} else if err == io.EOF {
			//that's ok
		} else {
			log.Fatal(err)
		}
	}
}

func main() {
	flag.Parse()
	if i, err := os.Open(filename); err == nil {
		rd = bufio.NewReader(i)
		mod = &ir.Module{}
		compile()
		if o, err := os.Create(filepath.Join(filepath.Dir(filename), strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filepath.Base(filename)))+".wast")); err == nil {
			buf := bytes.NewBuffer(nil)
			if err := gen.NewWriter(buf, gen.Opts{PrettyPrint: false}).WriteExpr(mod); err == nil {
				log.Println(buf.String())
				if parsed, err := ps.Parse(bytes.NewBuffer(buf.Bytes())); err == nil {
					buf2 := bytes.NewBuffer(nil)
					if err := gen.NewWriter(buf2, gen.Opts{PrettyPrint: false}).WriteExpr(parsed); err == nil {
						log.Println(buf2.String())
						io.Copy(o, buf2)
					} else {
						Halt(100, err)
					}
				} else {
					Halt(100, err)
				}
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}
