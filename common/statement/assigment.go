package statement

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type Assignment interface {
	Statement
}

type VarAssignment struct {
	Name  string
	Value *expression.ExpressionNode
}

func (a VarAssignment) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(a.Name, common.COLOR_WHITE))
	fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
	start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
	a.Value.Print(start)
}

type ArrayAssignment struct {
	Name    string
	Indexes []*expression.ExpressionNode
	Value   *expression.ExpressionNode
}

func (aa ArrayAssignment) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(aa.Name, common.COLOR_WHITE))

	if len(aa.Indexes) > 0 {
		fmt.Printf("%s%s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Indexes:", common.COLOR_YELLOW))
		start += string(common.SIMPLE_INDENT) + string(common.BRANCH_CONNECTOR)
		for i, indexExpr := range aa.Indexes {
			if i == len(aa.Indexes)-1 {
				fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize(fmt.Sprintf("Index %d:", i+1), common.COLOR_YELLOW))
				indexStart := start + string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
				indexExpr.Print(indexStart)
			} else {
				fmt.Printf("%s%s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize(fmt.Sprintf("Index %d:", i+1), common.COLOR_YELLOW))
				indexStart := start + string(common.SIMPLE_INDENT) + string(common.BRANCH_CONNECTOR)
				indexExpr.Print(indexStart)
			}
		}
	}

	fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
	start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
	aa.Value.Print(start)
}
