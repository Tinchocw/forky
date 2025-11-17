package primaryExpression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type GroupingExpression struct {
	Expression          *expression.Expression
	StartingParenthesis bool
	ClosingParenthesis  bool
}

func (ge GroupingExpression) Print(start string) {
	fmt.Printf("%s%s\n", start, common.Colorize("GroupingExpression", common.COLOR_GREEN))
	start = common.AdvanceSuffix(start)
	ge.Expression.Print(start + string(common.LAST_CONNECTOR))
}
