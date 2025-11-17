package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Comparison struct {
	Left     ExpressionNode
	Operator *common.Token
	Right    *Term
}

func (c *Comparison) Print(start string) {
	if c.skipPrinting() {
		c.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Comparison (%s)", c.Operator.FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	c.Left.Print(start + string(common.BRANCH_CONNECTOR))
	c.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (c *Comparison) skipPrinting() bool {
	return c.Operator == nil && c.Right == nil
}
