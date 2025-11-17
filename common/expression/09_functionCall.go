package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
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
	fc.Callee.Print(start + string(common.BRANCH_CONNECTOR))

	if len(fc.Arguments) > 0 {
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
	}
}

func (fc FunctionCallNode) skipPrinting() bool {
	return fc.Callee == nil
}
