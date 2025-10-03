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

type MergableNode interface {
	IsComplete() bool
	IsLeftComplete() bool
	IsRightComplete() bool
	IsOperatorComplete() bool

	GetLeft() MergableNode
	GetRight() MergableNode
	GetOperator() *Token

	SetLeft(n MergableNode)
	SetRight(n MergableNode)
	SetOperator(t *Token)
}

type Expression struct {
	Root *BinaryOr
}

func (bo *BinaryOr) skipPrinting() bool {
	return bo.Right == nil
}

func (ba *BinaryAnd) skipPrinting() bool {
	return ba.Right == nil
}

func (eq *Equality) skipPrinting() bool {
	return eq.Operator == nil && eq.Right == nil
}

func (c *Comparison) skipPrinting() bool {
	return c.Operator == nil && c.Right == nil
}

func (t *Term) skipPrinting() bool {
	return t.Operator == nil && t.Right == nil
}

func (f *Factor) skipPrinting() bool {
	return f.Operator == nil && f.Right == nil
}

func (e *Expression) Print(start string) {
	e.Root.Print(start)
}

type BinaryOr struct {
	Left     *BinaryAnd
	Operator *Token
	Right    *BinaryOr
}

func (bo BinaryOr) IsComplete() bool {
	return bo.IsLeftComplete() && bo.IsOperatorComplete() && bo.IsRightComplete()
}

func (bo BinaryOr) IsEmpty() bool {
	return !bo.IsLeftComplete() && !bo.IsOperatorComplete() && !bo.IsRightComplete()
}

func (bo BinaryOr) IsLeftComplete() bool {
	return bo.Left != nil
}

func (bo BinaryOr) IsRightComplete() bool {
	return bo.Right != nil
}

func (bo BinaryOr) IsOperatorComplete() bool {
	return bo.Operator != nil
}

func (bo *BinaryOr) GetLeft() MergableNode {
	if bo == nil {
		return nil
	}
	return bo.Left
}

func (bo *BinaryOr) GetRight() MergableNode {
	if bo == nil {
		return nil
	}
	return bo.Right
}

func (bo *BinaryOr) GetOperator() *Token {
	if bo == nil {
		return nil
	}
	return bo.Operator
}

func (bo *BinaryOr) SetLeft(val MergableNode) {
	if bo == nil {
		return
	}
	bo.Left = val.(*BinaryAnd)
}

func (bo *BinaryOr) SetRight(val MergableNode) {
	if bo == nil {
		return
	}
	bo.Right = val.(*BinaryOr)
}

func (bo *BinaryOr) SetOperator(op *Token) {
	if bo == nil {
		return
	}
	bo.Operator = op
}

func (bo *BinaryOr) Print(start string) {
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
	Left     *Equality
	Operator *Token
	Right    *BinaryAnd
}

func (ba *BinaryAnd) IsComplete() bool {
	return ba.IsLeftComplete() && ba.IsOperatorComplete() && ba.IsRightComplete()
}

func (ba *BinaryAnd) IsLeftComplete() bool {
	return ba.Left != nil
}

func (ba *BinaryAnd) IsRightComplete() bool {
	return ba.Right != nil
}

func (ba *BinaryAnd) IsOperatorComplete() bool {
	return ba.Operator != nil
}

func (ba *BinaryAnd) GetLeft() MergableNode {
	if ba == nil || ba.Left == nil {
		return nil
	}
	return ba.Left
}

func (ba *BinaryAnd) GetRight() MergableNode {
	if ba == nil || ba.Right == nil {
		return nil
	}
	return ba.Right
}

func (ba *BinaryAnd) GetOperator() *Token {
	if ba == nil {
		return nil
	}
	return ba.Operator
}

func (ba *BinaryAnd) SetLeft(n MergableNode) {
	if ba == nil {
		return
	}
	ba.Left = n.(*Equality)
}

func (ba *BinaryAnd) SetRight(n MergableNode) {
	if ba == nil {
		return
	}
	ba.Right = n.(*BinaryAnd)
}

func (ba *BinaryAnd) SetOperator(t *Token) {
	if ba == nil {
		return
	}
	ba.Operator = t
}

func (ba *BinaryAnd) Print(start string) {
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
	Left     *Comparison
	Operator *Token
	Right    *Equality
}

