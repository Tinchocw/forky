package interpreter

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Function struct {
	Parameters []string
	Statements []common.Statement
}

func NewFunction(params []string, statements []common.Statement) Function {
	return Function{
		Parameters: params,
		Statements: statements,
	}
}

func (f Function) Call(args []Value, env *Env) (Value, error) {
	if len(args) != len(f.Parameters) {
		return Value{}, fmt.Errorf("expected %d arguments, got %d", len(f.Parameters), len(args))
	}

	functionEnv := NewEnv(env)
	for idx, argValue := range args {
		functionEnv.DefineVariable(f.Parameters[idx], argValue)
	}

	return executeStatements(f.Statements, functionEnv)
}
