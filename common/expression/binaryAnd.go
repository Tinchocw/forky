package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type BinaryAnd struct {
	Left     ExpressionNode
	Operator *common.Token
	Right    *Equality
}

func (ba *BinaryAnd) Print(start string) {
	if ba.skipPrinting() {
		ba.Left.Print(start)
		return
	}

	nodeName := "BinaryAnd"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	ba.Left.Print(start + string(common.BRANCH_CONNECTOR))
	ba.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (ba *BinaryAnd) skipPrinting() bool {
	return ba.Right == nil
}
