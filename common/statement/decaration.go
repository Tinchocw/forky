package statement

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type DeclarationStatement interface {
	Statement
}

type ArrayDeclaration struct {
	Name    *string
	Lengths []*expression.ExpressionNode
	Value   *expression.ExpressionNode
}

func (ad ArrayDeclaration) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(*ad.Name, common.COLOR_WHITE))

	if len(ad.Lengths) > 0 {
		fmt.Printf("%s%s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Length:", common.COLOR_YELLOW))
		start += string(common.SIMPLE_INDENT) + string(common.BRANCH_CONNECTOR)
		for i, lengthExpr := range ad.Lengths {
			if i == len(ad.Lengths)-1 {
				fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize(fmt.Sprintf("Dimension %d:", i+1), common.COLOR_YELLOW))
				lengthStart := start + string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
				lengthExpr.Print(lengthStart)
			} else {
				fmt.Printf("%s%s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize(fmt.Sprintf("Dimension %d:", i+1), common.COLOR_YELLOW))
				lengthStart := start + string(common.SIMPLE_INDENT) + string(common.BRANCH_CONNECTOR)
				lengthExpr.Print(lengthStart)
			}
		}
	}

	if ad.Value != nil {
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
		start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
		ad.Value.Print(start)
	}
}

type VarDeclaration struct {
	Name  *string
	Value *expression.ExpressionNode
}

func (vd VarDeclaration) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(common.BRANCH_CONNECTOR), common.Colorize("Name:", common.COLOR_YELLOW), common.Colorize(*vd.Name, common.COLOR_WHITE))

	if vd.Value != nil {
		fmt.Printf("%s%s\n", start+string(common.LAST_CONNECTOR), common.Colorize("Value:", common.COLOR_YELLOW))
		start += string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
		vd.Value.Print(start)
	}
}
