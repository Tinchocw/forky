package statement

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type IfStatement struct {
	Condition *expression.Expression
	Body      *BlockStatement
	ElseIf    *ElseIfStatement
	Else      *ElseStatement
}

type ElseIfStatement struct {
	Condition *expression.Expression
	Body      *BlockStatement
	ElseIf    *ElseIfStatement
}

type ElseStatement struct {
	Body *BlockStatement
}

func (ifs IfStatement) Print(start string) {
	fmt.Printf("%s%s%s\n", start, string(common.BRANCH_CONNECTOR), common.Colorize("Condition:", common.COLOR_YELLOW))
	ifs.Condition.Print(start + string(common.SIMPLE_CONNECTOR) + string(common.LAST_CONNECTOR))

	if ifs.ElseIf != nil || ifs.Else != nil {
		fmt.Printf("%s%s%s\n", start, string(common.BRANCH_CONNECTOR), common.Colorize("Body:", common.COLOR_YELLOW))
		ifs.Body.Print(start + string(common.SIMPLE_CONNECTOR))
	} else {
		fmt.Printf("%s%s%s\n", start, string(common.LAST_CONNECTOR), common.Colorize("Body:", common.COLOR_YELLOW))
		ifs.Body.Print(start + string(common.SIMPLE_INDENT))
	}

	if ifs.ElseIf != nil {
		ifs.ElseIf.Print(start, ifs.Else != nil)
	}

	if ifs.Else != nil {
		ifs.Else.Print(start)
	}
}

func (eis ElseIfStatement) Print(start string, hasElse bool) {
	fmt.Printf("%s%s%s\n", start, string(common.BRANCH_CONNECTOR), common.Colorize("ElseIf Condition:", common.COLOR_YELLOW))
	eis.Condition.Print(start + string(common.SIMPLE_CONNECTOR) + string(common.LAST_CONNECTOR))

	if hasElse || eis.ElseIf != nil {
		fmt.Printf("%s%s%s\n", start, string(common.BRANCH_CONNECTOR), common.Colorize("Body:", common.COLOR_YELLOW))
		eis.Body.Print(start + string(common.SIMPLE_CONNECTOR))
	} else {
		fmt.Printf("%s%s%s\n", start, string(common.LAST_CONNECTOR), common.Colorize("Body:", common.COLOR_YELLOW))
		eis.Body.Print(start + string(common.SIMPLE_INDENT))
	}

	if eis.ElseIf != nil {
		eis.ElseIf.Print(start, hasElse)
	}
}

func (es ElseStatement) Print(start string) {
	fmt.Printf("%s%s%s\n", start, string(common.LAST_CONNECTOR), common.Colorize("Else Body:", common.COLOR_YELLOW))
	es.Body.Print(start + string(common.SIMPLE_INDENT))
}
