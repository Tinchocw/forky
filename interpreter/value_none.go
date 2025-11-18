package interpreter

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

func (nv NoneValue) TypeName() string {
	return "NONE"
}
