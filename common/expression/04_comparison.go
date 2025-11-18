package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type ComparisonNode struct {
	Left     Expression // ComparisonNode or TermNode
	Operator *common.Token
	Right    *TermNode
}

func (c *ComparisonNode) Print(start string) {
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

func (c *ComparisonNode) skipPrinting() bool {
	return c.Operator == nil && c.Right == nil
}
