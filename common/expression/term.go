package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Term struct {
	Left     *Factor
	Operator *common.Token
	Right    *Term
}

func (t *Term) HasLeft() bool {
	return t.Left != nil
}

func (t *Term) HasRight() bool {
	return t.Right != nil
}

func (t *Term) HasOperator() bool {
	return t.Operator != nil
}

func (t *Term) GetLeft() MergableNode {
	if t.Left == nil {
		panic("GetLeft: expected *Factor")
	}
	return t.Left
}

func (t *Term) GetRight() MergableNode {
	if t.Right == nil {
		panic("GetRight: expected *Term")
	}
	return t.Right
}

func (t *Term) GetOperator() *common.Token {
	if t.Operator == nil {
		panic("GetOperator: expected *Token")
	}
	return t.Operator
}

func (t *Term) SetLeft(val MergableNode) {
	factor, ok := val.(*Factor)
	if !ok {
		panic("SetLeft: expected *Factor")
	}
	t.Left = factor
}

func (t *Term) SetRight(val MergableNode) {
	term, ok := val.(*Term)
	if !ok {
		panic("SetRight: expected *Term")
	}
	t.Right = term
}

func (t *Term) SetOperator(op *common.Token) {
	if op == nil {
		panic("SetOperator: expected *Token")
	}
	t.Operator = op
}

func (t *Term) Print(start string) {
	if t.skipPrinting() {
		t.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Term (%s)", t.GetOperator().FriendlyOperatorName())

	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	t.Left.Print(start + string(common.BRANCH_CONNECTOR))
	t.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (t *Term) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if t.HasRight() {
			return t.Right.ExpresionDepth(direction)
		}

		if t.HasOperator() {
			return 0
		}
	}

	if !t.HasLeft() {
		return 0
	}

	return t.Left.ExpresionDepth(direction)
}

func (t *Term) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if t.HasRight() {
			return t.Right.GetSubExpression(level, direction)
		}

		if t.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !t.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return t.Left.GetSubExpression(level, direction)
}

func (t *Term) skipPrinting() bool {
	return t.Operator == nil && t.Right == nil
}
