package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type EqualityNode struct {
	Left     Expression
	Operator common.Token
	Right    Expression
}

func (eq *EqualityNode) Print(start string) {
	nodeName := fmt.Sprintf("Equality (%s)", eq.Operator.FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	eq.Left.Print(start + string(common.BRANCH_CONNECTOR))
	eq.Right.Print(start + string(common.LAST_CONNECTOR))
}
