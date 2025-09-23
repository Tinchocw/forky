package common

import (
	"fmt"
	"strings"
)

/*

	Program			-> 	Statements * EOF
	Statements		-> 	BlockStatement 			|
							IfStatement 		|
							WhileStatement 		|
							FunctionDef 		|
							VarDeclaration 		|
							Assignment 			|
							PrintStat	}
}

type Unary interface{}ent 		|
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

// Helper function for simple indentation
func simpleIndent(level int) string {
	return strings.Repeat("  ", level)
}

// Helper function to create tree connectors (only for lowest level)
func leafPrefix(isLast bool) string {
	if isLast {
		return "└── "
	} else {
		return "├── "
	}
}

type Expression struct {
	Root BinaryOr
}

// Helper functions to check if nodes are "pass-through" (only have Left, no Right/Operator)
func (bo BinaryOr) isSimple() bool {
	return bo.Right == nil
}

func (ba BinaryAnd) isSimple() bool {
	return ba.Right == nil
}

func (eq Equality) isSimple() bool {
	return eq.Operator == nil && eq.Right == nil
}

func (c Comparison) isSimple() bool {
	return c.Operator == nil && c.Right == nil
}

func (t Term) isSimple() bool {
	return t.Operator == nil && t.Right == nil
}

func (f Factor) isSimple() bool {
	return f.Operator == nil && f.Right == nil
}

func (e Expression) Print(level int) {
	e.Root.Print(level)
}

type BinaryOr struct {
	Left  BinaryAnd
	Right *BinaryOr
}

func (bo BinaryOr) Print(level int) {
	// Si es simple (no tiene Right), saltear este nivel
	if bo.isSimple() {
		bo.Left.Print(level)
		return
	}

	fmt.Printf("%s%sBinaryOr%s\n", simpleIndent(level), COLOR_MAGENTA, COLOR_RESET)

	fmt.Printf("%sLeft:\n", simpleIndent(level+1))
	bo.Left.Print(level + 2)

	if bo.Right != nil {
		fmt.Printf("%sOperator: %sor%s\n", simpleIndent(level+1), COLOR_RED, COLOR_RESET)
		fmt.Printf("%sRight:\n", simpleIndent(level+1))
		bo.Right.Print(level + 2)
	}
}

type BinaryAnd struct {
	Left  Equality
	Right *BinaryAnd
}

func (ba BinaryAnd) Print(level int) {
	// Si es simple (no tiene Right), saltear este nivel
	if ba.isSimple() {
		ba.Left.Print(level)
		return
	}

	fmt.Printf("%s%sBinaryAnd%s\n", simpleIndent(level), COLOR_MAGENTA, COLOR_RESET)

	fmt.Printf("%sLeft:\n", simpleIndent(level+1))
	ba.Left.Print(level + 2)

	if ba.Right != nil {
		fmt.Printf("%sOperator: %sand%s\n", simpleIndent(level+1), COLOR_RED, COLOR_RESET)
		fmt.Printf("%sRight:\n", simpleIndent(level+1))
		ba.Right.Print(level + 2)
	}
}

type Equality struct {
	Left     Comparison
	Operator *Token
	Right    *Equality
}

func (eq Equality) Print(level int) {
	// Si es simple (no tiene operador), saltear este nivel
	if eq.isSimple() {
		eq.Left.Print(level)
		return
	}

	fmt.Printf("%s%sEquality%s\n", simpleIndent(level), COLOR_MAGENTA, COLOR_RESET)

	fmt.Printf("%sLeft:\n", simpleIndent(level+1))
	eq.Left.Print(level + 2)

	if eq.Operator != nil && eq.Right != nil {
		fmt.Printf("%sOperator: %s\n", simpleIndent(level+1), eq.Operator.ColorString())
		fmt.Printf("%sRight:\n", simpleIndent(level+1))
		eq.Right.Print(level + 2)
	}
}

type Comparison struct {
	Left     Term
	Operator *Token
	Right    *Comparison
}

func (c Comparison) Print(level int) {
	// Si es simple (no tiene operador), saltear este nivel
	if c.isSimple() {
		c.Left.Print(level)
		return
	}

	fmt.Printf("%s%sComparison%s\n", simpleIndent(level), COLOR_MAGENTA, COLOR_RESET)

	fmt.Printf("%sLeft:\n", simpleIndent(level+1))
	c.Left.Print(level + 2)

	if c.Operator != nil && c.Right != nil {
		fmt.Printf("%sOperator: %s\n", simpleIndent(level+1), c.Operator.ColorString())
		fmt.Printf("%sRight:\n", simpleIndent(level+1))
		c.Right.Print(level + 2)
	}
}

type Term struct {
	Left     Factor
	Operator *Token
	Right    *Term
}

func (t Term) Print(level int) {
	// Si es simple (no tiene operador), saltear este nivel
	if t.isSimple() {
		t.Left.Print(level)
		return
	}

	fmt.Printf("%s%sTerm%s\n", simpleIndent(level), COLOR_MAGENTA, COLOR_RESET)

	fmt.Printf("%sLeft:\n", simpleIndent(level+1))
	t.Left.Print(level + 2)

	if t.Operator != nil && t.Right != nil {
		fmt.Printf("%sOperator: %s\n", simpleIndent(level+1), t.Operator.ColorString())
		fmt.Printf("%sRight:\n", simpleIndent(level+1))
		t.Right.Print(level + 2)
	}
}

type Factor struct {
	Left     Unary
	Operator *Token
	Right    *Factor
}

func (f Factor) Print(level int) {
	// Si es simple (no tiene operador), saltear este nivel
	if f.isSimple() {
		switch u := f.Left.(type) {
		case UnaryWithOperator:
			u.Print(level)
		case Primary:
			u.Print(level)
		}
		return
	}

	fmt.Printf("%s%sFactor%s\n", simpleIndent(level), COLOR_MAGENTA, COLOR_RESET)

	fmt.Printf("%sLeft:\n", simpleIndent(level+1))
	switch u := f.Left.(type) {
	case UnaryWithOperator:
		u.Print(level + 2)
	case Primary:
		u.Print(level + 2)
	}

	if f.Operator != nil && f.Right != nil {
		fmt.Printf("%sOperator: %s\n", simpleIndent(level+1), f.Operator.ColorString())
		fmt.Printf("%sRight:\n", simpleIndent(level+1))
		f.Right.Print(level + 2)
	}
}

type Unary interface{}

type UnaryWithOperator struct {
	Operator Token
	Right    Unary
}

func (uwo UnaryWithOperator) Print(level int) {
	fmt.Printf("%s%sUnary%s\n", simpleIndent(level), COLOR_MAGENTA, COLOR_RESET)
	fmt.Printf("%sOperator: %s\n", simpleIndent(level+1), uwo.Operator.ColorString())
	fmt.Printf("%sRight:\n", simpleIndent(level+1))
	switch u := uwo.Right.(type) {
	case UnaryWithOperator:
		u.Print(level + 2)
	case Primary:
		u.Print(level + 2)
	}
}

type Primary struct {
	Value PrimaryValue
}

func (p Primary) Print(level int) {
	// Solo aquí mostramos un conector simple - es el más bajo nivel
	fmt.Printf("%s└── ", simpleIndent(level))
	switch v := p.Value.(type) {
	case Token:
		fmt.Printf("%s%s%s\n", COLOR_GREEN, v.ColorString(), COLOR_RESET)
	case Call:
		fmt.Printf("%sCall%s\n", COLOR_GREEN, COLOR_RESET)
		v.Print(level + 1)
	case GroupingExpression:
		fmt.Printf("%sGrouping%s\n", COLOR_GREEN, COLOR_RESET)
		v.Print(level + 1)
	default:
		fmt.Printf("%s%v%s\n", COLOR_GREEN, v, COLOR_RESET)
	}
}

type PrimaryValue interface{}

type GroupingExpression struct {
	Expression Expression
}

func (ge GroupingExpression) Print(level int) {
	fmt.Printf("%s%sGroupingExpression%s\n", simpleIndent(level), COLOR_GREEN, COLOR_RESET)
	fmt.Printf("%sExpression:\n", simpleIndent(level+1))
	ge.Expression.Print(level + 2)
}

type Call struct {
	Callee    string
	Arguments []Expression
}

func (c Call) Print(level int) {
	fmt.Printf("%s%sCall%s\n", simpleIndent(level), COLOR_GREEN, COLOR_RESET)
	fmt.Printf("%sCallee: %s%s%s\n", simpleIndent(level+1), COLOR_WHITE, c.Callee, COLOR_RESET)

	if len(c.Arguments) > 0 {
		fmt.Printf("%sArguments:\n", simpleIndent(level+1))
		for i, arg := range c.Arguments {
			fmt.Printf("%sArg[%d]:\n", simpleIndent(level+2), i)
			arg.Print(level + 3)
		}
	}
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
