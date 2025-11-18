package extra

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/statement/block"
)

type ForkBlockStatement struct {
	Block *block.BlockStatement
}

func (fs ForkBlockStatement) Print(start string) {
	fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Body:", common.COLOR_YELLOW))
	fs.Block.Print(start + string(common.SIMPLE_INDENT))
}

func (fs ForkBlockStatement) Headline() string {
	return common.Colorize("Fork Statement", common.COLOR_CYAN)
}
