package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type EqualityNode struct {
	Left     Expression // EqualityNode or ComparisonNode
	Operator *common.Token
	Right    *ComparisonNode
}

func (eq *EqualityNode) Print(start string) {
	if eq.skipPrinting() {
		eq.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Equality (%s)", eq.Operator.FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	eq.Left.Print(start + string(common.BRANCH_CONNECTOR))
	eq.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (eq *EqualityNode) skipPrinting() bool {
	return eq.Operator == nil && eq.Right == nil
}
