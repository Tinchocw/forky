package interpreter

import "github.com/Tinchocw/forky/common/statement"

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

	if err != nil {
		return "", err
	}

	if value == nil {
		return "", nil
	}

	return value.Content(), nil
}
