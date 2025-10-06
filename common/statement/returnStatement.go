package statement

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type ReturnStatement struct {
	Value *expression.Expression
}

func (r ReturnStatement) Print(start string) {
	fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
	start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
	r.Value.Print(start)
}
