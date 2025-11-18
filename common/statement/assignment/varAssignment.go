package assignment

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
)

type VarAssignment struct {
	Name  string
	Value *expression.ExpressionNode
}

func (a VarAssignment) Print(start string) {
	conector := string(common.BRANCH_CONNECTOR)
	if a.Value == nil {
		conector = string(common.LAST_CONNECTOR)
	}

	fmt.Printf("%s%s %s\n", start+conector, common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(a.Name, common.COLOR_WHITE))

	if a.Value != nil {
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
		start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
		a.Value.Print(start)
	}
}

func (a VarAssignment) Headline() string {
	return common.Colorize("Var Assignment", common.COLOR_GREEN)
}
