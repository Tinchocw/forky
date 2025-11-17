package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type LogicalOrNode struct {
	Left     Expression // LogicalOrNode or LogicalAndNode
	Operator *common.Token
	Right    *LogicalAndNode
}

func (lo *LogicalOrNode) Print(start string) {
	if lo.skipPrinting() {
		lo.Left.Print(start)
		return
	}

	nodeName := "LogicalOrNode"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	lo.Left.Print(start + string(common.BRANCH_CONNECTOR))
	lo.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (lo *LogicalOrNode) skipPrinting() bool {
	return lo.Right == nil
}
