package declaration

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type ArrayDeclaration struct {
	Name    string
	Lengths []*expression.ExpressionNode
	Value   *expression.ExpressionNode
}

func (ad ArrayDeclaration) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(ad.Name, common.COLOR_WHITE))

	if len(ad.Lengths) > 0 {
		// choose header connector based on whether Value exists (so Value can be last)
		headerConn := string(common.BRANCH_CONNECTOR)
		if ad.Value == nil {
			headerConn = string(common.LAST_CONNECTOR)
		}
		fmt.Printf("%s%s\n", start+headerConn, common.Colorize("Length:", common.COLOR_YELLOW))
		// advance suffix so vertical lines propagate to children
		start = common.AdvanceSuffix(start)
		lengthsBase := start + string(common.SIMPLE_INDENT)
		for i, lengthExpr := range ad.Lengths {
			isLast := i == len(ad.Lengths)-1
			conn := string(common.BRANCH_CONNECTOR)
			if isLast {
				conn = string(common.LAST_CONNECTOR)
			}

			fmt.Printf("%s%s\n", lengthsBase+conn, common.Colorize(fmt.Sprintf("Dimension %d:", i+1), common.COLOR_YELLOW))
			// expression start: indent one more level and attach LAST_CONNECTOR so child prints correctly
			lengthStart := lengthsBase + string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
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
