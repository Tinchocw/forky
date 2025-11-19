package function

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
)

type ReturnStatement struct {
	Value expression.Expression
}

func (r ReturnStatement) Print(start string) {
	if r.Value != nil {
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
		start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
		r.Value.Print(start)
	}
}

func (r ReturnStatement) Headline() string {
	return common.Colorize("Return Statement", common.COLOR_CYAN)
}
