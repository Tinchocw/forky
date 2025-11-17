package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type ArrayLiteralNode struct {
	Elements []*ExpressionNode
}

func (al ArrayLiteralNode) Print(start string) {
	nodeName := "Array Literal"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_GREEN))
	start = common.AdvanceSuffix(start)

	if len(al.Elements) > 0 {
		nodeName := "Elements"
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize(nodeName, common.COLOR_GREEN))
		start += string(common.SIMPLE_INDENT)
		for i, arg := range al.Elements {
			connector := common.BRANCH_CONNECTOR
			identation := common.SIMPLE_CONNECTOR
			if i == len(al.Elements)-1 {
				connector = common.LAST_CONNECTOR
				identation = common.SIMPLE_INDENT
			}
			fmt.Printf("%sElement[%d]:\n", start+string(connector), i)
			arg.Print(start + string(identation) + string(common.LAST_CONNECTOR))
		}
	}
}
