package statement

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
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

func statementHeadline(s Statement) string {
	switch s.(type) {
	case *BlockStatement:
		return common.Colorize("BlockStatement", common.COLOR_MAGENTA)
	case *IfStatement:
		return common.Colorize("IfStatement", common.COLOR_BLUE)
	case *WhileStatement:
		return common.Colorize("WhileStatement", common.COLOR_BLUE)
	case *FunctionDef:
		return common.Colorize("FunctionDef", common.COLOR_CYAN)
	case *VarDeclaration:
		return common.Colorize("VarDeclaration", common.COLOR_GREEN)
	case *Assignment:
		return common.Colorize("Assignment", common.COLOR_GREEN)
	case *PrintStatement:
		return common.Colorize("PrintStatement", common.COLOR_RED)
	case *ReturnStatement:
		return common.Colorize("ReturnStatement", common.COLOR_CYAN)
	case *BreakStatement:
		return common.Colorize("BreakStatement", common.COLOR_CYAN)
	case *ExpressionStatement:
		return common.Colorize("ExpressionStatement", common.COLOR_YELLOW)
	default:
		return common.Colorize("UnknownStatement", common.COLOR_RED)
	}
}

type Statement interface {
	Print(start string)
}

func printStatements(start string, statements []Statement) {
	for i, stmt := range statements {

		fmt.Printf("%s%4d: %s\n", start, i+1, statementHeadline(stmt))
		stmt.Print(common.AdvanceSuffix(start + string(common.COUNTER_INDENT)))
	}
}
