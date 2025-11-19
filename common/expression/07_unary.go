package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type UnaryNode struct {
	Operator common.Token
	Right    Expression
}

func (u UnaryNode) Print(start string) {
	nodeName := fmt.Sprintf("Unary (%s)", u.Operator.FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	u.Right.Print(start + string(common.LAST_CONNECTOR))
}
