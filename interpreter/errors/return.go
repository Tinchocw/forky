package errors

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
