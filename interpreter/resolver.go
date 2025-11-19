package interpreter

import (
	"fmt"
	"strconv"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
	"github.com/Tinchocw/forky/interpreter/errors"
)

func resolveExpression(expr expression.ExpressionNode, env *Env) (Value, error) {
	return resolveLogicalOr(*expr.Root, env)
}

func resolveLogicalOr(bor expression.LogicalOrNode, env *Env) (Value, error) {
	var left Value
	var err error

	switch bor.Left.(type) {
	case *expression.LogicalOrNode:
		left, err = resolveLogicalOr(*bor.Left.(*expression.LogicalOrNode), env)
	case *expression.LogicalAndNode:
		left, err = resolveLogicalAnd(*bor.Left.(*expression.LogicalAndNode), env)
	default:
		return nil, fmt.Errorf("invalid left operand type for LogicalOr")
	}

	if err != nil {
		return nil, err
	}

	if bor.Right == nil {
		return left, nil
	}

	if left.IsTruthy() {
		return left, nil
	}

	right, err := resolveLogicalAnd(*bor.Right, env)
	if err != nil {
		return nil, err
	}

	return right, nil
}

func resolveLogicalAnd(band expression.LogicalAndNode, env *Env) (Value, error) {
	var left Value
	var err error

	switch band.Left.(type) {
	case *expression.LogicalAndNode:
		left, err = resolveLogicalAnd(*band.Left.(*expression.LogicalAndNode), env)
	case *expression.EqualityNode:
		left, err = resolveEquality(*band.Left.(*expression.EqualityNode), env)
	default:
		return nil, fmt.Errorf("invalid left operand type for LogicalAnd")
	}

	if err != nil {
		return nil, err
	}

	if band.Right == nil {
		return left, nil
	}

	if !left.IsTruthy() {
		return left, nil
	}

	right, err := resolveEquality(*band.Right, env)
	if err != nil {
		return nil, err
	}

	return right, nil
}

func resolveEquality(eq expression.EqualityNode, env *Env) (Value, error) {
	var left Value
	var err error

	switch eq.Left.(type) {
	case *expression.EqualityNode:
		left, err = resolveEquality(*eq.Left.(*expression.EqualityNode), env)
	case *expression.ComparisonNode:
		left, err = resolveComparison(*eq.Left.(*expression.ComparisonNode), env)
	default:
		return nil, fmt.Errorf("invalid left operand type for Equality")
	}

	if err != nil {
		return nil, err
	}

	if eq.Right == nil {
		return left, nil
	}

	right, err := resolveComparison(*eq.Right, env)
	if err != nil {
		return nil, err
	}

	if left.Type() != right.Type() {
		return nil, fmt.Errorf("type mismatch in equality comparison: %s vs %s", left.TypeName(), right.TypeName())
	}

	switch eq.Operator.Typ {
	case common.EQUAL_EQUAL:
		return &BoolValue{Value: left.Data() == right.Data()}, nil
	case common.BANG_EQUAL:
		return &BoolValue{Value: left.Data() != right.Data()}, nil
	default:
		return nil, fmt.Errorf("unknown equality operator: %s", eq.Operator.Value)
	}
}

func resolveComparison(cmp expression.ComparisonNode, env *Env) (Value, error) {
	var left Value
	var err error

	switch cmp.Left.(type) {
	case *expression.ComparisonNode:
		left, err = resolveComparison(*cmp.Left.(*expression.ComparisonNode), env)
	case *expression.TermNode:
		left, err = resolveTerm(*cmp.Left.(*expression.TermNode), env)
	default:
		return nil, fmt.Errorf("invalid left operand type for Comparison")
	}

	if err != nil {
		return nil, err
	}

	if cmp.Right == nil {
		return left, nil
	}

	right, err := resolveTerm(*cmp.Right, env)
	if err != nil {
		return nil, err
	}

	if left.Type() != right.Type() {
		return nil, fmt.Errorf("type mismatch in comparison: %s vs %s", left.TypeName(), right.TypeName())
	}

	switch cmp.Operator.Typ {
	case common.LESS:
		if left.Type() == VAL_INT {
			return &BoolValue{Value: left.(*IntValue).Value < right.(*IntValue).Value}, nil
		}
		if left.Type() == VAL_STRING {
			return &BoolValue{Value: left.(*StringValue).Value < right.(*StringValue).Value}, nil
		}
		return nil, fmt.Errorf("operator '<' not supported for type %s", left.TypeName())
	case common.LESS_EQUAL:
		if left.Type() == VAL_INT {
			return &BoolValue{Value: left.(*IntValue).Value <= right.(*IntValue).Value}, nil
		}
		if left.Type() == VAL_STRING {
			return &BoolValue{Value: left.(*StringValue).Value <= right.(*StringValue).Value}, nil
		}
		return nil, fmt.Errorf("operator '<=' not supported for type %s", left.TypeName())
	case common.GREATER:
		if left.Type() == VAL_INT {
			return &BoolValue{Value: left.(*IntValue).Value > right.(*IntValue).Value}, nil
		}
		if left.Type() == VAL_STRING {
			return &BoolValue{Value: left.(*StringValue).Value > right.(*StringValue).Value}, nil
		}
		return nil, fmt.Errorf("operator '>' not supported for type %s", left.TypeName())
	case common.GREATER_EQUAL:
		if left.Type() == VAL_INT {
			return &BoolValue{Value: left.(*IntValue).Value >= right.(*IntValue).Value}, nil
		}
		if left.Type() == VAL_STRING {
			return &BoolValue{Value: left.(*StringValue).Value >= right.(*StringValue).Value}, nil
		}
		return nil, fmt.Errorf("operator '>=' not supported for type %s", left.TypeName())
	default:
		return nil, fmt.Errorf("unknown comparison operator: %s", cmp.Operator.Value)
	}
}

