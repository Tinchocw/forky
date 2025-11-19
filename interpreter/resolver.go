package interpreter

import (
	"fmt"
	"strconv"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
	"github.com/Tinchocw/forky/interpreter/errors"
)

func resolveExpression(expr expression.Expression, env *Env) (Value, error) {
	switch e := expr.(type) {
	case *expression.LogicalOrNode:
		return resolveLogicalOr(*e, env)
	case *expression.LogicalAndNode:
		return resolveLogicalAnd(*e, env)
	case *expression.EqualityNode:
		return resolveEquality(*e, env)
	case *expression.ComparisonNode:
		return resolveComparison(*e, env)
	case *expression.TermNode:
		return resolveTerm(*e, env)
	case *expression.FactorNode:
		return resolveFactor(*e, env)
	case *expression.UnaryNode:
		return resolveUnary(*e, env)
	case *expression.ArrayAccessNode:
		return resolveArrayAccess(*e, env)
	case *expression.FunctionCallNode:
		return resolveFunctionCall(*e, env)
	case expression.Primary:
		return resolvePrimary(e, env)
	default:
		return nil, fmt.Errorf("unknown expression type")
	}
}

func resolveLogicalOr(bor expression.LogicalOrNode, env *Env) (Value, error) {
	left, err := resolveExpression(bor.Left, env)
	if err != nil {
		return nil, err
	}

	if left.IsTruthy() {
		return left, nil
	}

	right, err := resolveExpression(bor.Right, env)
	if err != nil {
		return nil, err
	}

	return right, nil
}

func resolveLogicalAnd(band expression.LogicalAndNode, env *Env) (Value, error) {
	left, err := resolveExpression(band.Left, env)
	if err != nil {
		return nil, err
	}

	if !left.IsTruthy() {
		return left, nil
	}

	right, err := resolveExpression(band.Right, env)
	if err != nil {
		return nil, err
	}

	return right, nil
}

func resolveEquality(eq expression.EqualityNode, env *Env) (Value, error) {
	left, err := resolveExpression(eq.Left, env)
	if err != nil {
		return nil, err
	}

	right, err := resolveExpression(eq.Right, env)
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
	left, err := resolveExpression(cmp.Left, env)
	if err != nil {
		return nil, err
	}

	right, err := resolveExpression(cmp.Right, env)
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
	left, err := resolveExpression(term.Left, env)
	if err != nil {
		return nil, err
	}

	right, err := resolveExpression(term.Right, env)
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
	left, err := resolveExpression(factor.Left, env)
	if err != nil {
		return nil, err
	}

	right, err := resolveExpression(factor.Right, env)
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
	right, err := resolveExpression(unary.Right, env)
	if err != nil {
		return nil, err
	}

	switch unary.Operator.Typ {
	case common.PLUS:
		if right.Type() == VAL_INT {
			return &IntValue{Value: right.Data().(int)}, nil
		}
		return nil, fmt.Errorf("unary '+' not supported for type %s", right.TypeName())
	case common.MINUS:
		if right.Type() == VAL_INT {
			return &IntValue{Value: -right.Data().(int)}, nil
		}
		return nil, fmt.Errorf("unary '-' not supported for type %s", right.TypeName())
	case common.BANG:
		return &BoolValue{Value: !right.IsTruthy()}, nil
	default:
		return nil, fmt.Errorf("unknown unary operator: %s", unary.Operator.Value)
	}
}

func resolveArrayAccess(aa expression.ArrayAccessNode, env *Env) (Value, error) {
	left, err := resolveExpression(aa.Left, env)
	if err != nil {
		return nil, err
	}

	if left.Type() != VAL_ARRAY {
		return nil, fmt.Errorf("attempted to index a non-array value")
	}

	indexValue, err := resolveExpression(aa.Index, env)
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
	callee, err := resolveExpression(fc.Callee, env)
	if err != nil {
		return nil, err
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
		argValue, err := resolveExpression(argExpr, env)
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
		return resolveTokenLiteral(*p, env)

	case *expression.GroupingExpressionNode:
		return resolveExpression(p.Expression, env)

	case *expression.ArrayLiteralNode:
		return resolveArrayLiteral(*p, env)

	default:
		return nil, fmt.Errorf("unknown primary type")
	}
}

func resolveTokenLiteral(tl expression.TokenLiteralNode, env *Env) (Value, error) {
	token := tl.Token
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
}

func resolveArrayLiteral(al expression.ArrayLiteralNode, env *Env) (Value, error) {
	elements := make([]Value, len(al.Elements))

	for i, elemExpr := range al.Elements {
		elemValue, err := resolveExpression(elemExpr, env)
		if err != nil {
			return nil, err
		}
		elements[i] = elemValue
	}

	return &ArrayValue{Values: elements}, nil

}
