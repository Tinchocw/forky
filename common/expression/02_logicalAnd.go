package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type LogicalAndNode struct {
	Left     Expression
	Operator common.Token
	Right    Expression
}

func (la *LogicalAndNode) Print(start string) {
	nodeName := "LogicalAnd"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	la.Left.Print(start + string(common.BRANCH_CONNECTOR))
	la.Right.Print(start + string(common.LAST_CONNECTOR))
}
