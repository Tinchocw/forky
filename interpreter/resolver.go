package interpreter

import (
	"fmt"
	"strconv"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
	primaryExpression "github.com/Tinchocw/Interprete-concurrente/common/expression/primary"
)

func resolveExpression(expr expression.Expression, env *Env) (Value, error) {
	return resolveBinaryOr(*expr.Root, env)
}

func resolveBinaryOr(bor expression.BinaryOr, env *Env) (Value, error) {
	var left Value
	var err error

	switch bor.Left.(type) {
	case *expression.BinaryOr:
		left, err = resolveBinaryOr(*bor.Left.(*expression.BinaryOr), env)
	case *expression.BinaryAnd:
		left, err = resolveBinaryAnd(*bor.Left.(*expression.BinaryAnd), env)
	default:
		return Value{}, fmt.Errorf("invalid left operand type for BinaryOr")
	}

	if err != nil {
		return Value{}, err
	}

	if bor.Right == nil {
		return left, nil
	}

	if isTruthy(left) {
		return left, nil
	}

	right, err := resolveBinaryAnd(*bor.Right, env)
	if err != nil {
		return Value{}, err
	}

	return right, nil
}

func resolveBinaryAnd(band expression.BinaryAnd, env *Env) (Value, error) {
	var left Value
	var err error

	switch band.Left.(type) {
	case *expression.BinaryAnd:
		left, err = resolveBinaryAnd(*band.Left.(*expression.BinaryAnd), env)
	case *expression.Equality:
		left, err = resolveEquality(*band.Left.(*expression.Equality), env)
	default:
		return Value{}, fmt.Errorf("invalid left operand type for BinaryAnd")
	}

	if err != nil {
		return Value{}, err
	}

	if band.Right == nil {
		return left, nil
	}

	if !isTruthy(left) {
		return left, nil
	}

	right, err := resolveEquality(*band.Right, env)
	if err != nil {
		return Value{}, err
	}

	return right, nil
}

func resolveEquality(eq expression.Equality, env *Env) (Value, error) {
	var left Value
	var err error

	switch eq.Left.(type) {
	case *expression.Equality:
		left, err = resolveEquality(*eq.Left.(*expression.Equality), env)
	case *expression.Comparison:
		left, err = resolveComparison(*eq.Left.(*expression.Comparison), env)
	default:
		return Value{}, fmt.Errorf("invalid left operand type for Equality")
	}

	if err != nil {
		return Value{}, err
	}

	if eq.Right == nil {
		return left, nil
	}

	right, err := resolveComparison(*eq.Right, env)
	if err != nil {
		return Value{}, err
	}

	if left.Typ != right.Typ {
		return Value{}, fmt.Errorf("type mismatch in equality comparison: %v vs %v", left.Typ, right.Typ)
	}

	switch eq.Operator.Typ {
	case common.EQUAL_EQUAL:
		return Value{Typ: VAL_BOOL, Data: left.Data == right.Data}, nil
	case common.BANG_EQUAL:
		return Value{Typ: VAL_BOOL, Data: left.Data != right.Data}, nil
	default:
		return Value{}, fmt.Errorf("unknown equality operator: %s", eq.Operator.Value)
	}
}

func resolveComparison(cmp expression.Comparison, env *Env) (Value, error) {
	var left Value
	var err error

	switch cmp.Left.(type) {
	case *expression.Comparison:
		left, err = resolveComparison(*cmp.Left.(*expression.Comparison), env)
	case *expression.Term:
		left, err = resolveTerm(*cmp.Left.(*expression.Term), env)
	default:
		return Value{}, fmt.Errorf("invalid left operand type for Comparison")
	}

	if err != nil {
		return Value{}, err
	}

	if cmp.Right == nil {
		return left, nil
	}

	right, err := resolveTerm(*cmp.Right, env)
	if err != nil {
		return Value{}, err
	}

	if left.Typ != right.Typ {
		return Value{}, fmt.Errorf("type mismatch in comparison: %v vs %v", left.Typ, right.Typ)
	}

	switch cmp.Operator.Typ {
	case common.LESS:
		if left.Typ == VAL_INT {
			return Value{Typ: VAL_BOOL, Data: left.Data.(int) < right.Data.(int)}, nil
		}
		if left.Typ == VAL_STRING {
			return Value{Typ: VAL_BOOL, Data: left.Data.(string) < right.Data.(string)}, nil
		}
		return Value{}, fmt.Errorf("operator '<' not supported for type %v", left.Typ)
	case common.LESS_EQUAL:
		if left.Typ == VAL_INT {
			return Value{Typ: VAL_BOOL, Data: left.Data.(int) <= right.Data.(int)}, nil
		}
		if left.Typ == VAL_STRING {
			return Value{Typ: VAL_BOOL, Data: left.Data.(string) <= right.Data.(string)}, nil
		}
		return Value{}, fmt.Errorf("operator '<=' not supported for type %v", left.Typ)
	case common.GREATER:
		if left.Typ == VAL_INT {
			return Value{Typ: VAL_BOOL, Data: left.Data.(int) > right.Data.(int)}, nil
		}
		if left.Typ == VAL_STRING {
			return Value{Typ: VAL_BOOL, Data: left.Data.(string) > right.Data.(string)}, nil
		}
		return Value{}, fmt.Errorf("operator '>' not supported for type %v", left.Typ)
	case common.GREATER_EQUAL:
		if left.Typ == VAL_INT {
			return Value{Typ: VAL_BOOL, Data: left.Data.(int) >= right.Data.(int)}, nil
		}
		if left.Typ == VAL_STRING {
			return Value{Typ: VAL_BOOL, Data: left.Data.(string) >= right.Data.(string)}, nil
		}
		return Value{}, fmt.Errorf("operator '>=' not supported for type %v", left.Typ)
	default:
		return Value{}, fmt.Errorf("unknown comparison operator: %s", cmp.Operator.Value)
	}
}

