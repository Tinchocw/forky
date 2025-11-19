package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type FactorNode struct {
	Left     Expression
	Operator common.Token
	Right    Expression
}

func (f *FactorNode) Print(start string) {
	nodeName := fmt.Sprintf("Factor (%s)", f.Operator.FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	f.Left.Print(start + string(common.BRANCH_CONNECTOR))
	f.Right.Print(start + string(common.LAST_CONNECTOR))
}
