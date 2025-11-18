package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type ArrayAccessNode struct {
	Left  Expression // ArrayAccessNode or FunctionCallNode
	Index *ExpressionNode
}

func (aa *ArrayAccessNode) Print(start string) {
	if aa.skipPrinting() {
		aa.Left.Print(start)
		return
	}

	nodeName := "Array Access"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_CYAN))
	start = common.AdvanceSuffix(start)
	aa.Left.Print(start + string(common.BRANCH_CONNECTOR))
	aa.Index.Print(start + string(common.LAST_CONNECTOR))
}

func (aa *ArrayAccessNode) skipPrinting() bool {
	return aa.Index == nil
}