func resolveTerm(term expression.Term, env *Env) (Value, error) {
	var left Value
	var err error

	switch term.Left.(type) {
	case *expression.Term:
		left, err = resolveTerm(*term.Left.(*expression.Term), env)
	case *expression.Factor:
		left, err = resolveFactor(*term.Left.(*expression.Factor), env)
	default:
		return Value{}, fmt.Errorf("invalid left operand type for Term")
	}

	if err != nil {
		return Value{}, err
	}

	if term.Right == nil {
		return left, nil
	}

	right, err := resolveFactor(*term.Right, env)
	if err != nil {
		return Value{}, err
	}

	switch term.Operator.Typ {
	case common.PLUS:
		if left.Typ != right.Typ {
			return Value{Typ: VAL_STRING, Data: left.Content() + right.Content()}, nil
		}

		if left.Typ == VAL_INT {
			return Value{Typ: VAL_INT, Data: left.Data.(int) + right.Data.(int)}, nil
		}

		if left.Typ == VAL_STRING {
			return Value{Typ: VAL_STRING, Data: left.Data.(string) + right.Data.(string)}, nil
		}

		return Value{}, fmt.Errorf("operator '+' not supported for type %v and type %v", left.Typ, right.Typ)

	case common.MINUS:
		if left.Typ == VAL_INT {
			return Value{Typ: VAL_INT, Data: left.Data.(int) - right.Data.(int)}, nil
		}
		return Value{}, fmt.Errorf("operator '-' not supported for type %v and type %v", left.Typ, right.Typ)
	default:
		return Value{}, fmt.Errorf("unknown term operator: %s", term.Operator.Value)
	}
}

func resolveFactor(factor expression.Factor, env *Env) (Value, error) {
	var left Value
	var err error

	switch factor.Left.(type) {
	case *expression.Factor:
		left, err = resolveFactor(*factor.Left.(*expression.Factor), env)
	case expression.Unary:
		left, err = resolveUnary(factor.Left.(expression.Unary), env)
	default:
		if u, ok := factor.Left.(expression.Unary); ok {
			left, err = resolveUnary(u, env)
		} else {
			return Value{}, fmt.Errorf("invalid left operand type for Factor")
		}
	}

	if err != nil {
		return Value{}, err
	}

	if factor.Right == nil {
		return left, nil
	}

	right, err := resolveUnary(factor.Right, env)
	if err != nil {
		return Value{}, err
	}

	if left.Typ != right.Typ {
		return Value{}, fmt.Errorf("type mismatch in factor operation: %v vs %v", left.Typ, right.Typ)
	}

	switch factor.Operator.Typ {
	case common.ASTERISK:
		if left.Typ == VAL_INT {
			return Value{Typ: VAL_INT, Data: left.Data.(int) * right.Data.(int)}, nil
		}
		return Value{}, fmt.Errorf("operator '*' not supported for type %v", left.Typ)
	case common.SLASH:
		if left.Typ == VAL_INT {
			if right.Data.(int) == 0 {
				return Value{}, fmt.Errorf("division by zero")
			}
			return Value{Typ: VAL_INT, Data: left.Data.(int) / right.Data.(int)}, nil
		}
		return Value{}, fmt.Errorf("operator '/' not supported for type %v", left.Typ)
	default:
		return Value{}, fmt.Errorf("unknown factor operator: %s", factor.Operator.Value)
	}
}

