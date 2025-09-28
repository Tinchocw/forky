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

type IncompleteFactor struct {
	Left     *IncompleteUnary
	Operator *Token
	Right    *IncompleteFactor
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

type IncompletePrimary struct {
	Value PrimaryValue
}

func (ip IncompletePrimary) Print(start string) {
	if ip.Value == nil {
		fmt.Printf("%s%s\n", start, Colorize("IncompletePrimary (empty)", COLOR_RED))
		return
	}
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
