package interpreter

type ReturnErr struct{}

func (e ReturnErr) Error() string {
	return "return"
}

func NewReturnErr() ReturnErr {
	return ReturnErr{}
}

func IsReturnErr(err error) bool {
	_, ok := err.(ReturnErr)
	return ok
}

type BreakErr struct{}

func (e BreakErr) Error() string {
	return "break"
}

func NewBreakErr() BreakErr {
	return BreakErr{}
}

func IsBreakErr(err error) bool {
	_, ok := err.(BreakErr)
	return ok
}
