package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type FunctionCallNode struct {
	Callee    Expression // FunctionCallNode or PrimaryNode
	Arguments []*ExpressionNode
}

func (fc FunctionCallNode) Print(start string) {
	if fc.skipPrinting() {
		fc.Callee.Print(start)
		return
	}

	nodeName := "Function Call"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_GREEN))
	start = common.AdvanceSuffix(start)

	// If there are arguments, print callee with a branch connector so the
	// Arguments header can be the last connector. If there are no
	// arguments, print callee with the last connector.
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

func (fc FunctionCallNode) skipPrinting() bool {
	return fc.Arguments == nil
}
