package interpreter

type ValueType int

const (
	VAL_INT ValueType = iota + 1
	VAL_STRING
	VAL_BOOL
	VAL_NONE
	VAL_ARRAY
	VAL_FUNCTION
)

type Value interface {
	Content() string
	IsTruthy() bool
	Type() ValueType
	Data() any
	TypeName() string
}
