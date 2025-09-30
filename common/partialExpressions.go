package common

import (
	"fmt"
)

type IncompleteExpression struct {
	Root *IncompleteBinaryOr
}

func (bo IncompleteBinaryOr) skipPrinting() bool {
	return bo.Right == nil
}

func (ba IncompleteBinaryAnd) skipPrinting() bool {
	return ba.Right == nil
}

func (eq IncompleteEquality) skipPrinting() bool {
	return eq.Operator == nil && eq.Right == nil
}

func (c IncompleteComparison) skipPrinting() bool {
	return c.Operator == nil && c.Right == nil
}

func (t IncompleteTerm) skipPrinting() bool {
	return t.Operator == nil && t.Right == nil
}

func (f IncompleteFactor) skipPrinting() bool {
	return f.Operator == nil && f.Right == nil
}

func (e IncompleteExpression) Print(start string) {
	e.Root.Print(start)
}

type IncompleteBinaryOr struct {
	Left  *IncompleteBinaryAnd
	Right *IncompleteBinaryOr
}

func (bo IncompleteBinaryOr) IsComplete() bool {
	if bo.Left != nil && bo.Right != nil {
		return true
	}
	return false
}

func (bo IncompleteBinaryOr) IsEmpty() bool {
	if bo.Left == nil && bo.Right == nil {
		return true
	}
	return false
}

func (bo IncompleteBinaryOr) IsLeftComplete() bool {
	if bo.Left != nil /* && bo.Left.isComplete() */ {
		return true
	}
	return false
}

func (bo IncompleteBinaryOr) IsRightComplete() bool {
	if bo.Right != nil && bo.Right.IsComplete() {
		return true
	}
	return false
}

func (bo IncompleteBinaryOr) Print(start string) {
	if bo.skipPrinting() {
		bo.Left.Print(start)
		return
	}

	nodeName := "IncompleteBinaryOr"
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	bo.Left.Print(start + string(BRANCH_CONNECTOR))
	bo.Right.Print(start + string(LAST_CONNECTOR))
}

type IncompleteBinaryAnd struct {
	Left  *IncompleteEquality
	Right *IncompleteBinaryAnd
}

func (ba IncompleteBinaryAnd) Print(start string) {
	if ba.skipPrinting() {
		ba.Left.Print(start)
		return
	}

	nodeName := "IncompleteBinaryAnd"

	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	ba.Left.Print(start + string(BRANCH_CONNECTOR))
	ba.Right.Print(start + string(LAST_CONNECTOR))
}

func (ba IncompleteBinaryAnd) IsComplete() bool {
	if ba.Left != nil && ba.Right != nil {
		return true
	}
	return false
}

func (ba IncompleteBinaryAnd) IsEmpty() bool {
	if ba.Left == nil && ba.Right == nil {
		return true
	}
	return false
}

func (ba IncompleteBinaryAnd) IsLeftComplete() bool {
	if ba.Left != nil /* && ba.Left.isComplete() */ {
		return true
	}
	return false
}

func (ba IncompleteBinaryAnd) IsRightComplete() bool {
	if ba.Right != nil && ba.Right.IsComplete() {
		return true
	}
	return false
}

type IncompleteEquality struct {
	Left     *IncompleteComparison
	Operator *Token
	Right    *IncompleteEquality
}

