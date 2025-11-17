package expression

/*
	Expression 		->	LogicalOr
	LogicalOr 		->	LogicalAnd ('or' LogicalAnd )*
	LogicalAnd 		->	Equality ('and' Equality )*
	Equality 		->	Comparison ( ( '!=' | '==' ) Comparison )*
	Comparison 		->	Term ( ( '>' | '>=' | '<' | '<=' ) Term )*
	Term 			->	Factor ( ( '-' | '+' ) Factor )*
	Factor 			->	Unary ( ( '/' | '*' ) Unary )*
	Unary 			->	( '!' | '~' ) Unary | ArrAccess
	ArrAccess		->	Primary ( '[' Expression ']' )*
	Call 			->	Primary ( '(' Arguments? ')' )*
	Primary 		->	IDENTIFIER 				|
							NUMBER 				|
							STRING 				|
							'true' 				|
							'false' 			|
							'None' 				|
							GroupingExpression

	GroupingExpression -> '(' Expression ')'
	Arguments 		->	Expression ( ',' Expression )*
*/

type Expression interface {
	Print(start string)
}

type ExpressionNode struct {
	Root *LogicalOrNode
}

func (e *ExpressionNode) Print(start string) {
	e.Root.Print(start)
}
