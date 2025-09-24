package common

import (
	"fmt"
)

/*
	Program			-> 	Statements*
	Statements		-> 	BlockStatement 			|
							IfStatement 		|
							WhileStatement 		|
							BreakStatement		|
							FunctionDef 		|
							ReturnStatement		|
							VarDeclaration 		|
							Assignment 			|
							PrintStatement 		|
							ExpressionStatement


	BlockStatement	-> '{' Statements * '}'
	IfStatement 	-> 'if' '(' Expression ')' BlockStatement
						( 'else' 'if' '(' Expression ')' BlockStatement )*
						( 'else' BlockStatement )?
	WhileStatement 	-> 'while' '(' Expression ')' BlockStatement
	BreakStatement  -> 'break' ';'
	FunctionDef 	-> 'func' IDENTIFIER '(' Parameters? ')' BlockStatement
	Return 			-> 'return' Expression ';'
	VarDeclaration 	-> 'var' IDENTIFIER ( '=' Expression )? ';'
	Assignment 		-> IDENTIFIER '=' Expression ';'
	PrintStatement 	-> 'print' '(' Expression ')' ';'
	ExpressionStatement -> Expression ';'

	Arguments 		-> Expression (',' Expression)*
*/

// Removed old tree printing helper functions - now using the unified tree system

func statementHeadline(s Statement) string {
	switch s.(type) {
	case BlockStatement:
		return Colorize("BlockStatement", COLOR_MAGENTA)
	case IfStatement:
		return Colorize("IfStatement", COLOR_BLUE)
	case WhileStatement:
		return Colorize("WhileStatement", COLOR_BLUE)
	case FunctionDef:
		return Colorize("FunctionDef", COLOR_CYAN)
	case VarDeclaration:
		return Colorize("VarDeclaration", COLOR_GREEN)
	case Assignment:
		return Colorize("Assignment", COLOR_GREEN)
	case PrintStatement:
		return Colorize("PrintStatement", COLOR_RED)
	case ReturnStatement:
		return Colorize("ReturnStatement", COLOR_CYAN)
	case BreakStatement:
		return Colorize("BreakStatement", COLOR_CYAN)
	case ExpressionStatement:
		return Colorize("ExpressionStatement", COLOR_YELLOW)
	default:
		return Colorize("UnknownStatement", COLOR_RED)
	}
}

func printBlock(start string, statements []Statement) {
	for i, stmt := range statements {

		fmt.Printf("%s%4d: %s\n", start, i+1, statementHeadline(stmt))
		stmt.Print(advanceSuffix(start + string(COUNTER_INDENT)))
	}
}

type Program struct {
	Statements []Statement
}

func (p Program) Print() {
	printBlock("", p.Statements)
}

type Statement interface {
	Print(start string)
}

type BlockStatement struct {
	Statements []Statement
}

func (bs BlockStatement) Print(start string) {
	printBlock(start, bs.Statements)
}

type IfStatement struct {
	Condition Expression
	Body      BlockStatement
	ElseIf    *ElseIfStatement
	Else      *ElseStatement
}

type ElseIfStatement struct {
	Condition Expression
	Body      BlockStatement
	ElseIf    *ElseIfStatement
}

type ElseStatement struct {
	Body BlockStatement
}

func (ifs IfStatement) Print(start string) {
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

func (eis ElseIfStatement) Print(start string, hasElse bool) {
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

func (es ElseStatement) Print(start string) {
	fmt.Printf("%s%s%s\n", start, string(LAST_CONNECTOR), Colorize("Else Body:", COLOR_YELLOW))
	es.Body.Print(start + string(SIMPLE_INDENT))
}

type WhileStatement struct {
	Condition Expression
	Body      BlockStatement
}

func (ws WhileStatement) Print(start string) {
	fmt.Printf("%s%s%s\n", start, string(BRANCH_CONNECTOR), Colorize("Condition:", COLOR_YELLOW))
	ws.Condition.Print(start + string(SIMPLE_CONNECTOR) + string(LAST_CONNECTOR))

	fmt.Printf("%s%s%s\n", start, string(LAST_CONNECTOR), Colorize("Body:", COLOR_YELLOW))
	ws.Body.Print(start + string(SIMPLE_INDENT))
}

type FunctionDef struct {
	Name       string
	Parameters []string
	Body       BlockStatement
}

func (fd FunctionDef) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(BRANCH_CONNECTOR), Colorize("Name:", COLOR_YELLOW), Colorize(fd.Name, COLOR_WHITE))

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

type VarDeclaration struct {
	Name  string
	Value Expression
}

func (vd VarDeclaration) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(BRANCH_CONNECTOR), Colorize("Name:", COLOR_YELLOW), Colorize(vd.Name, COLOR_WHITE))
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Value:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	vd.Value.Print(start)
}

type Assignment struct {
	Name  string
	Value Expression
}

func (a Assignment) Print(start string) {
	fmt.Printf("%s%s %s\n", start+string(BRANCH_CONNECTOR), Colorize("Name:", COLOR_YELLOW), Colorize(a.Name, COLOR_WHITE))
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Value:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	a.Value.Print(start)
}

type PrintStatement struct {
	Value Expression
}

func (ps PrintStatement) Print(start string) {
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Value:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	ps.Value.Print(start)
}

type ReturnStatement struct {
	Value Expression
}

func (r ReturnStatement) Print(start string) {
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Value:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	r.Value.Print(start)
}

type BreakStatement struct{}

func (bs BreakStatement) Print(start string) {
	fmt.Printf("%s%s\n", start, Colorize("BreakStatement", COLOR_CYAN))
}

type ExpressionStatement struct {
	Expression Expression
}

func (es ExpressionStatement) Print(start string) {
	fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize("Expression:", COLOR_YELLOW))
	start += string(SIMPLE_INDENT) + string(LAST_CONNECTOR)
	es.Expression.Print(start)
}

func PrintProgram(program Program) {
	fmt.Println()
	fmt.Printf("%s\n", Title("PROGRAM (Tree View)"))
	program.Print()
	fmt.Printf("%s\n", Title("End of PROGRAM"))
	fmt.Println()
}
