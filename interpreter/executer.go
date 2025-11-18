package interpreter

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common/statement"
	"github.com/Tinchocw/Interprete-concurrente/common/statement/assignment"
	"github.com/Tinchocw/Interprete-concurrente/common/statement/block"
	"github.com/Tinchocw/Interprete-concurrente/common/statement/declaration"
	"github.com/Tinchocw/Interprete-concurrente/common/statement/extra"
	"github.com/Tinchocw/Interprete-concurrente/common/statement/flow"
	"github.com/Tinchocw/Interprete-concurrente/common/statement/function"
	"github.com/Tinchocw/Interprete-concurrente/interpreter/errors"
)

func executeStatements(statements []statement.Statement, env *Env) (Value, error) {
	var value Value
	var err error

	for _, stmt := range statements {
		value, err = executeStatement(stmt, env)
		if err != nil {
			if errors.IsReturnErr(err) {
				return value, err
			}
			return nil, err
		}
	}
	return value, nil
}

func executeStatement(stmt statement.Statement, env *Env) (Value, error) {
	switch s := stmt.(type) {
	case *block.BlockStatement:
		return executeBlockStatement(s, env)
	case *declaration.VarDeclaration:
		return executeVarDeclaration(s, env)
	case *declaration.ArrayDeclaration:
		return executeArrayDeclaration(s, env)
	case *assignment.VarAssignment:
		return executeVarAssignment(s, env)
	case *assignment.ArrayAssignment:
		return executeArrayAssignment(s, env)
	case *extra.PrintStatement:
		return executePrintStatement(s, env)
	case *extra.ForkBlockStatement:
		return executeForkBlockStatement(s, env)
	case *extra.ForkArrayStatement:
		return excecuteForkArrayStatement(s, env)
	case *flow.IfStatement:
		return executeIfStatement(s, env)
	case *flow.WhileStatement:
		return executeWhileStatement(s, env)
	case *function.FunctionDef:
		return executeFunctionDef(s, env)
	case *statement.ExpressionStatement:
		return executeExpressionStatement(s, env)
	case *function.ReturnStatement:
		return executeReturnStatement(s, env)
	case *flow.BreakStatement:
		return executeBreakStatement(s, env)
	default:
		return nil, fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func executeBlockStatement(stmt *block.BlockStatement, env *Env) (Value, error) {
	newEnv := NewEnv(env)
	return executeStatements(stmt.Statements, newEnv)
}

func executeVarDeclaration(stmt *declaration.VarDeclaration, env *Env) (Value, error) {
	var value Value = NoneValue{}

	if stmt.Value != nil {
		var err error
		value, err = resolveExpression(*stmt.Value, env)
		if err != nil {
			return nil, err
		}
	}

	if !env.DefineVariable(stmt.Name, value) {
		return nil, fmt.Errorf("variable '%s' already defined in this scope", stmt.Name)
	}

	return nil, nil
}

func executeArrayDeclaration(stmt *declaration.ArrayDeclaration, env *Env) (Value, error) {
	lengths := []IntValue{}

	for _, lenExpr := range stmt.Lengths {
		lenValue, err := resolveExpression(*lenExpr, env)
		if err != nil {
			return nil, err
		}

		if lenValue.Type() != VAL_INT {
			return nil, fmt.Errorf("array length must be an integer, got %v", lenValue.Type())
		}

		lengths = append(lengths, lenValue.(IntValue))
	}

	var value Value = NoneValue{}

	if stmt.Value != nil {
		var err error
		value, err = resolveExpression(*stmt.Value, env)
		if err != nil {
			return nil, err
		}
	}

	array := createArrayRecursive(lengths, value)

	if !env.DefineVariable(stmt.Name, array) {
		return nil, fmt.Errorf("array '%s' already defined in this scope", stmt.Name)
	}

	return nil, nil
}

func createArrayRecursive(lengths []IntValue, cellValue Value) Value {
	if len(lengths) == 0 {
		return cellValue
	}

	size := int(lengths[0].Value)
	array := make([]Value, size)
	for i := range size {
		array[i] = createArrayRecursive(lengths[1:], cellValue)
	}

	return ArrayValue{Values: array}
}

func executeVarAssignment(stmt *assignment.VarAssignment, env *Env) (Value, error) {
	value, err := resolveExpression(*stmt.Value, env)
	if err != nil {
		return nil, err
	}

	if !env.AssignVariable(stmt.Name, value) {
		return nil, fmt.Errorf("variable '%s' not defined", stmt.Name)
	}

	return nil, nil
}

func executeArrayAssignment(stmt *assignment.ArrayAssignment, env *Env) (Value, error) {
	array, ok := env.GetVariable(stmt.Name)
	if !ok {
		return nil, fmt.Errorf("array '%s' not defined", stmt.Name)
	}

	indexes := []int{}
	for _, indexExpr := range stmt.Indexes {
		indexValue, err := resolveExpression(*indexExpr, env)
		if err != nil {
			return nil, err
		}

		if indexValue.Type() != VAL_INT {
			return nil, fmt.Errorf("array index must be an integer, got %v", indexValue.Type())
		}

		indexes = append(indexes, int(indexValue.(IntValue).Value))
	}

	value, err := resolveExpression(*stmt.Value, env)
	if err != nil {
		return nil, err
	}

	for _, idx := range indexes {
		if (*array).Type() != VAL_ARRAY {
			return nil, fmt.Errorf("expected array type during assignment, got %v", (*array).Type())
		}

		arrayValue := (*array).(ArrayValue)

		if idx < 0 || idx >= len(arrayValue.Values) {
			return nil, fmt.Errorf("array index %d out of bounds", idx)
		}

		array = &arrayValue.Values[idx]
	}

	*array = value

	return nil, nil
}

func executePrintStatement(stmt *extra.PrintStatement, env *Env) (Value, error) {
	value, err := resolveExpression(*stmt.Value, env)
	if err != nil {
		return nil, err
	}

	fmt.Println(value.Content())
	return nil, nil
}

func executeForkBlockStatement(stmt *extra.ForkBlockStatement, env *Env) (Value, error) {
	done := make(chan error)

	for _, s := range stmt.Block.Statements {
		newEnv := NewEnv(env)
		go func(st statement.Statement, e *Env) {
			_, err := executeStatement(st, e)
			done <- err
		}(s, newEnv)
	}

	for range stmt.Block.Statements {
		if err := <-done; err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func excecuteForkArrayStatement(stmt *extra.ForkArrayStatement, env *Env) (Value, error) {
	value, err := resolveExpression(*stmt.Array, env)
	if err != nil {
		return nil, err
	}

	if value.Type() != VAL_ARRAY {
		return nil, fmt.Errorf("expected array type in fork array statement, got %v", value.Type())
	}

	arrayValue := value.(ArrayValue).Values

	done := make(chan error)

	for index, elem := range arrayValue {
		newEnv := NewEnv(env)

		if stmt.IndexName != nil {
			if !newEnv.DefineVariable(*stmt.IndexName, IntValue{Value: index}) {
				return nil, fmt.Errorf("variable '%s' already defined in this scope", *stmt.IndexName)
			}
		}

		if stmt.ElemName != nil {
			if !newEnv.DefineVariable(*stmt.ElemName, elem) {
				return nil, fmt.Errorf("variable '%s' already defined in this scope", *stmt.ElemName)
			}
		}

		go func(e *Env) {
			_, err := executeBlockStatement(stmt.Block, e)
			done <- err
		}(newEnv)
	}

	for range arrayValue {
		if err := <-done; err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func executeIfStatement(stmt *flow.IfStatement, env *Env) (Value, error) {
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

func executeWhileStatement(stmt *flow.WhileStatement, env *Env) (Value, error) {
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
			if errors.IsBreakErr(err) {
				break
			} else {
				return result, err
			}
		}

	}
	return nil, nil
}

func executeFunctionDef(stmt *function.FunctionDef, env *Env) (Value, error) {
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

func executeReturnStatement(stmt *function.ReturnStatement, env *Env) (Value, error) {
	if stmt.Value == nil {
		return nil, errors.NewReturnErr()
	}

	value, err := resolveExpression(*stmt.Value, env)
	if err != nil {
		return nil, err
	}

	return value, errors.NewReturnErr()
}

func executeBreakStatement(_ *flow.BreakStatement, _ *Env) (Value, error) {
	return nil, errors.NewBreakErr()
}
