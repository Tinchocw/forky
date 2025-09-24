package interpreter

import "fmt"

type ReturnErr error

func NewReturnErr() ReturnErr {
	return ReturnErr(fmt.Errorf("return"))
}

func IsReturnErr(err error) bool {
	_, ok := err.(ReturnErr)
	return ok
}

type BreakErr error

func NewBreakErr() BreakErr {
	return BreakErr(fmt.Errorf("break"))
}

func IsBreakErr(err error) bool {
	_, ok := err.(BreakErr)
	return ok
}
