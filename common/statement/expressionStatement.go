package statement

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type ExpressionStatement struct {
	Expression *expression.Expression
}

func (es ExpressionStatement) Print(start string) {
	fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Expression:", common.COLOR_YELLOW))
	start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
	es.Expression.Print(start)
}
