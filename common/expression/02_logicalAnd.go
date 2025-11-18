package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type LogicalAndNode struct {
	Left     Expression // LogicalAndNode or EqualityNode
	Operator *common.Token
	Right    *EqualityNode
}

func (la *LogicalAndNode) Print(start string) {
	if la.skipPrinting() {
		la.Left.Print(start)
		return
	}

	nodeName := "LogicalAnd"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	la.Left.Print(start + string(common.BRANCH_CONNECTOR))
	la.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (la *LogicalAndNode) skipPrinting() bool {
	return la.Right == nil
}
