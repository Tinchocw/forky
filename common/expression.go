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

type Direction int

const (
	Left Direction = iota + 1
	Right
)

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
	HasLeft() bool
	HasRight() bool
	HasOperator() bool

	GetLeft() MergableNode
	GetRight() MergableNode
	GetOperator() *Token

	SetLeft(n MergableNode)
	SetRight(n MergableNode)
	SetOperator(t *Token)

	ExpresionDepth(direction Direction) int
	GetSubExpression(level int, direction Direction) MergableNode
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

// func (bo BinaryOr) IsComplete() bool {
// 	return bo.HasLeft() && bo.HasOperator() && bo.HasRight()
// }

func (bo BinaryOr) IsEmpty() bool {
	return !bo.HasLeft() && !bo.HasOperator() && !bo.HasRight()
}

func (bo BinaryOr) HasLeft() bool {
	return bo.Left != nil
}

func (bo BinaryOr) HasRight() bool {
	return bo.Right != nil
}

func (bo BinaryOr) HasOperator() bool {
	return bo.Operator != nil
}

func (bo *BinaryOr) GetLeft() MergableNode {
	return bo.Left
}

func (bo *BinaryOr) GetRight() MergableNode {
	return bo.Right
}

func (bo *BinaryOr) GetOperator() *Token {
	if bo.Operator == nil {
		panic("GetOperator: expected *Token")
	}

	return bo.Operator
}

func (bo *BinaryOr) SetLeft(val MergableNode) {
	binaryAnd, ok := val.(*BinaryAnd)
	if !ok {
		panic("SetLeft: expected *BinaryAnd")
	}

	bo.Left = binaryAnd
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

func (bo *BinaryOr) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if bo.HasRight() {
			return bo.Right.ExpresionDepth(direction)
		}

		if bo.HasOperator() {
			return 0
		}
	}

	if !bo.HasLeft() {
		return 0
	}

	return bo.Left.ExpresionDepth(direction)
}

func (bo *BinaryOr) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if bo.HasRight() {
			return bo.Right.GetSubExpression(level, direction)
		}

		if bo.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !bo.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return bo.Left.GetSubExpression(level, direction)
}

type BinaryAnd struct {
	Left     *Equality
	Operator *Token
	Right    *BinaryAnd
}

// func (ba *BinaryAnd) IsComplete() bool {
// 	return ba.HasLeft() && ba.HasOperator() && ba.HasRight()
// }

func (ba *BinaryAnd) HasLeft() bool {
	return ba.Left != nil
}

func (ba *BinaryAnd) HasRight() bool {
	return ba.Right != nil
}

func (ba *BinaryAnd) HasOperator() bool {
	return ba.Operator != nil
}

func (ba *BinaryAnd) GetLeft() MergableNode {
	if ba.Left == nil {
		panic("GetLeft: expected *Equality")
	}
	return ba.Left
}

func (ba *BinaryAnd) GetRight() MergableNode {
	if ba.Right == nil {
		panic("GetRight: expected *BinaryAnd")
	}
	return ba.Right
}

func (ba *BinaryAnd) GetOperator() *Token {
	if ba.Operator == nil {
		panic("GetOperator: expected *Token")
	}
	return ba.Operator
}

func (ba *BinaryAnd) SetLeft(val MergableNode) {
	equality, ok := val.(*Equality)
	if !ok {
		panic("SetLeft: expected *Equality")
	}
	ba.Left = equality
}

func (ba *BinaryAnd) SetRight(val MergableNode) {
	binaryAnd, ok := val.(*BinaryAnd)
	if !ok {
		panic("SetRight: expected *BinaryAnd")
	}
	ba.Right = binaryAnd
}

func (ba *BinaryAnd) SetOperator(t *Token) {
	if t == nil {
		panic("SetOperator: expected *Token")
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

func (ba *BinaryAnd) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if ba.HasRight() {
			return ba.Right.ExpresionDepth(direction)
		}

		if ba.HasOperator() {
			return 0
		}
	}

	if !ba.HasLeft() {
		return 0
	}

	return ba.Left.ExpresionDepth(direction)
}

func (ba *BinaryAnd) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if ba.HasRight() {
			return ba.Right.GetSubExpression(level, direction)
		}

		if ba.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !ba.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return ba.Left.GetSubExpression(level, direction)
}

type Equality struct {
	Left     *Comparison
	Operator *Token
	Right    *Equality
}

// func (eq *Equality) IsComplete() bool {
// 	return eq.HasLeft() && eq.HasOperator() && eq.HasRight()
// }

func (eq *Equality) HasLeft() bool {
	return eq.Left != nil
}

func (eq *Equality) HasRight() bool {
	return eq.Right != nil
}

func (eq *Equality) HasOperator() bool {
	return eq.Operator != nil
}

func (eq *Equality) GetLeft() MergableNode {
	if eq.Left == nil {
		panic("GetLeft: expected *Comparison")
	}
	return eq.Left
}

func (eq *Equality) GetRight() MergableNode {
	if eq.Right == nil {
		panic("GetRight: expected *Equality")
	}
	return eq.Right
}

