package primaryExpression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type ArrayAccess struct {
	ArrayName string
	Indexes   []*expression.Expression
}

func (a *ArrayAccess) Print(start string) {
	nodeName := "Array Access"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_CYAN))
	start = common.AdvanceSuffix(start)
	fmt.Printf("%s%s%s\n", start, string(common.BRANCH_CONNECTOR), common.Colorize("Array Name: "+a.ArrayName, common.COLOR_YELLOW))

	if len(a.Indexes) > 0 {
		nodeName := "Indexes"
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize(nodeName, common.COLOR_CYAN))
		start += string(common.SIMPLE_INDENT)
		for i, idx := range a.Indexes {
			connector := common.BRANCH_CONNECTOR
			identation := common.SIMPLE_CONNECTOR
			if i == len(a.Indexes)-1 {
				connector = common.LAST_CONNECTOR
				identation = common.SIMPLE_INDENT
			}
			fmt.Printf("%sIndex[%d]:\n", start+string(connector), i)
			idx.Print(start + string(identation) + string(common.LAST_CONNECTOR))
		}
	}

}
