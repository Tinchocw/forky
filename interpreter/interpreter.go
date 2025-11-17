package interpreter

import "github.com/Tinchocw/Interprete-concurrente/common/statement"

type Interpreter struct {
	globalEnv *Env
}

func NewInterpreter() Interpreter {
	return Interpreter{
		globalEnv: NewEnv(nil),
	}
}

func (i *Interpreter) Execute(program statement.Program) (string, error) {
	value, err := executeStatements(program.Statements, i.globalEnv)
	if value == nil {
		return "", err
	}
	return value.Content(), err
}
