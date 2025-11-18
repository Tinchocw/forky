package interpreter

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

func (fv FunctionValue) TypeName() string {
	return "FUNCTION"
}
