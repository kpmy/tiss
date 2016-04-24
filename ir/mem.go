package ir

import (
	"fmt"
	"strconv"

	"github.com/kpmy/tiss/ir/types"
	"github.com/kpmy/ypk/fn"
	. "github.com/kpmy/ypk/tc"
)

type CurrentMemoryExpr struct {
	ns `sexpr:"current_memory"`
}

func (*CurrentMemoryExpr) Validate() error { return nil }
func (*CurrentMemoryExpr) Eval()           {}

type GrowMemoryExpr struct {
	ns   `sexpr:"grow_memory"`
	Expr CodeExpr
}

func (m *GrowMemoryExpr) Validate() (err error) {
	if fn.IsNil(m.Expr) {
		err = Error("no expression for grow memory")
	}
	return
}

func (m *GrowMemoryExpr) Children() (ret []interface{}) {
	return append(ret, m.Expr)
}

func (*GrowMemoryExpr) Eval() {}

type LoadSize int

const (
	Load8  LoadSize = 8
	Load16 LoadSize = 16
	Load32 LoadSize = 32
)

func (s LoadSize) Valid() bool {
	return s == Load8 || s == Load16 || s == Load32
}

type LoadExpr struct {
	Type   types.Type
	Size   LoadSize
	Signed bool
	Offset uint
	Align  uint
	Expr   CodeExpr
}

func (l *LoadExpr) Name() (ret string) {
	ret = fmt.Sprint(l.Type, ".", "load", l.Size)
	if l.Type == types.I64 || l.Type == types.I32 {
		if l.Signed {
			ret = fmt.Sprint(ret, "_s")
		} else {
			ret = fmt.Sprint(ret, "_u")
		}
	}
	return
}

func (l *LoadExpr) Validate() error {
	if !l.Type.Valid() {
		return Error("invalid load type")
	}

	if !l.Size.Valid() {
		return Error("invalid load size")
	}

	if l.Align > 1 && l.Align%2 == 1 {
		return Error("load align must be power of 2")
	}

	if fn.IsNil(l.Expr) {
		return Error("load expression is nil")
	}
	return nil
}

func (l *LoadExpr) Children() (ret []interface{}) {

	if l.Offset > 0 {
		ret = append(ret, fmt.Sprint("offset=", strconv.FormatUint(uint64(l.Offset), 32)))
	}

	if l.Align > 0 {
		ret = append(ret, fmt.Sprint("align=", strconv.FormatUint(uint64(l.Align), 32)))
	}

	ret = append(ret, l.Expr)
	return
}

func (l *LoadExpr) Eval() {}

type StoreExpr struct {
	Type        types.Type
	Size        LoadSize
	Offset      uint
	Align       uint
	Expr, Value CodeExpr
}

func (s *StoreExpr) Name() (ret string) {
	ret = fmt.Sprint(s.Type, ".", "store", s.Size)
	return
}

func (s *StoreExpr) Validate() error {
	if !s.Type.Valid() {
		return Error("invalid store type")
	}

	if !s.Size.Valid() {
		return Error("invalid store size")
	}

	if s.Align > 1 && s.Align%2 == 1 {
		return Error("store align must be power of 2")
	}

	if fn.IsNil(s.Expr) {
		return Error("store expression is nil")
	}

	if fn.IsNil(s.Value) {
		return Error("store value expression is nil")
	}
	return nil
}

func (s *StoreExpr) Children() (ret []interface{}) {

	if s.Offset > 0 {
		ret = append(ret, fmt.Sprint("offset=", strconv.FormatUint(uint64(s.Offset), 32)))
	}

	if s.Align > 0 {
		ret = append(ret, fmt.Sprint("align=", strconv.FormatUint(uint64(s.Align), 32)))
	}

	ret = append(ret, s.Expr, s.Value)
	return
}

func (l *StoreExpr) Eval() {}

const PageSize = 65535

type Memory struct {
	ns       `sexpr:"memory"`
	Initial  uint
	Max      uint
	Segments []*Segment
}

type Segment struct {
	ns     `sexpr:"segment"`
	Offset uint
	Data   string
}

func (s *Segment) Validate() error { return nil }
func (s *Segment) Children() (ret []interface{}) {
	return append(ret, s.Offset, []rune(s.Data))
}

func (m *Memory) Validate() error {
	if m.Initial > m.Max {
		return Error("error page limit")
	}

	off := -1
	for i, o := range m.Segments {
		if off > int(o.Offset) {
			return Error(fmt.Sprint("segments overlap", i-1, o.Offset))
		}
		off += int(o.Offset) + len([]rune(o.Data))
	}

	if off >= int(m.Max)*PageSize {
		return Error("not enough memory")
	}
	return nil
}

func (m *Memory) Children() (ret []interface{}) {
	ret = append(ret, m.Initial, m.Max)

	for _, s := range m.Segments {
		ret = append(ret, s)
	}

	return
}