func (eq IncompleteEquality) Print(start string) {
	if eq.skipPrinting() {
		eq.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("IncompleteEquality (%s)", friendlyOperatorName(eq.Operator, false))

	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	eq.Left.Print(start + string(BRANCH_CONNECTOR))
	eq.Right.Print(start + string(LAST_CONNECTOR))
}

func (eq IncompleteEquality) IsComplete() bool {
	if eq.Left != nil && eq.Operator != nil && eq.Right != nil {
		return true
	}
	return false
}

func (eq IncompleteEquality) IsEmpty() bool {
	if eq.Left == nil && eq.Operator == nil && eq.Right == nil {
		return true
	}
	return false
}

func (eq IncompleteEquality) IsLeftComplete() bool {
	if eq.Left != nil /* && eq.Left.isComplete() */ {
		return true
	}
	return false
}

func (eq IncompleteEquality) IsRightComplete() bool {
	if eq.Right != nil && eq.Right.IsComplete() {
		return true
	}
	return false
}

func (eq IncompleteEquality) IsOperatorComplete() bool {
	if eq.Operator != nil {
		return true
	}
	return false
}

type IncompleteComparison struct {
	Left     *IncompleteTerm
	Operator *Token
	Right    *IncompleteComparison
}

func (c IncompleteComparison) Print(start string) {
	if c.skipPrinting() {
		c.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("IncompleteComparison (%s)", friendlyOperatorName(c.Operator, false))

	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	c.Left.Print(start + string(BRANCH_CONNECTOR))
	c.Right.Print(start + string(LAST_CONNECTOR))
}

func (c IncompleteComparison) IsComplete() bool {
	if c.Left != nil && c.Operator != nil && c.Right != nil {
		return true
	}
	return false
}

func (c IncompleteComparison) IsEmpty() bool {
	if c.Left == nil && c.Operator == nil && c.Right == nil {
		return true
	}
	return false
}

func (c IncompleteComparison) IsLeftComplete() bool {
	if c.Left != nil /* && c.Left.isComplete() */ {
		return true
	}
	return false
}

func (c IncompleteComparison) IsRightComplete() bool {
	if c.Right != nil && c.Right.IsComplete() {
		return true
	}
	return false
}

func (c IncompleteComparison) IsOperatorComplete() bool {
	if c.Operator != nil {
		return true
	}
	return false
}

type IncompleteTerm struct {
	Left     *IncompleteFactor
	Operator *Token
	Right    *IncompleteTerm
}

func (t IncompleteTerm) Print(start string) {
	if t.skipPrinting() {
		t.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("IncompleteTerm (%s)", friendlyOperatorName(t.Operator, false))
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	t.Left.Print(start + string(BRANCH_CONNECTOR))
	t.Right.Print(start + string(LAST_CONNECTOR))
}

func (t IncompleteTerm) IsComplete() bool {
	if t.Left != nil && t.Operator != nil && t.Right != nil {
		return true
	}
	return false
}

func (t IncompleteTerm) IsEmpty() bool {
	if t.Left == nil && t.Operator == nil && t.Right == nil {
		return true
	}
	return false
}

func (t IncompleteTerm) IsLeftComplete() bool {
	if t.Left != nil /* && t.Left.isComplete() */ {
		return true
	}
	return false
}

func (t IncompleteTerm) IsRightComplete() bool {
	if t.Right != nil && t.Right.IsComplete() {
		return true
	}
	return false
}

func (t IncompleteTerm) IsOperatorComplete() bool {
	if t.Operator != nil {
		return true
	}
	return false
}

type IncompleteFactor struct {
	Left     *IncompleteUnary
	Operator *Token
	Right    *IncompleteFactor
}

func (f IncompleteFactor) IsComplete() bool {
	if f.Left != nil && f.Operator != nil && f.Right != nil {
		return true
	}
	return false
}
func (f IncompleteFactor) IsEmpty() bool {
	if f.Left == nil && f.Operator == nil && f.Right == nil {
		return true
	}
	return false
}
func (f IncompleteFactor) IsLeftComplete() bool {
	if f.Left != nil /* && f.Left.isComplete() */ {
		return true
	}
	return false
}
func (f IncompleteFactor) IsRightComplete() bool {
	if f.Right != nil && f.Right.IsComplete() {
		return true
	}
	return false
}
func (f IncompleteFactor) IsOperatorComplete() bool {
	if f.Operator != nil {
		return true
	}
	return false
}

func (f IncompleteFactor) Print(start string) {
	if f.skipPrinting() {
		(*f.Left).Print(start)
		return
	}

	nodeName := fmt.Sprintf("IncompleteFactor (%s)", friendlyOperatorName(f.Operator, false))
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	(*f.Left).Print(start + string(BRANCH_CONNECTOR))
	(*f.Right).Print(start + string(LAST_CONNECTOR))
}

type IncompleteUnary interface {
	Print(start string)
	IsComplete() bool
	IsEmpty() bool
	IsRightComplete() bool
}

type IncompleteUnaryWithOperator struct {
	Operator *Token
	Right    *IncompleteUnary
}

func (uwo IncompleteUnaryWithOperator) Print(start string) {
	nodeName := fmt.Sprintf("IncompleteUnary (%s)", friendlyOperatorName(uwo.Operator, true))
	fmt.Printf("%s%s\n", start, Colorize(nodeName, COLOR_MAGENTA))
	start = advanceSuffix(start)
	(*uwo.Right).Print(start + string(LAST_CONNECTOR))
}

func (uwo IncompleteUnaryWithOperator) IsComplete() bool {
	if uwo.Operator != nil && uwo.Right != nil {
		return true
	}
	return false
}

func (uwo IncompleteUnaryWithOperator) IsEmpty() bool {
	if uwo.Operator == nil && uwo.Right == nil {
		return true
	}
	return false
}

func (uwo IncompleteUnaryWithOperator) IsRightComplete() bool {
	return uwo.Right != nil
}

type IncompletePrimary struct {
	Value PrimaryValue
}

func (ip IncompletePrimary) Print(start string) {
	if ip.Value == nil {
		fmt.Printf("%s%s\n", start, Colorize("IncompletePrimary (empty)", COLOR_RED))
		return
	}
}

func (ip IncompletePrimary) IsEmpty() bool {
	return ip.Value == nil
}

func (ip IncompletePrimary) IsComplete() bool {
	return ip.Value != nil
}

func (ip IncompletePrimary) IsRightComplete() bool {
	return ip.Value != nil
}

type IncompleteGroupingExpression struct {
	Expression *IncompleteExpression
}

func (ge IncompleteGroupingExpression) Print(start string) {
	fmt.Printf("%s%s\n", start, Colorize("IncompleteGroupingExpression", COLOR_GREEN))
	start = advanceSuffix(start)
	ge.Expression.Print(start + string(LAST_CONNECTOR))
}

type IncompleteCall struct {
	Callee    *string
	Arguments *[]IncompleteExpression
}

func (c IncompleteCall) Print(start string) {
	if c.Callee == nil {
		fmt.Printf("%s%s\n", start, Colorize("IncompleteCall (missing callee)", COLOR_RED))
		return
	}

	fmt.Printf("%s%s%s\n", start, string(BRANCH_CONNECTOR), Colorize(fmt.Sprintf("Call (%s)", *c.Callee), COLOR_GREEN))
	start = advanceSuffix(start)

	if c.Arguments == nil || len(*c.Arguments) == 0 {
		fmt.Printf("%s%s\n", start, Colorize("No arguments", COLOR_YELLOW))
		return
	}

	for i, arg := range *c.Arguments {
		if i == len(*c.Arguments)-1 {
			arg.Print(start + string(LAST_CONNECTOR))
		} else {
			arg.Print(start + string(SIMPLE_CONNECTOR))
		}
	}
}
