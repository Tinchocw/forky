package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type FactorNode struct {
	Left     Expression // FactorNode or Unary
	Operator *common.Token
	Right    *UnaryNode
}

func (f *FactorNode) Print(start string) {
	if f.skipPrinting() {
		f.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Factor (%s)", f.Operator.FriendlyOperatorName())

	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	f.Left.Print(start + string(common.BRANCH_CONNECTOR))
	f.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (f *FactorNode) skipPrinting() bool {
	return f.Operator == nil && f.Right == nil
}
