package statement

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

/*
	Program				-> 	Statements*
	Statements			-> 	BlockStatement 			|
								IfStatement 		|
								WhileStatement 		|
								BreakStatement		|
								FunctionDef 		|
								ReturnStatement		|
								VarDeclaration 		|
								Assignment 			|
								PrintStatement 		|
								ExpressionStatement


	BlockStatement		-> '{' Statements * '}'
	IfStatement 		-> 'if' '(' Expression ')' BlockStatement
							( 'else' 'if' '(' Expression ')' BlockStatement )*
							( 'else' BlockStatement )?
	WhileStatement 		-> 'while' '(' Expression ')' BlockStatement
	BreakStatement  	-> 'break' ';'
	FunctionDef 		-> 'func' IDENTIFIER '(' Parameters? ')' BlockStatement
	Return 				-> 'return' Expression ';'
	VarDeclaration 		-> 'var' IDENTIFIER ( '=' Expression )? ';'
	ArrayDeclaration	-> 'var' IDENTIFIER ( '[' Expression ']' )+
	Assignment 			-> 'set' IDENTIFIER '=' Expression ';'
	ArrayAssignment 	-> 'set' IDENTIFIER ('[' Expression ']')+ '=' Expression ';'
	PrintStatement 		-> 'print' '(' Expression ')' ';'
	ForkStatement   	-> 'fork' BlockStatement
	ForkArrayStatement  -> 'fork' IDENTIFIER ',' IDENTIFIER BlockStatement
	ExpressionStatement -> Expression ';'

	Arguments 		-> Expression (',' Expression)*
*/

/*

fork {
	foo();
	too();
}

fork arr index,elem {
	print(elem);
	print(index);
}

*/

type Statement interface {
	Print(start string)
	Headline() string
}

func PrintStatements(start string, statements []Statement) {
	for i, stmt := range statements {
		fmt.Printf("%s%4d: %s\n", start, i+1, stmt.Headline())
		stmt.Print(common.AdvanceSuffix(start + string(common.COUNTER_INDENT)))
	}
}
