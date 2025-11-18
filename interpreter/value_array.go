package interpreter

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

func (av ArrayValue) TypeName() string {
	return "ARRAY"
}
