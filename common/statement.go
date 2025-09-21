package common

/*
	Program			-> 	Statements * EOF
	Statements		-> 	BlockStatement 			|
							IfStatement 		|
							WhileStatement 		|
							FunctionDef 		|
							VarDeclaration 		|
							Assignment 			|
							PrintStatement 		|
							Expression

	BlockStatement	-> '{' Statements * '}'
	IfStatement 	-> 'if' '(' Expression ')' BlockStatement
						( 'else' 'if' '(' Expression ')' BlockStatement )*
						( 'else' BlockStatement )?
	WhileStatement 	-> 'while' '(' Expression ')' BlockStatement
	FunctionDef 	-> 'func' IDENTIFIER '(' Parameters? ')' BlockStatement
	VarDeclaration 	-> 'var' IDENTIFIER ( '=' Expression )? ';'
	Assignment 		-> IDENTIFIER '=' Expression ';'
	PrintStatement 	-> 'print' '(' Expression ')' ';'
	Return 			-> 'return' Expression ';'
*/

type Program struct {
	Statements []Statement
}

type Statement interface{}

type BlockStatement struct {
	Statements []Statement
}

type IfStatement struct {
	Condition Expression
	Body      BlockStatement
	Else      *IfStatement
}

type WhileStatement struct {
	Condition Expression
	Body      BlockStatement
}

type FunctionDef struct {
	Name       string
	Parameters []string
	Body       BlockStatement
}

type VarDeclaration struct {
	Name  string
	Value Expression
}

type Assignment struct {
	Name  string
	Value Expression
}

type PrintStatement struct {
	Value Expression
}

type Return struct {
	Value Expression
}

type BreakStatement struct{}