func (eq *Equality) IsComplete() bool {
	return eq.IsLeftComplete() && eq.IsOperatorComplete() && eq.IsRightComplete()
}

func (eq *Equality) IsLeftComplete() bool {
	return eq.Left != nil
}

func (eq *Equality) IsRightComplete() bool {
	return eq.Right != nil
}

func (eq *Equality) IsOperatorComplete() bool {
	return eq.Operator != nil
}

func (eq *Equality) GetLeft() MergableNode {
	if eq == nil || eq.Left == nil {
		return nil
	}
	return eq.Left
}

func (eq *Equality) GetRight() MergableNode {
	if eq == nil || eq.Right == nil {
		return nil
	}
	return eq.Right
}

func (eq *Equality) GetOperator() *Token {
	if eq == nil {
		return nil
	}
	return eq.Operator
}

func (eq *Equality) SetLeft(n MergableNode) {
	if eq == nil {
		return
	}
	eq.Left = n.(*Comparison)
}

func (eq *Equality) SetRight(n MergableNode) {
	if eq == nil {
		return
	}
	eq.Right = n.(*Equality)
}

func (eq *Equality) SetOperator(t *Token) {
	if eq == nil {
		return
	}
	eq.Operator = t
}

func (eq *Equality) Print(start string) {
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
	Left     *Term
	Operator *Token
	Right    *Comparison
}

func (c *Comparison) IsComplete() bool {
	return c.IsLeftComplete() && c.IsOperatorComplete() && c.IsRightComplete()
}

func (c *Comparison) IsLeftComplete() bool {
	return c.Left != nil
}

func (c *Comparison) IsRightComplete() bool {
	return c.Right != nil
}

func (c *Comparison) IsOperatorComplete() bool {
	return c.Operator != nil
}

func (c *Comparison) GetLeft() MergableNode {
	if c == nil || c.Left == nil {
		return nil
	}
	return c.Left
}

func (c *Comparison) GetRight() MergableNode {
	if c == nil || c.Right == nil {
		return nil
	}
	return c.Right
}

func (c *Comparison) GetOperator() *Token {
	if c == nil {
		return nil
	}
	return c.Operator
}

func (c *Comparison) SetLeft(n MergableNode) {
	if c == nil {
		return
	}
	c.Left = n.(*Term)
}

func (c *Comparison) SetRight(n MergableNode) {
	if c == nil {
		return
	}
	c.Right = n.(*Comparison)
}

func (c *Comparison) SetOperator(t *Token) {
	if c == nil {
		return
	}
	c.Operator = t
}

func (c *Comparison) Print(start string) {
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
	Left     *Factor
	Operator *Token
	Right    *Term
}

func (t *Term) IsComplete() bool {
	return t.IsLeftComplete() && t.IsOperatorComplete() && t.IsRightComplete()
}

func (t *Term) IsLeftComplete() bool {
	return t.Left != nil
}

func (t *Term) IsRightComplete() bool {
	return t.Right != nil
}

func (t *Term) IsOperatorComplete() bool {
	return t.Operator != nil
}

func (t *Term) GetLeft() MergableNode {
	if t == nil || t.Left == nil {
		return nil
	}
	return t.Left
}

func (t *Term) GetRight() MergableNode {
	if t == nil || t.Right == nil {
		return nil
	}
	return t.Right
}

func (t *Term) GetOperator() *Token {
	if t == nil {
		return nil
	}
	return t.Operator
}

func (t *Term) SetLeft(n MergableNode) {
	if t == nil {
		return
	}
	t.Left = n.(*Factor)
}

func (t *Term) SetRight(n MergableNode) {
	if t == nil {
		return
	}
	t.Right = n.(*Term)
}

func (t *Term) SetOperator(op *Token) {
	if t == nil {
		return
	}
	t.Operator = op
}

