package interpreter

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

func (sv StringValue) TypeName() string {
	return "STRING"
}
