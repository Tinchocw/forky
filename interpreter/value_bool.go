package interpreter

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

func (bv BoolValue) TypeName() string {
	return "BOOL"
}