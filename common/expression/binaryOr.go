package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type BinaryOr struct {
	Left     *BinaryAnd
	Operator *common.Token
	Right    *BinaryOr
}

// func (bo BinaryOr) IsComplete() bool {
// 	return bo.HasLeft() && bo.HasOperator() && bo.HasRight()
// }

func (bo BinaryOr) HasLeft() bool {
	return bo.Left != nil
}

func (bo BinaryOr) HasRight() bool {
	return bo.Right != nil
}

func (bo BinaryOr) HasOperator() bool {
	return bo.Operator != nil
}

func (bo *BinaryOr) GetLeft() MergableNode {
	return bo.Left
}

func (bo *BinaryOr) GetRight() MergableNode {
	return bo.Right
}

func (bo *BinaryOr) GetOperator() *common.Token {
	if bo.Operator == nil {
		panic("GetOperator: expected *common.Token")
	}

	return bo.Operator
}

func (bo *BinaryOr) SetLeft(val MergableNode) {
	binaryAnd, ok := val.(*BinaryAnd)
	if !ok {
		panic("SetLeft: expected *BinaryAnd")
	}

	bo.Left = binaryAnd
}

func (bo *BinaryOr) SetRight(val MergableNode) {
	if bo == nil {
		return
	}
	bo.Right = val.(*BinaryOr)
}

func (bo *BinaryOr) SetOperator(op *common.Token) {
	if bo == nil {
		return
	}
	bo.Operator = op
}

func (bo *BinaryOr) Print(start string) {
	if bo.skipPrinting() {
		bo.Left.Print(start)
		return
	}

	nodeName := "BinaryOr"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	bo.Left.Print(start + string(common.BRANCH_CONNECTOR))
	bo.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (bo *BinaryOr) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if bo.HasRight() {
			return bo.Right.ExpresionDepth(direction)
		}

		if bo.HasOperator() {
			return 0
		}
	}

	if !bo.HasLeft() {
		return 0
	}

	return bo.Left.ExpresionDepth(direction)
}

func (bo *BinaryOr) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if bo.HasRight() {
			return bo.Right.GetSubExpression(level, direction)
		}

		if bo.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !bo.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return bo.Left.GetSubExpression(level, direction)
}

func (bo *BinaryOr) skipPrinting() bool {
	return bo.Right == nil
}