func resolveUnary(unary expression.Unary, env *Env) (Value, error) {
	switch u := unary.(type) {
	case *primaryExpression.Primary:
		return resolvePrimary(*u, env)
	case *expression.UnaryWithOperator:
		right, err := resolveUnary(u.Right, env)
		if err != nil {
			return Value{}, err
		}
		switch u.Operator.Typ {
		case common.TILDE:
			if right.Typ == VAL_INT {
				return Value{Typ: VAL_INT, Data: -right.Data.(int)}, nil
			}
			return Value{}, fmt.Errorf("unary '~' not supported for type %v", right.Typ)
		case common.BANG:
			return Value{Typ: VAL_BOOL, Data: !isTruthy(right)}, nil
		default:
			return Value{}, fmt.Errorf("unknown unary operator: %s", u.Operator.Value)
		}
	default:
		return Value{}, fmt.Errorf("unknown unary type")
	}
}

func resolvePrimary(primary primaryExpression.Primary, env *Env) (Value, error) {
	switch p := primary.Value.(type) {
	case *primaryExpression.TokenValue:
		token := p.Token
		switch token.Typ {
		case common.NUMBER:
			num, err := strconv.Atoi(token.Value)
			if err != nil {
				return Value{}, fmt.Errorf("invalid number: %s", token.Value)
			}
			return Value{Typ: VAL_INT, Data: num}, nil
		case common.LITERAL:
			return Value{Typ: VAL_STRING, Data: token.Value}, nil
		case common.TRUE:
			return Value{Typ: VAL_BOOL, Data: true}, nil
		case common.FALSE:
			return Value{Typ: VAL_BOOL, Data: false}, nil
		case common.NONE:
			return Value{Typ: VAL_NONE, Data: nil}, nil
		case common.IDENTIFIER:
			value, found := env.GetVariable(token.Value)
			if !found {
				return Value{}, fmt.Errorf("undefined variable: %s", token.Value)
			}
			return value, nil
		default:
			return Value{}, fmt.Errorf("unknown literal type: %v", token.Typ)
		}
	case *primaryExpression.Call:
		function, found := env.GetFunction(p.Callee)
		if !found {
			return Value{}, fmt.Errorf("undefined function: %s", p.Callee)
		}

		if len(p.Arguments) != len(function.Parameters) {
			return Value{}, fmt.Errorf("expected %d arguments, got %d", len(function.Parameters), len(p.Arguments))
		}

		args := make([]Value, 0, len(p.Arguments))
		for _, argExpr := range p.Arguments {
			argValue, err := resolveExpression(*argExpr, env)
			if err != nil {
				return Value{}, err
			}
			args = append(args, argValue)
		}

		value, err := function.Call(args, env)

		if err == nil || !IsReturnErr(err) {
			return Value{}, err
		}

		return value, nil

	case *primaryExpression.GroupingExpression:
		return resolveExpression(*p.Expression, env)

		// foo(a,b)[30]
		//var a = foo(1,2)

	case *primaryExpression.ArrayAccess:
		array, found := env.GetVariable(p.ArrayName)
		if !found {
			return Value{}, fmt.Errorf("undefined array: %s", p.ArrayName)
		}

		if array.Typ != VAL_ARRAY {
			return Value{}, fmt.Errorf("variable %s is not an array", p.ArrayName)
		}

		indexValues := make([]int, 0, len(p.Indexes))
		for _, indexExpr := range p.Indexes {
			indexValue, err := resolveExpression(*indexExpr, env)
			if err != nil {
				return Value{}, err
			}
			if indexValue.Typ != VAL_INT {
				return Value{}, fmt.Errorf("array index must be an integer")
			}
			indexValues = append(indexValues, indexValue.Data.(int))
		}

		for _, index := range indexValues {
			if array.Typ != VAL_ARRAY {
				return Value{}, fmt.Errorf("attempted to index a non-array value")
			}

			arrayData := array.Data.([]Value)
			if index < 0 || index >= len(arrayData) {
				return Value{}, fmt.Errorf("array index out of bounds")
			}
			array = arrayData[index]
		}

		return array, nil
	case *primaryExpression.ArrayLiteral:
		elements := make([]Value, 0, len(p.Elements))
		for _, elemExpr := range p.Elements {
			elemValue, err := resolveExpression(*elemExpr, env)
			if err != nil {
				return Value{}, err
			}
			elements = append(elements, elemValue)
		}
		return Value{Typ: VAL_ARRAY, Data: elements}, nil
	default:
		return Value{}, fmt.Errorf("unknown primary type")
	}
}
