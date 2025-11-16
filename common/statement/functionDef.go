package statement

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type FunctionDef struct {
	Name       *string
	Parameters []string
	Body       *BlockStatement
}

func (fd FunctionDef) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(*fd.Name, common.COLOR_WHITE))

	if len(fd.Parameters) > 0 {
		fmt.Printf("%s%s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Parameters:", common.COLOR_YELLOW))
		for i, param := range fd.Parameters {
			isLast := i == len(fd.Parameters)-1
			connector := string(common.BRANCH_CONNECTOR)
			if isLast {
				connector = string(common.LAST_CONNECTOR)
			}
			fmt.Printf("%s%s\n", start+string(common.SIMPLE_CONNECTOR)+connector, common.Colorize(param, common.COLOR_WHITE))
		}
	}

	fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Body:", common.COLOR_YELLOW))
	fd.Body.Print(start + string(common.SIMPLE_INDENT))
}
