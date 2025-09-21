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

	Expression 		->	BinaryOr
	BinaryOr 		->	BinaryAnd ('or' BinaryAnd )*
	BinaryAnd 		->	Equality ('and' Equality )*
	Equality 		->	Comparison ( ( '!=' | '==' ) Comparison )*
	Comparison 		->	Term ( ( '>' | '>=' | '<' | '<=' ) Term )*
	Term 			->	Factor ( ( '-' | '+' ) Factor )*
	Factor 			->	Unary ( ( '/' | '*' ) Unary )*
	Unary 			->	( '!' | '-' ) Unary | Primary
	Primary 		->	IDENTIFIER 				|
							NUMBER 				|
							STRING 				|
							'true' 				|
							'false' 			|
							'None' 				|
							Call				|
							GroupingExpression

	Call 	-> IDENTIFIER '(' Arguments? ')'
	GroupingExpression -> '(' Expression ')'

	Parameters 		-> IDENTIFIER (',' IDENTIFIER)*
	Arguments 		-> Expression (',' Expression)*
*/

type Expression struct {
	Root BinaryOr
}

type BinaryOr struct {
	Left  BinaryAnd
	Right *BinaryOr
}

type BinaryAnd struct {
	Left  Equality
	Right *BinaryAnd
}

type Equality struct {
	Left     Comparison
	Operator *Token
	Right    *Equality
}

type Comparison struct {
	Left     Term
	Operator *Token
	Right    *Comparison
}

type Term struct {
	Left     Factor
	Operator *Token
	Right    *Term
}

type Factor struct {
	Left     Unary
	Operator *Token
	Right    *Factor
}

type Unary interface{}

type UnaryWithOperator struct {
	Operator Token
	Right    Unary
}

type Primary struct {
	Value PrimaryValue
}

type PrimaryValue interface{}

type GroupingExpression struct {
	Expression Expression
}

type Call struct {
	Callee    string
	Arguments []Expression
}

func NewLiteralExpression(value PrimaryValue) Expression {
	return Expression{
		Root: BinaryOr{
			Left: BinaryAnd{
				Left: Equality{
					Left: Comparison{
						Left: Term{
							Left: Factor{
								Left: Primary{
									Value: value,
								},
							},
						},
					},
				},
			},
		},
	}
}
