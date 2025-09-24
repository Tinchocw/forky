package interpreter

import (
	"fmt"
	"strconv"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

func resolveExpression(expr common.Expression, env *Env) (Value, error) {
	return resolveBinaryOr(expr.Root, env)
}

func resolveBinaryOr(bor common.BinaryOr, env *Env) (Value, error) {
	left, err := resolveBinaryAnd(bor.Left, env)
	if err != nil {
		return Value{}, err
	}

	if bor.Right == nil {
		return left, nil
	}

	if isTruthy(left) {
		return Value{Typ: VAL_BOOL, Data: true}, nil
	}

	right, err := resolveBinaryOr(*bor.Right, env)
	if err != nil {
		return Value{}, err
	}

	return right, nil
}

func resolveBinaryAnd(band common.BinaryAnd, env *Env) (Value, error) {
	left, err := resolveEquality(band.Left, env)
	if err != nil {
		return Value{}, err
	}

	if band.Right == nil {
		return left, nil
	}

	if !isTruthy(left) {
		return Value{Typ: VAL_BOOL, Data: false}, nil
	}

	right, err := resolveBinaryAnd(*band.Right, env)
	if err != nil {
		return Value{}, err
	}

	return right, nil
}

func resolveEquality(eq common.Equality, env *Env) (Value, error) {
	left, err := resolveComparison(eq.Left, env)
	if err != nil {
		return Value{}, err
	}

	if eq.Right == nil || eq.Operator == nil {
		return left, nil
	}

	right, err := resolveEquality(*eq.Right, env)
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

func resolveComparison(cmp common.Comparison, env *Env) (Value, error) {
	left, err := resolveTerm(cmp.Left, env)
	if err != nil {
		return Value{}, err
	}

	if cmp.Right == nil || cmp.Operator == nil {
		return left, nil
	}

	right, err := resolveComparison(*cmp.Right, env)
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

func resolveTerm(term common.Term, env *Env) (Value, error) {
	left, err := resolveFactor(term.Left, env)
	if err != nil {
		return Value{}, err
	}

	if term.Right == nil || term.Operator == nil {
		return left, nil
	}

	right, err := resolveTerm(*term.Right, env)
	if err != nil {
		return Value{}, err
	}

	if left.Typ != right.Typ {
		return Value{}, fmt.Errorf("type mismatch in term operation: %v vs %v", left.Typ, right.Typ)
	}

	switch term.Operator.Typ {
	case common.PLUS:
		if left.Typ == VAL_INT {
			return Value{Typ: VAL_INT, Data: left.Data.(int) + right.Data.(int)}, nil
		}
		if left.Typ == VAL_STRING {
			return Value{Typ: VAL_STRING, Data: left.Data.(string) + right.Data.(string)}, nil
		}
		return Value{}, fmt.Errorf("operator '+' not supported for type %v", left.Typ)
	case common.MINUS:
		if left.Typ == VAL_INT {
			return Value{Typ: VAL_INT, Data: left.Data.(int) - right.Data.(int)}, nil
		}
		return Value{}, fmt.Errorf("operator '-' not supported for type %v", left.Typ)
	default:
		return Value{}, fmt.Errorf("unknown term operator: %s", term.Operator.Value)
	}
}

func resolveFactor(factor common.Factor, env *Env) (Value, error) {
	left, err := resolveUnary(factor.Left, env)
	if err != nil {
		return Value{}, err
	}

	if factor.Right == nil || factor.Operator == nil {
		return left, nil
	}

	right, err := resolveFactor(*factor.Right, env)
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

func resolveUnary(unary common.Unary, env *Env) (Value, error) {
	switch u := unary.(type) {
	case common.Primary:
		return resolvePrimary(u, env)
	case common.UnaryWithOperator:
		right, err := resolveUnary(u.Right, env)
		if err != nil {
			return Value{}, err
		}
		switch u.Operator.Typ {
		case common.MINUS:
			if right.Typ == VAL_INT {
				return Value{Typ: VAL_INT, Data: -right.Data.(int)}, nil
			}
			return Value{}, fmt.Errorf("unary '-' not supported for type %v", right.Typ)
		case common.BANG:
			return Value{Typ: VAL_BOOL, Data: !isTruthy(right)}, nil
		default:
			return Value{}, fmt.Errorf("unknown unary operator: %s", u.Operator.Value)
		}
	default:
		return Value{}, fmt.Errorf("unknown unary type")
	}
}

func resolvePrimary(primary common.Primary, env *Env) (Value, error) {
	switch p := primary.Value.(type) {
	case common.Token:
		switch p.Typ {
		case common.NUMBER:
			num, err := strconv.Atoi(p.Value)
			if err != nil {
				return Value{}, fmt.Errorf("invalid number: %s", p.Value)
			}
			return Value{Typ: VAL_INT, Data: num}, nil
		case common.LITERAL:
			return Value{Typ: VAL_STRING, Data: p.Value}, nil
		case common.TRUE:
			return Value{Typ: VAL_BOOL, Data: true}, nil
		case common.FALSE:
			return Value{Typ: VAL_BOOL, Data: false}, nil
		case common.NONE:
			return Value{Typ: VAL_NONE, Data: nil}, nil
		case common.IDENTIFIER:
			value, found := env.GetVariable(p.Value)
			if !found {
				return Value{}, fmt.Errorf("undefined variable: %s", p.Value)
			}
			return value, nil
		default:
			return Value{}, fmt.Errorf("unknown literal type: %v", p.Typ)
		}
	case common.Call:
		function, found := env.GetFunction(p.Callee)
		if !found {
			return Value{}, fmt.Errorf("undefined function: %s", p.Callee)
		}

		if len(p.Arguments) != len(function.Parameters) {
			return Value{}, fmt.Errorf("expected %d arguments, got %d", len(function.Parameters), len(p.Arguments))
		}

		args := make([]Value, 0, len(p.Arguments))
		for _, argExpr := range p.Arguments {
			argValue, err := resolveExpression(argExpr, env)
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

	case common.GroupingExpression:
		return resolveExpression(p.Expression, env)
	default:
		return Value{}, fmt.Errorf("unknown primary type")
	}
}
