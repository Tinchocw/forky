package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type LogicalOrNode struct {
	Left     Expression
	Operator common.Token
	Right    Expression
}

func (lo *LogicalOrNode) Print(start string) {
	nodeName := "LogicalOrNode"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	lo.Left.Print(start + string(common.BRANCH_CONNECTOR))
	lo.Right.Print(start + string(common.LAST_CONNECTOR))
}
