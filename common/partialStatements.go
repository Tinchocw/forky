package common

import "fmt"

type PartialProgram struct {
	Statements []IncompleteStatement
}

func (p PartialProgram) Print() {
}

type IncompleteStatement interface {
	Print(start string)
}

type IncompleteBlockStatement struct {
	Statements []IncompleteStatement
}

func (bs IncompleteBlockStatement) Print(start string) {

}

type IncompleteIfStatement struct {
	Condition *IncompleteExpression
	Body      *IncompleteBlockStatement
	ElseIf    *IncompleteElseIfStatement
	Else      *IncompleteElseStatement
}

type IncompleteElseIfStatement struct {
	Condition *IncompleteExpression
	Body      *IncompleteBlockStatement
	ElseIf    *IncompleteElseIfStatement
}

type IncompleteElseStatement struct {
	Body *IncompleteBlockStatement
}

func (ifs IncompleteIfStatement) Print(start string) {
	fmt.Printf("%s%s%s\n", start, string(BRANCH_CONNECTOR), Colorize("Condition:", COLOR_YELLOW))
	ifs.Condition.Print(start + string(SIMPLE_CONNECTOR) + string(LAST_CONNECTOR))

	if ifs.ElseIf != nil || ifs.Else != nil {
		fmt.Printf("%s%s%s\n", start, string(BRANCH_CONNECTOR), Colorize("Body:", COLOR_YELLOW))
		ifs.Body.Print(start + string(SIMPLE_CONNECTOR))
	} else {
		fmt.Printf("%s%s%s\n", start, string(LAST_CONNECTOR), Colorize("Body:", COLOR_YELLOW))
		ifs.Body.Print(start + string(SIMPLE_INDENT))
	}

	if ifs.ElseIf != nil {
		ifs.ElseIf.Print(start, ifs.Else != nil)
	}

	if ifs.Else != nil {
		ifs.Else.Print(start)
	}
}

func (eis IncompleteElseIfStatement) Print(start string, hasElse bool) {
	fmt.Printf("%s%s%s\n", start, string(BRANCH_CONNECTOR), Colorize("ElseIf Condition:", COLOR_YELLOW))
	eis.Condition.Print(start + string(SIMPLE_CONNECTOR) + string(LAST_CONNECTOR))

	if hasElse || eis.ElseIf != nil {
		fmt.Printf("%s%s%s\n", start, string(BRANCH_CONNECTOR), Colorize("Body:", COLOR_YELLOW))
		eis.Body.Print(start + string(SIMPLE_CONNECTOR))
	} else {
		fmt.Printf("%s%s%s\n", start, string(LAST_CONNECTOR), Colorize("Body:", COLOR_YELLOW))
		eis.Body.Print(start + string(SIMPLE_INDENT))
	}

	if eis.ElseIf != nil {
		eis.ElseIf.Print(start, hasElse)
	}
}

func (es IncompleteElseStatement) Print(start string) {
	fmt.Printf("%s%s%s\n", start, string(LAST_CONNECTOR), Colorize("Else Body:", COLOR_YELLOW))
	es.Body.Print(start + string(SIMPLE_INDENT))
}

type IncompleteWhileStatement struct {
	Condition *IncompleteExpression
	Body      *IncompleteBlockStatement
}

func (ws IncompleteWhileStatement) Print(start string) {
	fmt.Printf("%s%s%s\n", start, string(BRANCH_CONNECTOR), Colorize("Condition:", COLOR_YELLOW))
	ws.Condition.Print(start + string(SIMPLE_CONNECTOR) + string(LAST_CONNECTOR))

	fmt.Printf("%s%s%s\n", start, string(LAST_CONNECTOR), Colorize("Body:", COLOR_YELLOW))
	ws.Body.Print(start + string(SIMPLE_INDENT))
}

type IncompleteFunctionDef struct {
	Name       *string
	Parameters []string
	Body       *IncompleteBlockStatement
}

func (fd IncompleteFunctionDef) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(BRANCH_CONNECTOR), Colorize("Name:", COLOR_YELLOW), Colorize(*fd.Name, COLOR_WHITE))

	if len(fd.Parameters) > 0 {
		fmt.Printf("%s%s\n", start+string(BRANCH_CONNECTOR), Colorize("Parameters:", COLOR_YELLOW))
		for i, param := range fd.Parameters {
			isLast := i == len(fd.Parameters)-1
			connector := string(BRANCH_CONNECTOR)
			if isLast {
				connector = string(LAST_CONNECTOR)
			}
			fmt.Printf("%s%s\n", start+string(SIMPLE_CONNECTOR)+connector, Colorize(param, COLOR_WHITE))
		}
	}

	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Body:", COLOR_YELLOW))
	fd.Body.Print(start + string(SIMPLE_INDENT))
}

type IncompleteVarDeclaration struct {
	Name  *string
	Value *IncompleteExpression
}

func (vd IncompleteVarDeclaration) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(BRANCH_CONNECTOR), Colorize("Name:", COLOR_YELLOW), Colorize(*vd.Name, COLOR_WHITE))
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Value:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	vd.Value.Print(start)
}

type IncompleteAssignment struct {
	Name  *string
	Value *IncompleteExpression
}

func (a IncompleteAssignment) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(BRANCH_CONNECTOR), Colorize("Name:", COLOR_YELLOW), Colorize(*a.Name, COLOR_WHITE))
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Value:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	a.Value.Print(start)
}

type IncompletePrintStatement struct {
	Value *IncompleteExpression
}

func (ps IncompletePrintStatement) Print(start string) {
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Value:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	ps.Value.Print(start)
}

type IncompleteReturnStatement struct {
	Value *IncompleteExpression
}

func (r IncompleteReturnStatement) Print(start string) {
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Value:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	r.Value.Print(start)
}

type IncompleteBreakStatement struct{}

func (bs IncompleteBreakStatement) Print(start string) {
	fmt.Printf("%s%s\n", start, Colorize("BreakStatement", COLOR_CYAN))
}

type IncompleteExpressionStatement struct {
	Expression *IncompleteExpression
}

func (es IncompleteExpressionStatement) Print(start string) {
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Expression:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	es.Expression.Print(start)
}

func PrintPartialProgram(program PartialProgram) {
	fmt.Println()
	fmt.Printf("%s\n", Title("PROGRAM (Tree View)"))
	program.Print()
	fmt.Printf("%s\n", Title("End of PROGRAM"))
	fmt.Println()
}