func resolveTerm(term expression.TermNode, env *Env) (Value, error) {
	var left Value
	var err error

	switch term.Left.(type) {
	case *expression.TermNode:
		left, err = resolveTerm(*term.Left.(*expression.TermNode), env)
	case *expression.FactorNode:
		left, err = resolveFactor(*term.Left.(*expression.FactorNode), env)
	default:
		return nil, fmt.Errorf("invalid left operand type for Term")
	}

	if err != nil {
		return nil, err
	}

	if term.Right == nil {
		return left, nil
	}

	right, err := resolveFactor(*term.Right, env)
	if err != nil {
		return nil, err
	}

	switch term.Operator.Typ {
	case common.PLUS:
		if left.Type() != right.Type() {
			return &StringValue{Value: left.Content() + right.Content()}, nil
		}

		if left.Type() == VAL_INT {
			return &IntValue{Value: left.Data().(int) + right.Data().(int)}, nil
		}

		if left.Type() == VAL_STRING {
			return &StringValue{Value: left.Data().(string) + right.Data().(string)}, nil
		}

		return nil, fmt.Errorf("operator '+' not supported for type %s and type %s", left.TypeName(), right.TypeName())
	case common.MINUS:
		if left.Type() == VAL_INT && right.Type() == VAL_INT {
			return &IntValue{Value: left.Data().(int) - right.Data().(int)}, nil
		}
		return nil, fmt.Errorf("operator '-' not supported for type %s and type %s", left.TypeName(), right.TypeName())
	default:
		return nil, fmt.Errorf("unknown term operator: %s", term.Operator.Value)
	}
}

func resolveFactor(factor expression.FactorNode, env *Env) (Value, error) {
	var left Value
	var err error

	switch factor.Left.(type) {
	case *expression.FactorNode:
		left, err = resolveFactor(*factor.Left.(*expression.FactorNode), env)
	case *expression.UnaryNode:
		left, err = resolveUnary(*factor.Left.(*expression.UnaryNode), env)
	default:
		return nil, fmt.Errorf("invalid left operand type for Factor")
	}

	if err != nil {
		return nil, err
	}

	if factor.Right == nil {
		return left, nil
	}

	right, err := resolveUnary(*factor.Right, env)
	if err != nil {
		return nil, err
	}

	if left.Type() != right.Type() {
		return nil, fmt.Errorf("type mismatch in factor operation: %s vs %s", left.TypeName(), right.TypeName())
	}

	switch factor.Operator.Typ {
	case common.ASTERISK:
		if left.Type() == VAL_INT {
			return &IntValue{Value: left.Data().(int) * right.Data().(int)}, nil
		}
		return nil, fmt.Errorf("operator '*' not supported for type %s", left.TypeName())
	case common.SLASH:
		if left.Type() == VAL_INT {
			if right.Data().(int) == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return &IntValue{Value: left.Data().(int) / right.Data().(int)}, nil
		}
		return nil, fmt.Errorf("operator '/' not supported for type %s", left.TypeName())
	default:
		return nil, fmt.Errorf("unknown factor operator: %s", factor.Operator.Value)
	}
}

