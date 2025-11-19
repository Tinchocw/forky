package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type FunctionCallNode struct {
	Callee    Expression
	Arguments []Expression
}

func (fc FunctionCallNode) Print(start string) {
	nodeName := "Function Call"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_GREEN))
	start = common.AdvanceSuffix(start)

	if len(fc.Arguments) > 0 {
		fc.Callee.Print(start + string(common.BRANCH_CONNECTOR))

		nodeName := "Arguments"
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize(nodeName, common.COLOR_GREEN))
		start += string(common.SIMPLE_INDENT)
		for i, arg := range fc.Arguments {
			connector := common.BRANCH_CONNECTOR
			identation := common.SIMPLE_CONNECTOR
			if i == len(fc.Arguments)-1 {
				connector = common.LAST_CONNECTOR
				identation = common.SIMPLE_INDENT
			}
			fmt.Printf("%sArg[%d]:\n", start+string(connector), i)
			arg.Print(start + string(identation) + string(common.LAST_CONNECTOR))
		}
	} else {
		fc.Callee.Print(start + string(common.LAST_CONNECTOR))
	}
}
