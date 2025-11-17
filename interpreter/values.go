package interpreter

import "fmt"

type ValueTypes int

const (
	VAL_INT ValueTypes = iota + 1
	VAL_STRING
	VAL_BOOL
	VAL_NONE
	VAL_ARRAY
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
	case VAL_ARRAY:
		return "ARRAY"
	default:
		return "UNKNOWN"
	}
}

type Value struct {
	Typ  ValueTypes
	Data any
}

func (v Value) Content() string {
	switch v.Typ {
	case VAL_INT:
		return fmt.Sprintf("%d", v.Data.(int))
	case VAL_STRING:
		return v.Data.(string)
	case VAL_BOOL:
		if v.Data.(bool) {
			return "true"
		}
		return "false"
	case VAL_NONE:
		return "none"
	case VAL_ARRAY:
		arr := v.Data.([]Value)
		str := "["
		for i, val := range arr {
			str += val.Content()
			if i < len(arr)-1 {
				str += ", "
			}
		}
		str += "]"
		return str
	default:
		return ""
	}
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
	case VAL_ARRAY:
		return len(value.Data.([]Value)) > 0
	default:
		return true
	}
}
