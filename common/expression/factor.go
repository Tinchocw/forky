package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Factor struct {
	Left     Unary
	Operator *common.Token
	Right    *Factor
}

func (f *Factor) HasLeft() bool {
	return f.Left != nil
}

func (f *Factor) HasRight() bool {
	return f.Right != nil
}

func (f *Factor) HasOperator() bool {
	return f.Operator != nil
}

func (f *Factor) GetLeft() MergableNode {
	if f.Left == nil {
		panic("GetLeft: expected Unary")
	}
	node, ok := f.Left.(MergableNode)
	if !ok {
		panic("GetLeft: expected MergableNode (Unary)")
	}
	return node
}

func (f *Factor) GetRight() MergableNode {
	if f.Right == nil {
		panic("GetRight: expected *Factor")
	}
	return f.Right
}

func (f *Factor) GetOperator() *common.Token {
	if f.Operator == nil {
		panic("GetOperator: expected *Token")
	}
	return f.Operator
}

func (f *Factor) SetLeft(val MergableNode) {
	unary, ok := val.(Unary)
	if !ok {
		panic("SetLeft: expected Unary")
	}
	f.Left = unary
}

func (f *Factor) SetRight(val MergableNode) {
	factor, ok := val.(*Factor)
	if !ok {
		panic("SetRight: expected *Factor")
	}
	f.Right = factor
}

func (f *Factor) SetOperator(op *common.Token) {
	if op == nil {
		panic("SetOperator: expected *Token")
	}
	f.Operator = op
}

func (f *Factor) Print(start string) {
	if f.skipPrinting() {
		f.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Factor (%s)", f.GetOperator().FriendlyOperatorName())

	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	f.Left.Print(start + string(common.BRANCH_CONNECTOR))
	f.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (f *Factor) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if f.HasRight() {
			return f.Right.ExpresionDepth(direction)
		}

		if f.HasOperator() {
			return 0
		}
	}

	if !f.HasLeft() {
		return 0
	}

	return f.Left.ExpresionDepth(direction)
}

func (f *Factor) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if f.HasRight() {
			return f.Right.GetSubExpression(level, direction)
		}

		if f.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !f.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return f.Left.GetSubExpression(level, direction)
}

func (f *Factor) skipPrinting() bool {
	return f.Operator == nil && f.Right == nil
}
