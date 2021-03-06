package sexpr

import (
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/kpmy/tiss/gen"
	"github.com/kpmy/tiss/ir"
	"github.com/kpmy/ypk/fn"
	. "github.com/kpmy/ypk/tc"
)

const InvalidAtom Error = "invalid name of atom"

const NameTag = "sexpr"

type Named interface {
	Name() string
}

type wr struct {
	base  io.Writer
	opts  gen.Opts
	depth int
	line  int
}

func (w *wr) Ln() {
	if w.opts.PrettyPrint && w.line >= 0 {
		w.base.Write([]byte("\n"))
	}
	w.line++
}

func (w *wr) Tab() {
	if w.opts.PrettyPrint {
		var buf []rune
		for i := 0; i < w.depth; i++ {
			buf = append(buf, '\t')
		}
		w.base.Write([]byte(string(buf)))
	}
}

func (w *wr) Raw(s string) {
	if _, err := w.base.Write([]byte(s)); err != nil {
		panic(err)
	}
}

const invalidAtomChars = ""

func validateAtomString(s string) {
	if strings.ContainsAny(s, invalidAtomChars) {
		panic(InvalidAtom)
	} else {
		r := []rune(s)
		for i := 0; i < len(r); i++ {
			Assert(!unicode.IsSpace(r[i]), 20, InvalidAtom)
		}
	}
}

func validateString(r []rune) {

}

func (w *wr) Atom(_v interface{}) {
	switch v := _v.(type) {
	case string:
		Assert(v != "", 20)
		validateAtomString(v)
		w.Raw(v)
	case ir.Variable:
		w.Atom(v.ValueOf())
	case int:
		w.Raw(strconv.Itoa(v))
	case uint:
		w.Raw(strconv.FormatUint(uint64(v), 10))
	case []rune:
		validateString(v)
		w.Raw(`"` + string(v) + `"`)
	case float64:
		w.Raw(strconv.FormatFloat(v, 'g', 14, 64))
	case float32:
		w.Raw(strconv.FormatFloat(float64(v), 'g', 14, 32))
	default:
		Halt(100, "wrong atom ", reflect.TypeOf(v))
	}
}

func getName(i interface{}) (ret string) {
	if named, ok := i.(Named); ok {
		ret = named.Name()
	} else {
		t := reflect.ValueOf(i).Elem().Type()
		for i := 0; i < t.NumField() && ret == ""; i++ {
			s := t.Field(i)
			ret = s.Tag.Get(NameTag)
		}
	}
	return
}

func (w *wr) WriteValue(v interface{}) (err error) {
	Try(func(x ...interface{}) (ret interface{}) {
		x0 := x[0]
		switch v := x0.(type) {
		case ir.Expression:
			err = w.WriteExpr(v)
		default:
			w.Atom(x0)
		}
		return
	}, v).Do()
	//.Catch(nil, func(e error) {err = e})
	return
}

func (w *wr) WriteExpr(e ir.Expression) (err error) {
	Assert(!fn.IsNil(e), 20)
	if err = e.Validate(); err == nil {
		Do(func() {
			w.Ln()
			w.Tab()
			w.Raw("(")
			w.Atom(getName(e))
			if n, ok := e.(ir.Node); ok {
				if el := n.Children(); len(el) > 0 {
					w.depth++
					for _, _e := range el {
						w.Raw(" ")
						if err = w.WriteValue(_e); err != nil {
							panic(err)
						}
					}
					w.depth--
				}
			}
			w.Raw(")")
		}).Do()
		//.Catch(nil, func(e error) {err = e})
	}
	return err
}

func New(w io.Writer, o ...gen.Opts) gen.Writer {
	Assert(!fn.IsNil(w), 20)
	ret := &wr{base: w}
	if len(o) > 0 {
		ret.opts = o[0]
	}
	if ret.opts.PrettyPrint {
		//ret.Raw(";; github.com/kpmy/tiss/generator")
		ret.line = -1
	}

	return ret
}
