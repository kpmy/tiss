package ops

import (
	"fmt"
	"strings"

	"github.com/kpmy/tiss/ir/types"
	. "github.com/kpmy/ypk/tc"
)

type Op string

//monadic int
const Clz Op = "clz"
const Ctz Op = "ctz"
const PopCnt Op = "popcnt"
const Eqz Op = "eqz"

//monadic float
const Neg Op = "neg"
const Abs Op = "abs"
const Sqrt Op = "sqrt"
const Ceil Op = "ceil"
const Floor Op = "floor"
const Nearest Op = "nearest"
const Trunc Op = "trunc"

//dyadic int
const Add Op = "add"
const Sub Op = "sub"
const Mul Op = "mul"

//const DivS Op = "div_s"
//const DivU Op = "div_u"
const Div Op = "div"

//const RemS Op = "rem_s"
//const RemU Op = "rem_u"
const Rem Op = "rem"
const And Op = "and"
const Or Op = "or"
const Xor Op = "xor"
const Shl Op = "shl"

//const ShrS Op = "shr_s"
//const ShrU Op = "shr_u"
const Shr Op = "shr"
const RotL Op = "rotl"
const RotR Op = "rotr"

//dyadic float
//const Add Op = "add"
//const Sub Op = "sub"
//const Mul Op = "mul"
//const Div Op = "div"
const Min Op = "min"
const Max Op = "max"
const CopySign Op = "copysign"

//compare int
const Eq Op = "eq"
const Ne Op = "ne"

//const LtS Op = "lt_s"
//const LtU Op = "lt_u"
//const LeS Op = "le_s"
//const LeU Op = "le_u"
//const GtS Op = "gt_s"
//const GtU Op = "gt_u"
//const GeS Op = "ge_s"
//const GeU Op = "ge_u"
const Lt Op = "lt"
const Le Op = "le"
const Gt Op = "gt"
const Ge Op = "ge"

//compare float
//const Eq Op = "eq"
//const Ne Op = "ne"
//const Lt Op = "lt"
//const Le Op = "le"
//const Gt Op = "gt"
//const Ge Op = "ge"

//conversion
const Wrap Op = "wrap" // i32_ /i64
//const Trunc Op = "trunc"             // i32_, i64_ * _s, _u * /f32, /f64
const Extend Op = "extend"           // i64_ * _s, _u * /i32
const Convert Op = "convert"         // f32_, f64_ * _s, _u * /i32, /i64
const Demote Op = "demote"           // f32_ /f64
const Promote Op = "promote"         // f64_ /f32
const Reinterpret Op = "reinterpret" // i32_, i64_ * /f32, /f64

type MonadicOpCode struct {
	typ types.Type
	op  Op
}

func (m MonadicOpCode) String() string {
	return fmt.Sprint(m.typ, ".", m.op)
}

var i_ops = fmt.Sprint(Clz, Ctz, PopCnt, Eqz)
var f_ops = fmt.Sprint(Neg, Abs, Sqrt, Floor, Trunc, Ceil, Nearest)

func Monadic(typ types.Type, op Op) (ret MonadicOpCode) {
	switch typ {
	case types.I32, types.I64:
		Assert(strings.Contains(i_ops, string(op)), 20)
	case types.F32, types.F64:
		Assert(strings.Contains(f_ops, string(op)), 20)
	default:
		Halt(100)
	}
	ret.typ = typ
	ret.op = op
	return
}

type DyadicOpCode struct {
	l, r   types.Type
	op     Op
	signed bool
}

var ii_ops = fmt.Sprint(Add, Sub, Mul, Div, Rem, And, Or, Xor, Shl, Shr, RotL, RotR, Eq, Ne, Le, Lt, Ge, Gt)
var ii_sign = fmt.Sprint(Div, Rem, Shr, Le, Lt, Ge, Gt)

func (d DyadicOpCode) String() (ret string) {
	if d.l == d.r {
		ret = fmt.Sprint(d.l, ".", d.op)
		if strings.Contains(ii_sign, string(d.op)) {
			if d.signed {
				ret = fmt.Sprint(ret, "_s")
			} else {
				ret = fmt.Sprint(ret, "_u")
			}
		}
	} else {
		return "unsupported"
	}
	return
}

func Dyadic(l, r types.Type, op Op, signed ...bool) (ret DyadicOpCode) {
	if l == r {
		switch l {
		case types.I32, types.I64:
			Assert(strings.Contains(ii_ops, string(op)), 21)
			if strings.Contains(ii_sign, string(op)) {
				Assert(len(signed) > 0, 21)
			}
		default:
			Halt(100)
		}
	} else {
		Halt(100)
	}
	ret.l = l
	ret.r = r
	ret.op = op
	if len(signed) > 0 {
		ret.signed = signed[0]
	}
	return
}
