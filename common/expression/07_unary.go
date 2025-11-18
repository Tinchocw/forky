package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type UnaryNode struct {
	Operator *common.Token
	Right    Expression // UnaryNode or ArrayAccessNode
}

func (u UnaryNode) Print(start string) {
	if u.skipPrinting() {
		u.Right.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Unary (%s)", u.Operator.FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	u.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (u UnaryNode) skipPrinting() bool {
	return u.Operator == nil
}
