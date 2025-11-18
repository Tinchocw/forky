package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type GroupingExpressionNode struct {
	Expression *ExpressionNode
}

func (ge GroupingExpressionNode) Print(start string) {
	fmt.Printf("%s%s\n", start, common.Colorize("GroupingExpression", common.COLOR_GREEN))
	start = common.AdvanceSuffix(start)
	ge.Expression.Print(start + string(common.LAST_CONNECTOR))
}
