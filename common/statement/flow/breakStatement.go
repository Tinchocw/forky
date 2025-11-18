package flow

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type BreakStatement struct{}

func (bs BreakStatement) Print(start string) {
	fmt.Printf("%s%s\n", start, common.Colorize("BreakStatement", common.COLOR_CYAN))
}

func (bs BreakStatement) Headline() string {
	return common.Colorize("Break Statement", common.COLOR_CYAN)
}
