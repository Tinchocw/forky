package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Unary interface {
	MergableNode
	Print(start string)
}

type UnaryWithOperator struct {
	Operator *common.Token
	Right    Unary
}

func (uwo UnaryWithOperator) Print(start string) {
	nodeName := fmt.Sprintf("Unary (%s)", uwo.GetOperator().FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	uwo.Right.Print(start + string(common.LAST_CONNECTOR))
}

//	func (uwo *UnaryWithOperator) IsComplete() bool {
//		return uwo.Operator != nil && uwo.Right != nil
//	}

func (uwo *UnaryWithOperator) HasRight() bool {
	return uwo.Right != nil
}
func (uwo *UnaryWithOperator) HasLeft() bool {
	return false
}

func (uwo *UnaryWithOperator) HasOperator() bool {
	return uwo.Operator != nil
}

func (uwo *UnaryWithOperator) GetLeft() MergableNode {
	return nil
}

func (uwo *UnaryWithOperator) GetRight() MergableNode {
	if uwo == nil || uwo.Right == nil {
		return nil
	}
	return uwo.Right
}

func (uwo *UnaryWithOperator) GetOperator() *common.Token {
	return uwo.Operator
}

func (uwo *UnaryWithOperator) SetLeft(n MergableNode) {
	panic("SetLeft: UnaryWithOperator does not have a left node")
}

func (uwo *UnaryWithOperator) SetRight(val MergableNode) {
	unary, ok := val.(Unary)
	if !ok {
		panic("SetRight: expected Unary")
	}
	uwo.Right = unary
}

func (uwo *UnaryWithOperator) SetOperator(t *common.Token) {
	if t == nil {
		panic("SetOperator: expected *Token")
	}
	uwo.Operator = t
}

func (uwo *UnaryWithOperator) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if uwo.HasRight() {
			return uwo.Right.ExpresionDepth(direction)
		}

		if uwo.HasOperator() {
			return 0
		}
	}

	if !uwo.HasLeft() {
		return 0
	}

	return uwo.GetLeft().ExpresionDepth(direction)
}

func (uwo *UnaryWithOperator) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if uwo.HasRight() {
			return uwo.Right.GetSubExpression(level, direction)
		}

		if uwo.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !uwo.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return uwo.GetLeft().GetSubExpression(level, direction)
}
