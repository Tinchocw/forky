package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Unary interface {
	ExpressionNode
	Print(start string)
}

type UnaryWithOperator struct {
	Operator *common.Token
	Right    Unary
}

func (uwo UnaryWithOperator) Print(start string) {
	nodeName := fmt.Sprintf("Unary (%s)", uwo.Operator.FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	uwo.Right.Print(start + string(common.LAST_CONNECTOR))
}