func (eq *Equality) GetOperator() *Token {
	if eq.Operator == nil {
		panic("GetOperator: expected *Token")
	}
	return eq.Operator
}

func (eq *Equality) SetLeft(val MergableNode) {
	comparison, ok := val.(*Comparison)
	if !ok {
		panic("SetLeft: expected *Comparison")
	}
	eq.Left = comparison
}

func (eq *Equality) SetRight(val MergableNode) {
	equality, ok := val.(*Equality)
	if !ok {
		panic("SetRight: expected *Equality")
	}
	eq.Right = equality
}

func (eq *Equality) SetOperator(t *Token) {
	if t == nil {
		panic("SetOperator: expected *Token")
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

func (eq *Equality) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if eq.HasRight() {
			return eq.Right.ExpresionDepth(direction)
		}

		if eq.HasOperator() {
			return 0
		}
	}

	if !eq.HasLeft() {
		return 0
	}

	return eq.Left.ExpresionDepth(direction)
}

func (eq *Equality) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if eq.HasRight() {
			return eq.Right.GetSubExpression(level, direction)
		}

		if eq.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !eq.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return eq.Left.GetSubExpression(level, direction)
}

type Comparison struct {
	Left     *Term
	Operator *Token
	Right    *Comparison
}

// func (c *Comparison) IsComplete() bool {
// 	return c.HasLeft() && c.HasOperator() && c.HasRight()
// }

func (c *Comparison) HasLeft() bool {
	return c.Left != nil
}

func (c *Comparison) HasRight() bool {
	return c.Right != nil
}

func (c *Comparison) HasOperator() bool {
	return c.Operator != nil
}

func (c *Comparison) GetLeft() MergableNode {
	if c.Left == nil {
		panic("GetLeft: expected *Term")
	}
	return c.Left
}

func (c *Comparison) GetRight() MergableNode {
	if c.Right == nil {
		panic("GetRight: expected *Comparison")
	}
	return c.Right
}

func (c *Comparison) GetOperator() *Token {
	if c.Operator == nil {
		panic("GetOperator: expected *Token")
	}
	return c.Operator
}

func (c *Comparison) SetLeft(val MergableNode) {
	term, ok := val.(*Term)
	if !ok {
		panic("SetLeft: expected *Term")
	}
	c.Left = term
}

func (c *Comparison) SetRight(val MergableNode) {
	comparison, ok := val.(*Comparison)
	if !ok {
		panic("SetRight: expected *Comparison")
	}
	c.Right = comparison
}

func (c *Comparison) SetOperator(t *Token) {
	if t == nil {
		panic("SetOperator: expected *Token")
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

func (c *Comparison) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if c.HasRight() {
			return c.Right.ExpresionDepth(direction)
		}

		if c.HasOperator() {
			return 0
		}
	}

	if !c.HasLeft() {
		return 0
	}

	return c.Left.ExpresionDepth(direction)
}

func (c *Comparison) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if c.HasRight() {
			return c.Right.GetSubExpression(level, direction)
		}

		if c.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !c.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return c.Left.GetSubExpression(level, direction)
}

type Term struct {
	Left     *Factor
	Operator *Token
	Right    *Term
}

// func (t *Term) IsComplete() bool {
// 	return t.HasLeft() && t.HasOperator() && t.HasRight()
// }

func (t *Term) HasLeft() bool {
	return t.Left != nil
}

func (t *Term) HasRight() bool {
	return t.Right != nil
}

func (t *Term) HasOperator() bool {
	return t.Operator != nil
}

func (t *Term) GetLeft() MergableNode {
	if t.Left == nil {
		panic("GetLeft: expected *Factor")
	}
	return t.Left
}

func (t *Term) GetRight() MergableNode {
	if t.Right == nil {
		panic("GetRight: expected *Term")
	}
	return t.Right
}

func (t *Term) GetOperator() *Token {
	if t.Operator == nil {
		panic("GetOperator: expected *Token")
	}
	return t.Operator
}

func (t *Term) SetLeft(val MergableNode) {
	factor, ok := val.(*Factor)
	if !ok {
		panic("SetLeft: expected *Factor")
	}
	t.Left = factor
}

func (t *Term) SetRight(val MergableNode) {
	term, ok := val.(*Term)
	if !ok {
		panic("SetRight: expected *Term")
	}
	t.Right = term
}

func (t *Term) SetOperator(op *Token) {
	if op == nil {
		panic("SetOperator: expected *Token")
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

func (t *Term) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if t.HasRight() {
			return t.Right.ExpresionDepth(direction)
		}

		if t.HasOperator() {
			return 0
		}
	}

	if !t.HasLeft() {
		return 0
	}

	return t.Left.ExpresionDepth(direction)
}

func (t *Term) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if t.HasRight() {
			return t.Right.GetSubExpression(level, direction)
		}

		if t.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !t.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return t.Left.GetSubExpression(level, direction)
}

type Factor struct {
	Left     Unary
	Operator *Token
	Right    *Factor
}

