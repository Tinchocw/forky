package interpreter

type ValueTypes int

const (
	VAL_INT ValueTypes = iota
	VAL_STRING
	VAL_BOOL
	VAL_NONE
)

// String implements fmt.Stringer for ValueTypes to provide readable type names in error messages
func (vt ValueTypes) String() string {
	switch vt {
	case VAL_INT:
		return "INT"
	case VAL_STRING:
		return "STRING"
	case VAL_BOOL:
		return "BOOL"
	case VAL_NONE:
		return "NONE"
	default:
		return "UNKNOWN"
	}
}

type Value struct {
	Typ  ValueTypes
	Data any
}

func isTruthy(value Value) bool {
	switch value.Typ {
	case VAL_BOOL:
		return value.Data.(bool)
	case VAL_NONE:
		return false
	case VAL_INT:
		return value.Data.(int) != 0
	case VAL_STRING:
		return value.Data.(string) != ""
	default:
		return true
	}
}
