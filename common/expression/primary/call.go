package primaryExpression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type Call struct {
	Callee    *string
	Arguments []*expression.Expression
}

func (c Call) Print(start string) {
	nodeName := "Call"
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_GREEN))
	start = common.AdvanceSuffix(start)
	fmt.Printf("%sCallee: %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize(*c.Callee, common.COLOR_WHITE))

	if len(c.Arguments) > 0 {
		nodeName := "Arguments"
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize(nodeName, common.COLOR_GREEN))
		start += string(common.SIMPLE_INDENT)
		for i, arg := range c.Arguments {
			connector := common.BRANCH_CONNECTOR
			identation := common.SIMPLE_CONNECTOR
			if i == len(c.Arguments)-1 {
				connector = common.LAST_CONNECTOR
				identation = common.SIMPLE_INDENT
			}
			fmt.Printf("%sArg[%d]:\n", start+string(connector), i)
			arg.Print(start + string(identation) + string(common.LAST_CONNECTOR))
		}
	}
}

func (c *Call) HasLeft() bool {
	return true
}

func (c *Call) HasRight() bool {
	return false
}

func (c *Call) HasOperator() bool {
	return false
}

func (c *Call) GetLeft() expression.MergableNode {
	panic("TODO")
}

func (c *Call) GetRight() expression.MergableNode {
	panic("GetRight: Call does not have a right node")
}

func (c *Call) GetOperator() *common.Token {
	panic("GetOperator: Call does not have an operator")
}

func (c *Call) SetLeft(n expression.MergableNode) {
	panic("TODO")
}

func (c *Call) SetRight(n expression.MergableNode) {
	panic("SetRight: Call does not have a right node")
}

func (c *Call) SetOperator(t *common.Token) {
	panic("SetOperator: Call does not have an operator")
}

func (c *Call) ExpresionDepth(direction expression.Direction) int {
	return 1
}

func (c *Call) GetSubExpression(level int, direction expression.Direction) expression.MergableNode {
	panic("TODO")
}