// func (f *Factor) IsComplete() bool {
// 	return f.HasLeft() && f.HasOperator() && f.HasRight()
// }

func (f *Factor) HasLeft() bool {
	return f.Left != nil
}

func (f *Factor) HasRight() bool {
	return f.Right != nil
}

func (f *Factor) HasOperator() bool {
	return f.Operator != nil
}

func (f *Factor) GetLeft() MergableNode {
	if f.Left == nil {
		panic("GetLeft: expected Unary")
	}
	node, ok := f.Left.(MergableNode)
	if !ok {
		panic("GetLeft: expected MergableNode (Unary)")
	}
	return node
}

func (f *Factor) GetRight() MergableNode {
	if f.Right == nil {
		panic("GetRight: expected *Factor")
	}
	return f.Right
}

func (f *Factor) GetOperator() *Token {
	if f.Operator == nil {
		panic("GetOperator: expected *Token")
	}
	return f.Operator
}

func (f *Factor) SetLeft(val MergableNode) {
	unary, ok := val.(Unary)
	if !ok {
		panic("SetLeft: expected Unary")
	}
	f.Left = unary
}

func (f *Factor) SetRight(val MergableNode) {
	factor, ok := val.(*Factor)
	if !ok {
		panic("SetRight: expected *Factor")
	}
	f.Right = factor
}

func (f *Factor) SetOperator(op *Token) {
	if op == nil {
		panic("SetOperator: expected *Token")
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

func (f *Factor) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if f.HasRight() {
			return f.Right.ExpresionDepth(direction)
		}

		if f.HasOperator() {
			return 0
		}
	}

	if !f.HasLeft() {
		return 0
	}

	return f.Left.ExpresionDepth(direction)
}

func (f *Factor) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if f.HasRight() {
			return f.Right.GetSubExpression(level, direction)
		}

		if f.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !f.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return f.Left.GetSubExpression(level, direction)
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

//	func (uwo *UnaryWithOperator) IsComplete() bool {
//		return uwo.Operator != nil && uwo.Right != nil
//	}

func (uwo *UnaryWithOperator) IsEmpty() bool {
	return uwo.Operator == nil && uwo.Right == nil
}
func (uwo *UnaryWithOperator) HasRight() bool {
	return uwo.Right != nil
}
func (uwo *UnaryWithOperator) HasLeft() bool {
	return false
}

func (uwo *UnaryWithOperator) HasOperator() bool {
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
	panic("SetLeft: UnaryWithOperator does not have a left node")
}

func (uwo *UnaryWithOperator) SetRight(val MergableNode) {
	unary, ok := val.(Unary)
	if !ok {
		panic("SetRight: expected Unary")
	}
	uwo.Right = unary
}

func (uwo *UnaryWithOperator) SetOperator(t *Token) {
	if t == nil {
		panic("SetOperator: expected *Token")
	}
	uwo.Operator = t
}

func (uwo *UnaryWithOperator) ExpresionDepth(direction Direction) int {
	if direction == Right {
		if uwo.HasRight() {
			return uwo.Right.ExpresionDepth(direction)
		}

		if uwo.HasOperator() {
			return 0
		}
	}

	if !uwo.HasLeft() {
		return 0
	}

	return uwo.GetLeft().ExpresionDepth(direction)
}

func (uwo *UnaryWithOperator) GetSubExpression(level int, direction Direction) MergableNode {
	if direction == Right {
		if uwo.HasRight() {
			return uwo.Right.GetSubExpression(level, direction)
		}

		if uwo.HasOperator() {
			panic("GetSubExpression: no sub-expression found at the specified level")
		}
	}

	if !uwo.HasLeft() {
		panic("GetSubExpression: no sub-expression found at the specified level")
	}

	return uwo.GetLeft().GetSubExpression(level, direction)
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

// func (ip *Primary) IsComplete() bool {
// 	return ip.Value != nil
// }

func (ip *Primary) HasRight() bool {
	return ip.Value != nil
}

func (ip *Primary) HasLeft() bool {
	return false
}

func (ip *Primary) HasOperator() bool {
	return false
}

func (ip *Primary) GetLeft() MergableNode {
	// to be implemented
	panic("GetLeft: Primary does not have a left node")
}

func (ip *Primary) GetRight() MergableNode {
	// to be implemented
	panic("GetRight: Primary does not have a right node")
}

func (ip *Primary) GetOperator() *Token {
	// to be implemented
	panic("GetOperator: Primary does not have an operator")
}

func (ip *Primary) SetLeft(n MergableNode) {
	// to be implemented
	panic("SetLeft: Primary does not have a left node")
}

func (ip *Primary) SetRight(n MergableNode) {
	// to be implemented
	panic("SetRight: Primary does not have a right node")
}

func (ip *Primary) SetOperator(t *Token) {
	// to be implemented
	panic("SetOperator: Primary does not have an operator")
}

func (ip *Primary) ExpresionDepth(direction Direction) int {
	// to be implemented
	return 0
}

func (ip *Primary) GetSubExpression(level int, direction Direction) MergableNode {
	// to be implemented
	return nil
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
