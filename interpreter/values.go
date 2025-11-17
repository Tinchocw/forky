package interpreter

import "fmt"

type ValueType int

const (
	VAL_INT ValueType = iota + 1
	VAL_STRING
	VAL_BOOL
	VAL_NONE
	VAL_ARRAY
	VAL_FUNCTION
)

// // String implements fmt.Stringer for ValueTypes to provide readable type names in error messages
// func (vt ValueTypes) String() string {
// 	switch vt {
// 	case VAL_INT:
// 		return "INT"
// 	case VAL_STRING:
// 		return "STRING"
// 	case VAL_BOOL:
// 		return "BOOL"
// 	case VAL_NONE:
// 		return "NONE"
// 	case VAL_ARRAY:
// 		return "ARRAY"
// 	default:
// 		return "UNKNOWN"
// 	}
// }

// type Value struct {
// 	Data ValueData
// }

type Value interface {
	Content() string
	IsTruthy() bool
	Type() ValueType
	Data() any
}

type IntValue struct {
	Value int
}

func (iv IntValue) Content() string {
	return fmt.Sprintf("%d", iv.Value)
}

func (iv IntValue) IsTruthy() bool {
	return iv.Value != 0
}

func (iv IntValue) Type() ValueType {
	return VAL_INT
}

func (iv IntValue) Data() any {
	return iv.Value
}

type StringValue struct {
	Value string
}

func (sv StringValue) Content() string {
	return sv.Value
}

func (sv StringValue) IsTruthy() bool {
	return sv.Value != ""
}

func (sv StringValue) Type() ValueType {
	return VAL_STRING
}

func (sv StringValue) Data() any {
	return sv.Value
}

type BoolValue struct {
	Value bool
}

func (bv BoolValue) Content() string {
	if bv.Value {
		return "true"
	}
	return "false"
}

func (bv BoolValue) IsTruthy() bool {
	return bv.Value
}

func (bv BoolValue) Type() ValueType {
	return VAL_BOOL
}

func (bv BoolValue) Data() any {
	return bv.Value
}

type NoneValue struct{}

func (nv NoneValue) Content() string {
	return "none"
}

func (nv NoneValue) IsTruthy() bool {
	return false
}

func (nv NoneValue) Type() ValueType {
	return VAL_NONE
}

func (nv NoneValue) Data() any {
	return nil
}

type ArrayValue struct {
	Values []Value
}

func (av ArrayValue) Content() string {
	str := "["
	for i, val := range av.Values {
		str += val.Content()
		if i < len(av.Values)-1 {
			str += ", "
		}
	}
	str += "]"
	return str
}

func (av ArrayValue) IsTruthy() bool {
	return len(av.Values) > 0
}

func (av ArrayValue) Type() ValueType {
	return VAL_ARRAY
}

func (av ArrayValue) Data() any {
	return av.Values
}

type FunctionValue struct {
	Function Function
}

func (fv FunctionValue) Content() string {
	return "<function>"
}

func (fv FunctionValue) IsTruthy() bool {
	return true
}

func (fv FunctionValue) Type() ValueType {
	return VAL_FUNCTION
}

func (fv FunctionValue) Data() any {
	return fv.Function
}
