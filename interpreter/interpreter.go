package interpreter

import (
	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Interpreter struct {
	globalEnv *Env
}

func NewInterpreter() Interpreter {
	return Interpreter{
		globalEnv: NewEnv(nil),
	}
}

func (i *Interpreter) Execute(program common.Program) (string, error) {
	value, err := executeStatements(program.Statements, i.globalEnv)
	return value.Content(), err
}
