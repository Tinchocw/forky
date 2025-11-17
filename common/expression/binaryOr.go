package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type BinaryOr struct {
	Left     ExpressionNode
	Operator *common.Token
	Right    *BinaryAnd
}

func (bo *BinaryOr) Print(start string) {
	if bo.skipPrinting() {
		bo.Left.Print(start)
		return
	}

	nodeName := "BinaryOr"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	bo.Left.Print(start + string(common.BRANCH_CONNECTOR))
	bo.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (bo *BinaryOr) skipPrinting() bool {
	return bo.Right == nil
}
