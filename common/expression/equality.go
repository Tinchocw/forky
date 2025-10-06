package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Equality struct {
	Left     *Comparison
	Operator *common.Token
	Right    *Equality
}

func (eq *Equality) HasLeft() bool {
	return eq.Left != nil
}

func (eq *Equality) HasRight() bool {
	return eq.Right != nil
}

func (eq *Equality) HasOperator() bool {
	return eq.Operator != nil
}

func (eq *Equality) GetLeft() MergableNode {
	if eq.Left == nil {
		panic("GetLeft: expected *Comparison")
	}
	return eq.Left
}

func (eq *Equality) GetRight() MergableNode {
	if eq.Right == nil {
		panic("GetRight: expected *Equality")
	}
	return eq.Right
}

func (eq *Equality) GetOperator() *common.Token {
	if eq.Operator == nil {
		panic("GetOperator: expected *Token")
	}
	return eq.Operator
}

func (eq *Equality) SetLeft(val MergableNode) {
	comparison, ok := val.(*Comparison)
	if !ok {
		panic("SetLeft: expected *Comparison")
	}
	eq.Left = comparison
}

func (eq *Equality) SetRight(val MergableNode) {
	equality, ok := val.(*Equality)
	if !ok {
		panic("SetRight: expected *Equality")
	}
	eq.Right = equality
}

func (eq *Equality) SetOperator(t *common.Token) {
	if t == nil {
		panic("SetOperator: expected *Token")
	}
	eq.Operator = t
}

func (eq *Equality) Print(start string) {
	if eq.skipPrinting() {
		eq.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Equality (%s)", eq.GetOperator().FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	eq.Left.Print(start + string(common.BRANCH_CONNECTOR))
	eq.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (eq *Equality) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if eq.HasRight() {
			return eq.Right.ExpresionDepth(direction)
		}

		if eq.HasOperator() {
			return 0
		}
	}

	if !eq.HasLeft() {
		return 0
	}

	return eq.Left.ExpresionDepth(direction)
}

func (eq *Equality) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if eq.HasRight() {
			return eq.Right.GetSubExpression(level, direction)
		}

		if eq.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !eq.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return eq.Left.GetSubExpression(level, direction)
}

func (eq *Equality) skipPrinting() bool {
	return eq.Operator == nil && eq.Right == nil
}
