package flow

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
	"github.com/Tinchocw/Interprete-concurrente/common/statement/block"
)

type WhileStatement struct {
	Condition *expression.ExpressionNode
	Body      *block.BlockStatement
}

func (ws WhileStatement) Print(start string) {
	fmt.Printf("%s%s%s\n", start, string(common.BRANCH_CONNECTOR), common.Colorize("Condition:", common.COLOR_YELLOW))
	ws.Condition.Print(start + string(common.SIMPLE_CONNECTOR) + string(common.LAST_CONNECTOR))

	fmt.Printf("%s%s%s\n", start, string(common.LAST_CONNECTOR), common.Colorize("Body:", common.COLOR_YELLOW))
	ws.Body.Print(start + string(common.SIMPLE_INDENT))
}

func (ws WhileStatement) Headline() string {
	return common.Colorize("While Statement", common.COLOR_BLUE)
}
