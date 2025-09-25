package interpreter

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

func executeStatements(statements []common.Statement, env *Env) (Value, error) {
	var value Value
	var err error

	for _, stmt := range statements {
		value, err = executeStatement(stmt, env)
		if err != nil {
			if IsReturnErr(err) {
				return value, err
			}
			return Value{}, err
		}
	}
	return value, nil
}

func executeStatement(stmt common.Statement, env *Env) (Value, error) {
	switch s := stmt.(type) {
	case common.BlockStatement:
		return executeBlockStatement(s, env)
	case common.VarDeclaration:
		return executeVarDeclaration(s, env)
	case common.Assignment:
		return executeAssignment(s, env)
	case common.PrintStatement:
		return executePrintStatement(s, env)
	case common.IfStatement:
		return executeIfStatement(s, env)
	case common.WhileStatement:
		return executeWhileStatement(s, env)
	case common.FunctionDef:
		return executeFunctionDef(s, env)
	case common.ExpressionStatement:
		return executeExpressionStatement(s, env)
	case common.ReturnStatement:
		return executeReturnStatement(s, env)
	case common.BreakStatement:
		return executeBreakStatement(s, env)
	default:
		return Value{}, fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func executeBlockStatement(stmt common.BlockStatement, env *Env) (Value, error) {
	newEnv := NewEnv(env)
	return executeStatements(stmt.Statements, newEnv)
}

func executeVarDeclaration(stmt common.VarDeclaration, env *Env) (Value, error) {
	value, err := resolveExpression(stmt.Value, env)
	if err != nil {
		return Value{}, err
	}

	if !env.DefineVariable(stmt.Name, value) {
		return Value{}, fmt.Errorf("variable '%s' already defined in this scope", stmt.Name)
	}

	return Value{}, nil
}

func executeAssignment(stmt common.Assignment, env *Env) (Value, error) {
	value, err := resolveExpression(stmt.Value, env)
	if err != nil {
		return Value{}, err
	}

	if !env.AssignVariable(stmt.Name, value) {
		return Value{}, fmt.Errorf("variable '%s' not defined", stmt.Name)
	}

	return Value{}, nil
}

func executePrintStatement(stmt common.PrintStatement, env *Env) (Value, error) {
	value, err := resolveExpression(stmt.Value, env)
	if err != nil {
		return Value{}, err
	}

	fmt.Println(value.Content())
	return Value{}, nil
}

func executeIfStatement(stmt common.IfStatement, env *Env) (Value, error) {
	conditionValue, err := resolveExpression(stmt.Condition, env)
	if err != nil {
		return Value{}, err
	}

	if isTruthy(conditionValue) {
		return executeBlockStatement(stmt.Body, env)
	}

	elseIf := stmt.ElseIf
	for elseIf != nil {
		condValue, err := resolveExpression(elseIf.Condition, env)
		if err != nil {
			return Value{}, err
		}
		if isTruthy(condValue) {
			return executeBlockStatement(elseIf.Body, env)
		}
		elseIf = elseIf.ElseIf
	}

	if stmt.Else != nil {
		return executeBlockStatement(stmt.Else.Body, env)
	}

	return Value{}, nil
}

func executeWhileStatement(stmt common.WhileStatement, env *Env) (Value, error) {
	for {
		conditionValue, err := resolveExpression(stmt.Condition, env)
		if err != nil {
			return Value{}, err
		}

		if !isTruthy(conditionValue) {
			break
		}

		result, err := executeBlockStatement(stmt.Body, env)
		if err != nil {
			if IsBreakErr(err) {
				break
			} else {
				return result, err
			}
		}

	}
	return Value{}, nil
}

func executeFunctionDef(stmt common.FunctionDef, env *Env) (Value, error) {
	function := NewFunction(stmt.Parameters, stmt.Body.Statements)
	if !env.DefineFunction(stmt.Name, function) {
		return Value{}, fmt.Errorf("function '%s' already defined in this scope", stmt.Name)
	}
	return Value{}, nil
}

func executeExpressionStatement(stmt common.ExpressionStatement, env *Env) (Value, error) {
	value, err := resolveExpression(stmt.Expression, env)
	return value, err
}

func executeReturnStatement(stmt common.ReturnStatement, env *Env) (Value, error) {
	value, err := resolveExpression(stmt.Value, env)
	if err != nil {
		return Value{}, err
	}

	return value, NewReturnErr()
}

func executeBreakStatement(_ common.BreakStatement, _ *Env) (Value, error) {
	return Value{}, NewBreakErr()
}
