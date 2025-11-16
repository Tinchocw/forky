package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Comparison struct {
	Left     *Term
	Operator *common.Token
	Right    *Comparison
}

func (c *Comparison) HasLeft() bool {
	return c.Left != nil
}

func (c *Comparison) HasRight() bool {
	return c.Right != nil
}

func (c *Comparison) HasOperator() bool {
	return c.Operator != nil
}

func (c *Comparison) GetLeft() MergableNode {
	if c.Left == nil {
		panic("GetLeft: expected *Term")
	}
	return c.Left
}

func (c *Comparison) GetRight() MergableNode {
	if c.Right == nil {
		panic("GetRight: expected *Comparison")
	}
	return c.Right
}

func (c *Comparison) GetOperator() *common.Token {
	if c.Operator == nil {
		panic("GetOperator: expected *common.Token")
	}
	return c.Operator
}

func (c *Comparison) SetLeft(val MergableNode) {
	term, ok := val.(*Term)
	if !ok {
		panic("SetLeft: expected *Term")
	}
	c.Left = term
}

func (c *Comparison) SetRight(val MergableNode) {
	comparison, ok := val.(*Comparison)
	if !ok {
		panic("SetRight: expected *Comparison")
	}
	c.Right = comparison
}

func (c *Comparison) SetOperator(t *common.Token) {
	if t == nil {
		panic("SetOperator: expected *common.Token")
	}
	c.Operator = t
}

func (c *Comparison) Print(start string) {
	if c.skipPrinting() {
		c.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Comparison (%s)", c.GetOperator().FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	c.Left.Print(start + string(common.BRANCH_CONNECTOR))
	c.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (c *Comparison) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if c.HasRight() {
			return c.Right.ExpresionDepth(direction)
		}

		if c.HasOperator() {
			return 0
		}
	}

	if !c.HasLeft() {
		return 0
	}

	return c.Left.ExpresionDepth(direction)
}

func (c *Comparison) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if c.HasRight() {
			return c.Right.GetSubExpression(level, direction)
		}

		if c.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !c.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return c.Left.GetSubExpression(level, direction)
}

func (c *Comparison) skipPrinting() bool {
	return c.Operator == nil && c.Right == nil
}