func resolveUnary(unary expression.UnaryNode, env *Env) (Value, error) {
	var right Value
	var err error

	switch unary.Right.(type) {
	case *expression.UnaryNode:
		right, err = resolveUnary(*unary.Right.(*expression.UnaryNode), env)
		if err != nil {
			return nil, err
		}
	case *expression.ArrayAccessNode:
		right, err = resolveArrayAccess(*unary.Right.(*expression.ArrayAccessNode), env)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid right operand type for Unary")
	}

	if unary.Operator == nil {
		return right, nil
	}

	switch unary.Operator.Typ {
	case common.TILDE:
		if right.Type() == VAL_INT {
			return &IntValue{Value: -right.Data().(int)}, nil
		}
		return nil, fmt.Errorf("unary '~' not supported for type %s", right.TypeName())
	case common.BANG:
		return &BoolValue{Value: !right.IsTruthy()}, nil
	default:
		return nil, fmt.Errorf("unknown unary operator: %s", unary.Operator.Value)
	}
}

func resolveArrayAccess(aa expression.ArrayAccessNode, env *Env) (Value, error) {
	var left Value
	var err error

	switch aa.Left.(type) {
	case *expression.ArrayAccessNode:
		left, err = resolveArrayAccess(*aa.Left.(*expression.ArrayAccessNode), env)
	case *expression.FunctionCallNode:
		left, err = resolveFunctionCall(*aa.Left.(*expression.FunctionCallNode), env)
	default:
		return nil, fmt.Errorf("invalid left operand type for ArrayAccess")
	}

	if err != nil {
		return nil, err
	}

	if aa.Index == nil {
		return left, nil
	}

	if left.Type() != VAL_ARRAY {
		return nil, fmt.Errorf("attempted to index a non-array value")
	}

	indexValue, err := resolveExpression(*aa.Index, env)
	if err != nil {
		return nil, err
	}

	if indexValue.Type() != VAL_INT {
		return nil, fmt.Errorf("array index must be an integer")
	}

	index := indexValue.(*IntValue).Value
	av := left.(*ArrayValue)

	if index < 0 || index >= len(av.Values) {
		return nil, fmt.Errorf("array index %d out of bounds", index)
	}

	val := av.Values[index]
	return val, nil
}

func resolveFunctionCall(fc expression.FunctionCallNode, env *Env) (Value, error) {
	callee, error := resolvePrimary(fc.Callee, env)
	if error != nil {
		return nil, error
	}

	if fc.Arguments == nil {
		return callee, nil
	}

	if callee.Type() != VAL_FUNCTION {
		return nil, fmt.Errorf("attempted to call a non-function value")
	}

	function := callee.(*FunctionValue).Function

	if len(fc.Arguments) != len(function.Parameters) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(function.Parameters), len(fc.Arguments))
	}

	args := make([]Value, 0, len(fc.Arguments))
	for _, argExpr := range fc.Arguments {
		argValue, err := resolveExpression(*argExpr, env)
		if err != nil {
			return nil, err
		}
		args = append(args, argValue)
	}

	value, err := function.Call(args, env)

	if err == nil || !errors.IsReturnErr(err) {
		return nil, err
	}

	return value, nil
}

func resolvePrimary(primary expression.Primary, env *Env) (Value, error) {
	switch p := primary.(type) {
	case *expression.TokenLiteralNode:
		token := p.Token
		switch token.Typ {
		case common.NUMBER:
			num, err := strconv.Atoi(token.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", token.Value)
			}
			return &IntValue{Value: num}, nil
		case common.LITERAL:
			return &StringValue{Value: token.Value}, nil
		case common.TRUE:
			return &BoolValue{Value: true}, nil
		case common.FALSE:
			return &BoolValue{Value: false}, nil
		case common.NONE:
			return &NoneValue{}, nil
		case common.IDENTIFIER:
			value, err := env.GetVariable(token.Value)
			if err != nil {
				return nil, err
			}
			return value, nil
		default:
			return nil, fmt.Errorf("unknown literal type: %v", token.Typ)
		}

	case *expression.GroupingExpressionNode:
		return resolveExpression(*p.Expression, env)

	case *expression.ArrayLiteralNode:
		elements := make([]Value, 0, len(p.Elements))
		for _, elemExpr := range p.Elements {
			elemValue, err := resolveExpression(*elemExpr, env)
			if err != nil {
				return nil, err
			}
			elements = append(elements, elemValue)
		}
		return &ArrayValue{Values: elements}, nil
	default:
		return nil, fmt.Errorf("unknown primary type")
	}
}
