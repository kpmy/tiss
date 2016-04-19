package types

type Type string

const (
	I32 Type = "i32"
	I64 Type = "i64"
	F32 Type = "f32"
	F64 Type = "f64"
)

func (t Type) Valid() bool {
	return t == I32 || t == I64 || t == F32 || t == F64
}
