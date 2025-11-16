package expression

import "github.com/Tinchocw/Interprete-concurrente/common"

/*
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
*/

type Expression struct {
	Root *BinaryOr
}

func (e *Expression) HasLeft() bool {
	return e.Root != nil
}

func (e *Expression) HasRight() bool {
	return false
}

func (e *Expression) HasOperator() bool {
	return false
}

func (e *Expression) GetLeft() MergableNode {
	if e.Root == nil {
		panic("GetLeft: expected *BinaryOr")
	}
	return e.Root
}

func (e *Expression) GetRight() MergableNode {
	panic("GetRight: Expression does not have a right node")
}

func (e *Expression) GetOperator() *common.Token {
	panic("GetOperator: Expression does not have an operator")
}

func (e *Expression) SetLeft(n MergableNode) {
	binaryOr, ok := n.(*BinaryOr)
	if !ok {
		panic("SetLeft: expected *BinaryOr")
	}
	e.Root = binaryOr
}

func (e *Expression) SetRight(n MergableNode) {
	panic("SetRight: Expression does not have a right node")
}

func (e *Expression) SetOperator(t *common.Token) {
	panic("SetOperator: Expression does not have an operator")
}

func (e *Expression) ExpresionDepth(direction Direction) int {
	if e.Root == nil {
		return 0
	}

	return e.Root.ExpresionDepth(direction)
}

func (e *Expression) GetSubExpression(level int, direction Direction) MergableNode {
	if level == 0 {
		return e
	}

	return e.Root.GetSubExpression(level, direction)
}
func (e *Expression) Print(start string) {
	e.Root.Print(start)
}
