package interpreter

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common/statement"
)

func executeStatements(statements []statement.Statement, env *Env) (Value, error) {
	var value Value
	var err error

	for _, stmt := range statements {
		value, err = executeStatement(stmt, env)
		if err != nil {
			if IsReturnErr(err) {
				return value, err
			}
			return nil, err
		}
	}
	return value, nil
}

func executeStatement(stmt statement.Statement, env *Env) (Value, error) {
	switch s := stmt.(type) {
	case *statement.BlockStatement:
		return executeBlockStatement(s, env)
	case *statement.VarDeclaration:
		return executeVarDeclaration(s, env)
	case *statement.Assignment:
		return executeAssignment(s, env)
	case *statement.PrintStatement:
		return executePrintStatement(s, env)
	case *statement.IfStatement:
		return executeIfStatement(s, env)
	case *statement.WhileStatement:
		return executeWhileStatement(s, env)
	case *statement.FunctionDef:
		return executeFunctionDef(s, env)
	case *statement.ExpressionStatement:
		return executeExpressionStatement(s, env)
	case *statement.ReturnStatement:
		return executeReturnStatement(s, env)
	case *statement.BreakStatement:
		return executeBreakStatement(s, env)
	default:
		return nil, fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func executeBlockStatement(stmt *statement.BlockStatement, env *Env) (Value, error) {
	newEnv := NewEnv(env)
	return executeStatements(stmt.Statements, newEnv)
}

func executeVarDeclaration(stmt *statement.VarDeclaration, env *Env) (Value, error) {
	value, err := resolveExpression(*stmt.Value, env)
	if err != nil {
		return nil, err
	}

	if !env.DefineVariable(*stmt.Name, value) {
		return nil, fmt.Errorf("variable '%s' already defined in this scope", *stmt.Name)
	}

	return nil, nil
}

func executeAssignment(stmt *statement.Assignment, env *Env) (Value, error) {
	value, err := resolveExpression(*stmt.Value, env)
	if err != nil {
		return nil, err
	}

	if !env.AssignVariable(*stmt.Name, value) {
		return nil, fmt.Errorf("variable '%s' not defined", *stmt.Name)
	}

	return nil, nil
}

func executePrintStatement(stmt *statement.PrintStatement, env *Env) (Value, error) {
	value, err := resolveExpression(*stmt.Value, env)
	if err != nil {
		return nil, err
	}

	fmt.Println(value.Content())
	return nil, nil
}

func executeIfStatement(stmt *statement.IfStatement, env *Env) (Value, error) {
	conditionValue, err := resolveExpression(*stmt.Condition, env)
	if err != nil {
		return nil, err
	}

	if conditionValue.IsTruthy() {
		return executeBlockStatement(stmt.Body, env)
	}

	elseIf := stmt.ElseIf
	for elseIf != nil {
		condValue, err := resolveExpression(*elseIf.Condition, env)
		if err != nil {
			return nil, err
		}
		if condValue.IsTruthy() {
			return executeBlockStatement(elseIf.Body, env)
		}
		elseIf = elseIf.ElseIf
	}

	if stmt.Else != nil {
		return executeBlockStatement(stmt.Else.Body, env)
	}

	return nil, nil
}

func executeWhileStatement(stmt *statement.WhileStatement, env *Env) (Value, error) {
	for {
		conditionValue, err := resolveExpression(*stmt.Condition, env)
		if err != nil {
			return nil, err
		}

		if !conditionValue.IsTruthy() {
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
	return nil, nil
}

func executeFunctionDef(stmt *statement.FunctionDef, env *Env) (Value, error) {
	function := NewFunction(stmt.Parameters, stmt.Body.Statements)
	if !env.DefineVariable(*stmt.Name, FunctionValue{Function: function}) {
		return nil, fmt.Errorf("function '%s' already defined in this scope", *stmt.Name)
	}
	return nil, nil
}

func executeExpressionStatement(stmt *statement.ExpressionStatement, env *Env) (Value, error) {
	value, err := resolveExpression(*stmt.Expression, env)
	return value, err
}

func executeReturnStatement(stmt *statement.ReturnStatement, env *Env) (Value, error) {
	value, err := resolveExpression(*stmt.Value, env)
	if err != nil {
		return nil, err
	}

	return value, NewReturnErr()
}

func executeBreakStatement(_ *statement.BreakStatement, _ *Env) (Value, error) {
	return nil, NewBreakErr()
}
