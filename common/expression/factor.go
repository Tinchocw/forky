package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Factor struct {
	Left     ExpressionNode
	Operator *common.Token
	Right    Unary
}

func (f *Factor) Print(start string) {
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

func (f *Factor) skipPrinting() bool {
	return f.Operator == nil && f.Right == nil
}