func (t *Term) Print(start string) {
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

func (f *Factor) IsComplete() bool {
	return f.IsLeftComplete() && f.IsOperatorComplete() && f.IsRightComplete()
}

func (f *Factor) IsLeftComplete() bool {
	return f.Left != nil
}

func (f *Factor) IsRightComplete() bool {
	return f.Right != nil
}

func (f *Factor) IsOperatorComplete() bool {
	return f.Operator != nil
}

func (f *Factor) GetLeft() MergableNode {
	if f == nil || f.Left == nil {
		return nil
	}
	if node, ok := f.Left.(MergableNode); ok {
		return node
	}
	return nil
}

func (f *Factor) GetRight() MergableNode {
	if f == nil || f.Right == nil {
		return nil
	}
	return f.Right
}

func (f *Factor) GetOperator() *Token {
	if f == nil {
		return nil
	}
	return f.Operator
}

func (f *Factor) SetLeft(n MergableNode) {
	if f == nil {
		return
	}
	if unary, ok := n.(Unary); ok {
		f.Left = unary
	}
}

func (f *Factor) SetRight(n MergableNode) {
	if f == nil {
		return
	}
	f.Right = n.(*Factor)
}

func (f *Factor) SetOperator(op *Token) {
	if f == nil {
		return
	}
	f.Operator = op
}

func (f *Factor) Print(start string) {
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
	MergableNode
	Print(start string)
}

type UnaryWithOperator struct {
	Operator *Token
	Right    Unary
}

func (uwo UnaryWithOperator) Print(start string) {
	nodeName := fmt.Sprintf("Unary (%s)", friendlyOperatorName(uwo.Operator, true))
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	uwo.Right.Print(start + string(LAST_CONNECTOR))
}

func (uwo *UnaryWithOperator) IsComplete() bool {
	return uwo.Operator != nil && uwo.Right != nil
}
func (uwo *UnaryWithOperator) IsEmpty() bool {
	return uwo.Operator == nil && uwo.Right == nil
}
func (uwo *UnaryWithOperator) IsRightComplete() bool {
	return uwo.Right != nil
}
func (uwo *UnaryWithOperator) IsLeftComplete() bool {
	return false
}

func (uwo *UnaryWithOperator) IsOperatorComplete() bool {
	return uwo.Operator != nil
}

func (uwo *UnaryWithOperator) GetLeft() MergableNode {
	return nil
}

func (uwo *UnaryWithOperator) GetRight() MergableNode {
	if uwo == nil || uwo.Right == nil {
		return nil
	}
	return uwo.Right
}

func (uwo *UnaryWithOperator) GetOperator() *Token {
	return uwo.Operator
}

func (uwo *UnaryWithOperator) SetLeft(n MergableNode) {
	// UnaryWithOperator does not have a left node
}

func (uwo *UnaryWithOperator) SetRight(n MergableNode) {
	if unary, ok := n.(Unary); ok {
		uwo.Right = unary
	}
}

func (uwo *UnaryWithOperator) SetOperator(t *Token) {
	uwo.Operator = t
}

type Primary struct {
	Value PrimaryValue
}

func (p *Primary) Print(start string) {
	p.Value.Print(start)
}

func (ip *Primary) IsEmpty() bool {
	return ip.Value == nil
}

func (ip *Primary) IsComplete() bool {
	return ip.Value != nil
}

func (ip *Primary) IsRightComplete() bool {
	return ip.Value != nil
}

func (ip *Primary) IsLeftComplete() bool {
	return false
}

func (ip *Primary) IsOperatorComplete() bool {
	return false
}

func (ip *Primary) GetLeft() MergableNode {
	return nil
}

func (ip *Primary) GetRight() MergableNode {
	return nil
}

func (ip *Primary) GetOperator() *Token {
	return nil
}

func (ip *Primary) SetLeft(n MergableNode) {
	// Primary does not have a left node
}

func (ip *Primary) SetRight(n MergableNode) {
	// Primary does not have a right node
}

func (ip *Primary) SetOperator(t *Token) {
	// Primary does not have an operator
}

type PrimaryValue interface {
	Print(start string)
}

type GroupingExpression struct {
	Expression *Expression
}

func (ge GroupingExpression) Print(start string) {
	fmt.Printf("%s%s\n", start, Colorize("GroupingExpression", COLOR_GREEN))
	start = advanceSuffix(start)
	ge.Expression.Print(start + string(LAST_CONNECTOR))
}

type Call struct {
	Callee    *string
	Arguments []Expression
}

func (c Call) Print(start string) {
	nodeName := "Call"
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_GREEN))
	start = advanceSuffix(start)
	fmt.Printf("%sCallee: %s\n", start+string(BRANCH_CONNECTOR), Colorize(*c.Callee, COLOR_WHITE))

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
