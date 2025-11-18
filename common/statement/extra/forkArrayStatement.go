package extra

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
	"github.com/Tinchocw/forky/common/statement/block"
)

type ForkArrayStatement struct {
	Array     *expression.ExpressionNode
	IndexName *string
	ElemName  *string
	Block     *block.BlockStatement
}

func (fas *ForkArrayStatement) Print(start string) {
	fmt.Printf("%s%s%s\n", start, string(common.BRANCH_CONNECTOR), common.Colorize("Array:", common.COLOR_YELLOW))
	fas.Array.Print(start + string(common.SIMPLE_CONNECTOR) + string(common.LAST_CONNECTOR))

	if fas.IndexName != nil {

		fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Index Name:", common.COLOR_YELLOW), common.Colorize(*fas.IndexName, common.COLOR_WHITE))

	}

	if fas.ElemName != nil {

		fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Elem Name:", common.COLOR_YELLOW), common.Colorize(*fas.ElemName, common.COLOR_WHITE))

	}

	fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Body:", common.COLOR_YELLOW))
	fas.Block.Print(start + string(common.SIMPLE_INDENT))
}

func (fas *ForkArrayStatement) Headline() string {
	return common.Colorize("Fork Array Statement", common.COLOR_CYAN)
}
