package parser

import (
	"github.com/Tinchocw/Interprete-concurrente/common"
)

type segment struct {
	CouldMergeStart bool
	CouldMergeEnd   bool
	Program         common.Program
	Tokens          []common.Token
}

func NewSegment() segment {
	return segment{
		CouldMergeStart: true,
		CouldMergeEnd:   true,
		Program:         common.Program{Statements: []common.Statement{}},
		Tokens:          []common.Token{},
	}
}

func (current *segment) AddStatement(content common.Statement) {
	current.Program.Statements = append(current.Program.Statements, content)
}
func (current *segment) AddStatements(contents []common.Statement) {
	current.Program.Statements = append(current.Program.Statements, contents...)
}

func (current *segment) hasStatements() bool {
	return len(current.Program.Statements) > 0
}

func (current *segment) firstStatement() common.Statement {
	if len(current.Program.Statements) > 0 {
		return current.Program.Statements[0]
	}
	return nil
}

func (current *segment) lastStatement() common.Statement {
	if len(current.Program.Statements) > 0 {
		return current.Program.Statements[len(current.Program.Statements)-1]
	}
	return nil
}

func (current *segment) mergeExpressions(leftNode, rightNode *common.Expression) {

	current.mergeBinaryOr(leftNode.Root, rightNode.Root)
}

func (current *segment) mergeBinaryOr(leftAst, rightAst *common.BinaryOr) {

	if leftAst.IsLeftComplete() && !leftAst.IsRightComplete() && !leftAst.IsOperatorComplete() && !rightAst.IsLeftComplete() && rightAst.IsRightComplete() {
		current.mergeBinaryAnd(leftAst.Left, rightAst.Left)
	}

	// LeftNode with left and without right
	if leftAst.IsLeftComplete() && !leftAst.IsRightComplete() {
		leftAst.Right = rightAst
	}

	if leftAst.IsLeftComplete() && leftAst.IsRightComplete() && rightAst.IsLeftComplete() && !rightAst.IsRightComplete() {
		current.mergeBinaryAnd(leftAst.Left, rightAst.Left)
	}

	// both size complete
	if leftAst.IsComplete() && rightAst.IsComplete() {
		current.mergeBinaryOr(leftAst.Right, rightAst)
	}

}

func (current *segment) mergeBinaryAnd(leftNode, rightNode *common.BinaryAnd) {

	if leftNode.IsLeftComplete() && !leftNode.IsRightComplete() {
		leftNode.Right = rightNode
	}

	if leftNode.IsLeftComplete() && !leftNode.IsRightComplete() && !rightNode.IsLeftComplete() && rightNode.IsRightComplete() {
		current.mergeEquality(leftNode.Left, rightNode.Left)
	}

	if leftNode.IsComplete() && rightNode.IsComplete() {
		current.mergeBinaryAnd(leftNode.Right, rightNode)
	}
}

func (current *segment) mergeEquality(leftNode, rightNode *common.Equality) {

	if leftNode.IsLeftComplete() && leftNode.IsOperatorComplete() && !leftNode.IsRightComplete() {
		leftNode.Right = rightNode
	}

	if leftNode.IsLeftComplete() && !leftNode.IsRightComplete() && !rightNode.IsLeftComplete() && rightNode.IsRightComplete() && leftNode.IsOperatorComplete() && !rightNode.IsOperatorComplete() {
		current.mergeComparison(leftNode.Left, rightNode.Left)
	}

	if leftNode.IsComplete() && rightNode.IsComplete() {
		current.mergeEquality(leftNode.Right, rightNode)
	}
}

func (current *segment) mergeComparison(leftNode, rightNode *common.Comparison) {

	if leftNode.IsLeftComplete() && !leftNode.IsRightComplete() {

		leftNode.Right = rightNode
	}

	if leftNode.IsLeftComplete() && !leftNode.IsRightComplete() && !rightNode.IsLeftComplete() && rightNode.IsRightComplete() {
		current.mergeTerm(leftNode.Left, rightNode.Left)
	}

	if leftNode.IsComplete() && rightNode.IsComplete() {
		current.mergeComparison(leftNode.Right, rightNode)
	}
}

func (current *segment) mergeTerm(leftNode, rightNode *common.Term) {

	if leftNode.IsLeftComplete() && leftNode.IsOperatorComplete() && !leftNode.IsRightComplete() {
		leftNode.Right = rightNode
	}

	if leftNode.IsLeftComplete() && leftNode.IsOperatorComplete() && !leftNode.IsRightComplete() && !rightNode.IsLeftComplete() && rightNode.IsRightComplete() && rightNode.IsOperatorComplete() {
		current.mergeFactor(leftNode.Left, rightNode.Left)
	}

	if leftNode.IsComplete() && rightNode.IsComplete() {
		current.mergeTerm(leftNode.Right, rightNode)
	}
}

func (current *segment) mergeFactor(leftNode, rightNode *common.Factor) {

	if leftNode.IsLeftComplete() && leftNode.IsOperatorComplete() && !leftNode.IsRightComplete() {
		leftNode.Right = rightNode
	}

	// TODO: See all the posible cases
	if leftNode.IsLeftComplete() && leftNode.IsOperatorComplete() && !leftNode.IsRightComplete() && !rightNode.IsLeftComplete() && rightNode.IsRightComplete() && rightNode.IsOperatorComplete() {
		current.mergeUnary(leftNode.Left)
	}

	if leftNode.IsComplete() && rightNode.IsComplete() {
		current.mergeFactor(leftNode.Right, rightNode)
	}
}

func (current *segment) mergeUnary(rightNode common.Unary) {
	switch rightNode := rightNode.(type) {
	case *common.UnaryWithOperator:
		current.mergeUnaryWithOperator(rightNode)
	case *common.Primary:
		current.mergePrimary(rightNode)
	}
}

func (current *segment) mergeUnaryWithOperator(rightNode *common.UnaryWithOperator) {
}

func (current *segment) mergePrimary(leftNode *common.Primary) {

}

func (current *segment) Merge(other segment) {
	defer func() {
		current.Tokens = append(current.Tokens, other.Tokens...)
	}()

	if current.CouldMergeEnd && other.CouldMergeStart {

		switch leftExpr := current.lastStatement().(type) {
		case common.ExpressionStatement:
			switch rightExpr := other.firstStatement().(type) {
			case common.ExpressionStatement:

				current.mergeExpressions(leftExpr.Expression, rightExpr.Expression)
				current.CouldMergeEnd = other.CouldMergeEnd
			}

		}

	}
}
