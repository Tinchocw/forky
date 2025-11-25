package declaration

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
)

type VarDeclaration struct {
	Name  string
	Value expression.Expression
}

func (vd VarDeclaration) Print(start string) {
	conector := string(common.BRANCH_CONNECTOR)
	if vd.Value == nil {
		conector = string(common.LAST_CONNECTOR)
	}

	fmt.Printf("%s%s %s\n", start+conector, common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(vd.Name, common.COLOR_WHITE))

	if vd.Value != nil {
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
		start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
		vd.Value.Print(start)
	}
}

func (vd VarDeclaration) Headline() string {
	return common.Colorize("Var Declaration", common.COLOR_GREEN)
}
