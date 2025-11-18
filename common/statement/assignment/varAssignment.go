package assignment

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type VarAssignment struct {
	Name  string
	Value *expression.ExpressionNode
}

func (a VarAssignment) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(a.Name, common.COLOR_WHITE))
	fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
	start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
	a.Value.Print(start)
}

func (a VarAssignment) Headline() string {
	return common.Colorize("Var Assignment", common.COLOR_GREEN)
}
