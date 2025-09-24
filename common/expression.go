package common

import (
	"fmt"
)

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

func friendlyOperatorName(token *Token, isUnary bool) string {
	if token == nil {
		return ""
	}
	switch token.Typ {
	case EQUAL_EQUAL:
		return "EQUALS"
	case BANG_EQUAL:
		return "NOT_EQUALS"
	case LESS:
		return "LESS_THAN"
	case GREATER:
		return "GREATER_THAN"
	case LESS_EQUAL:
		return "LESS_OR_EQUAL"
	case GREATER_EQUAL:
		return "GREATER_OR_EQUAL"
	case PLUS:
		if isUnary {
			return "POSITIVE"
		}
		return "PLUS"
	case MINUS:
		if isUnary {
			return "NEGATIVE"
		}
		return "MINUS"
	case ASTERISK:
		return "MULTIPLY"
	case SLASH:
		return "DIVIDE"
	case BANG:
		return "NOT"
	default:
		return token.String()
	}
}

type Expression struct {
	Root BinaryOr
}

func (bo BinaryOr) skipPrinting() bool {
	return bo.Right == nil
}

func (ba BinaryAnd) skipPrinting() bool {
	return ba.Right == nil
}

func (eq Equality) skipPrinting() bool {
	return eq.Operator == nil && eq.Right == nil
}

func (c Comparison) skipPrinting() bool {
	return c.Operator == nil && c.Right == nil
}

func (t Term) skipPrinting() bool {
	return t.Operator == nil && t.Right == nil
}

func (f Factor) skipPrinting() bool {
	return f.Operator == nil && f.Right == nil
}

func (e Expression) Print(start string) {
	e.Root.Print(start)
}

type BinaryOr struct {
	Left  BinaryAnd
	Right *BinaryOr
}

func (bo BinaryOr) Print(start string) {
	if bo.skipPrinting() {
		bo.Left.Print(start)
		return
	}

	nodeName := "BinaryOr"
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	bo.Left.Print(start + string(BRANCH_CONNECTOR))
	bo.Right.Print(start + string(LAST_CONNECTOR))
}

type BinaryAnd struct {
	Left  Equality
	Right *BinaryAnd
}

func (ba BinaryAnd) Print(start string) {
	if ba.skipPrinting() {
		ba.Left.Print(start)
		return
	}

	nodeName := "BinaryAnd"
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	ba.Left.Print(start + string(BRANCH_CONNECTOR))
	ba.Right.Print(start + string(LAST_CONNECTOR))
}

type Equality struct {
	Left     Comparison
	Operator *Token
	Right    *Equality
}

func (eq Equality) Print(start string) {
	if eq.skipPrinting() {
		eq.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Equality (%s)", friendlyOperatorName(eq.Operator, false))
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	eq.Left.Print(start + string(BRANCH_CONNECTOR))
	eq.Right.Print(start + string(LAST_CONNECTOR))
}

type Comparison struct {
	Left     Term
	Operator *Token
	Right    *Comparison
}

func (c Comparison) Print(start string) {
	if c.skipPrinting() {
		c.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Comparison (%s)", friendlyOperatorName(c.Operator, false))
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	c.Left.Print(start + string(BRANCH_CONNECTOR))
	c.Right.Print(start + string(LAST_CONNECTOR))
}

type Term struct {
	Left     Factor
	Operator *Token
	Right    *Term
}

func (t Term) Print(start string) {
	if t.skipPrinting() {
		t.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Term (%s)", friendlyOperatorName(t.Operator, false))

	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	t.Left.Print(start + string(BRANCH_CONNECTOR))
	t.Right.Print(start + string(LAST_CONNECTOR))
}

type Factor struct {
	Left     Unary
	Operator *Token
	Right    *Factor
}

func (f Factor) Print(start string) {
	if f.skipPrinting() {
		f.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Factor (%s)", friendlyOperatorName(f.Operator, false))

	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	f.Left.Print(start + string(BRANCH_CONNECTOR))
	f.Right.Print(start + string(LAST_CONNECTOR))
}

type Unary interface {
	Print(start string)
}

type UnaryWithOperator struct {
	Operator Token
	Right    Unary
}

func (uwo UnaryWithOperator) Print(start string) {
	nodeName := fmt.Sprintf("Unary (%s)", friendlyOperatorName(&uwo.Operator, true))
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	uwo.Right.Print(start + string(LAST_CONNECTOR))
}

type Primary struct {
	Value PrimaryValue
}

func (p Primary) Print(start string) {
	p.Value.Print(start)
}

type PrimaryValue interface {
	Print(start string)
}

type GroupingExpression struct {
	Expression Expression
}

func (ge GroupingExpression) Print(start string) {
	fmt.Printf("%s%s\n", start, Colorize("GroupingExpression", COLOR_GREEN))
	start = advanceSuffix(start)
	ge.Expression.Print(start + string(LAST_CONNECTOR))
}

type Call struct {
	Callee    string
	Arguments []Expression
}

func (c Call) Print(start string) {
	nodeName := "Call"
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_GREEN))
	start = advanceSuffix(start)
	fmt.Printf("%sCallee: %s\n", start+string(BRANCH_CONNECTOR), Colorize(c.Callee, COLOR_WHITE))

	if len(c.Arguments) > 0 {
		nodeName := "Arguments"
		fmt.Printf("%s%s\n", start+string(LAST_CONNECTOR), Colorize(nodeName, COLOR_GREEN))
		start += string(SIMPLE_INDENT)
		for i, arg := range c.Arguments {
			connector := BRANCH_CONNECTOR
			identation := SIMPLE_CONNECTOR
			if i == len(c.Arguments)-1 {
				connector = LAST_CONNECTOR
				identation = SIMPLE_INDENT
			}
			fmt.Printf("%sArg[%d]:\n", start+string(connector), i)
			arg.Print(start + string(identation) + string(LAST_CONNECTOR))
		}
	}
}
