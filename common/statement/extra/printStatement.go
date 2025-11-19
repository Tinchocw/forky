package extra

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
)

type PrintStatement struct {
	Value expression.Expression
}

func (ps PrintStatement) Print(start string) {
	if ps.Value != nil {
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
		start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
		ps.Value.Print(start)
	}
}

func (ps PrintStatement) Headline() string {
	return common.Colorize("Print Statement", common.COLOR_RED)
}
