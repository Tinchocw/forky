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

func (i *Interpreter) Interpret(program common.Program) error {
	_, err := executeStatements(program.Statements, i.globalEnv)
	return err
}
