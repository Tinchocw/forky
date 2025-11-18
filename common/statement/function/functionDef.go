package function

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/statement/block"
)

type FunctionDef struct {
	Name       *string
	Parameters []string
	Body       *block.BlockStatement
}

func (fd FunctionDef) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(*fd.Name, common.COLOR_WHITE))

	if len(fd.Parameters) > 0 {
		fmt.Printf("%s%s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Parameters:", common.COLOR_YELLOW))
		start = common.AdvanceSuffix(start)
		paramsBase := start + string(common.SIMPLE_CONNECTOR)
		for i, param := range fd.Parameters {
			isLast := i == len(fd.Parameters)-1
			conn := string(common.BRANCH_CONNECTOR)
			if isLast {
				conn = string(common.LAST_CONNECTOR)
			}

			fmt.Printf("%s%s %s\n", paramsBase+conn, common.Colorize(fmt.Sprintf("Parameter %d:", i+1), common.COLOR_YELLOW), common.Colorize(param, common.COLOR_WHITE))
		}
	}

	fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Body:", common.COLOR_YELLOW))
	fd.Body.Print(start + string(common.SIMPLE_INDENT))
}

func (fd FunctionDef) Headline() string {
	return common.Colorize("Function Definition", common.COLOR_CYAN)
}
