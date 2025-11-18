package errors

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
