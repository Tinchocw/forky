package declaration

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
)

type ArrayDeclaration struct {
	Name    string
	Lengths []expression.Expression
	Value   expression.Expression
}

func (ad ArrayDeclaration) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(ad.Name, common.COLOR_WHITE))

	if len(ad.Lengths) > 0 {
		headerConn := string(common.BRANCH_CONNECTOR)
		if ad.Value == nil {
			headerConn = string(common.LAST_CONNECTOR)
		}
		fmt.Printf("%s%s\n", start+headerConn, common.Colorize("Lengths:", common.COLOR_YELLOW))

		start = common.AdvanceSuffix(start)
		lengthsBase := start + string(common.SIMPLE_INDENT)
		if ad.Value != nil {
			lengthsBase = start + string(common.SIMPLE_CONNECTOR)
		}
		for i, lengthExpr := range ad.Lengths {
			isLast := i == len(ad.Lengths)-1
			conn := string(common.BRANCH_CONNECTOR)
			nextLevel := string(common.SIMPLE_CONNECTOR)
			if isLast {
				conn = string(common.LAST_CONNECTOR)
				nextLevel = string(common.SIMPLE_INDENT)
			}

			fmt.Printf("%s%s\n", lengthsBase+conn, common.Colorize(fmt.Sprintf("Dimension %d:", i+1), common.COLOR_YELLOW))

			lengthStart := lengthsBase + nextLevel + string(common.LAST_CONNECTOR)
			if lengthExpr != nil {
				lengthExpr.Print(lengthStart)
			}
		}
	}

	if ad.Value != nil {
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
		valueStart := start + string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
		ad.Value.Print(valueStart)
	}
}

func (ad ArrayDeclaration) Headline() string {
	return common.Colorize("Array Declaration", common.COLOR_GREEN)
}
