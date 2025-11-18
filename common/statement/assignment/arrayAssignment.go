package assignment

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
)

type ArrayAssignment struct {
	Name    string
	Indexes []*expression.ExpressionNode
	Value   *expression.ExpressionNode
}

func (aa ArrayAssignment) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(aa.Name, common.COLOR_WHITE))

	if len(aa.Indexes) > 0 {
		headerConn := string(common.BRANCH_CONNECTOR)
		if aa.Value == nil {
			headerConn = string(common.LAST_CONNECTOR)
		}
		fmt.Printf("%s%s\n", start+headerConn, common.Colorize("Indexes:", common.COLOR_YELLOW))

		start = common.AdvanceSuffix(start)
		indexesBase := start + string(common.SIMPLE_CONNECTOR)
		for i, indexExpr := range aa.Indexes {
			isLast := i == len(aa.Indexes)-1
			conn := string(common.BRANCH_CONNECTOR)
			nextLevel := string(common.SIMPLE_CONNECTOR)
			if isLast {
				conn = string(common.LAST_CONNECTOR)
				nextLevel = string(common.SIMPLE_INDENT)
			}

			fmt.Printf("%s%s\n", indexesBase+conn, common.Colorize(fmt.Sprintf("Index %d:", i+1), common.COLOR_YELLOW))

			indexStart := indexesBase + nextLevel + string(common.LAST_CONNECTOR)
			if indexExpr != nil {
				indexExpr.Print(indexStart)
			}
		}
	}

	if aa.Value != nil {
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
		valueStart := start + string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
		aa.Value.Print(valueStart)
	}
}

func (aa ArrayAssignment) Headline() string {
	return common.Colorize("Array Assignment", common.COLOR_GREEN)
}
