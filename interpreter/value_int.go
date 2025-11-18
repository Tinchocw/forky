package interpreter

import "fmt"

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

func (iv IntValue) TypeName() string {
	return "INT"
}
