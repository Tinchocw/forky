package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type BinaryAnd struct {
	Left     *Equality
	Operator *common.Token
	Right    *BinaryAnd
}

func (ba *BinaryAnd) HasLeft() bool {
	return ba.Left != nil
}

func (ba *BinaryAnd) HasRight() bool {
	return ba.Right != nil
}

func (ba *BinaryAnd) HasOperator() bool {
	return ba.Operator != nil
}

func (ba *BinaryAnd) GetLeft() MergableNode {
	if ba.Left == nil {
		panic("GetLeft: expected *Equality")
	}
	return ba.Left
}

func (ba *BinaryAnd) GetRight() MergableNode {
	if ba.Right == nil {
		panic("GetRight: expected *BinaryAnd")
	}
	return ba.Right
}

func (ba *BinaryAnd) GetOperator() *common.Token {
	if ba.Operator == nil {
		panic("GetOperator: expected *Token")
	}
	return ba.Operator
}

func (ba *BinaryAnd) SetLeft(val MergableNode) {
	equality, ok := val.(*Equality)
	if !ok {
		panic("SetLeft: expected *Equality")
	}
	ba.Left = equality
}

func (ba *BinaryAnd) SetRight(val MergableNode) {
	binaryAnd, ok := val.(*BinaryAnd)
	if !ok {
		panic("SetRight: expected *BinaryAnd")
	}
	ba.Right = binaryAnd
}

func (ba *BinaryAnd) SetOperator(t *common.Token) {
	if t == nil {
		panic("SetOperator: expected *Token")
	}
	ba.Operator = t
}

func (ba *BinaryAnd) Print(start string) {
	if ba.skipPrinting() {
		ba.Left.Print(start)
		return
	}

	nodeName := "BinaryAnd"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	ba.Left.Print(start + string(common.BRANCH_CONNECTOR))
	ba.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (ba *BinaryAnd) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if ba.HasRight() {
			return ba.Right.ExpresionDepth(direction)
		}

		if ba.HasOperator() {
			return 0
		}
	}

	if !ba.HasLeft() {
		return 0
	}

	return ba.Left.ExpresionDepth(direction)
}

func (ba *BinaryAnd) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if ba.HasRight() {
			return ba.Right.GetSubExpression(level, direction)
		}

		if ba.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !ba.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return ba.Left.GetSubExpression(level, direction)
}

func (ba *BinaryAnd) skipPrinting() bool {
	return ba.Right == nil
}
