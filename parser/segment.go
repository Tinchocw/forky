package parser

import "github.com/Tinchocw/Interprete-concurrente/common"

type segment struct {
	CouldMergeStart bool
	CouldMergeEnd   bool
	Program         common.PartialProgram
	Tokens          []common.Token
}

func NewSegment() segment {
	return segment{
		CouldMergeStart: true,
		CouldMergeEnd:   true,
		Program:         common.PartialProgram{Statements: []common.IncompleteStatement{}},
		Tokens:          []common.Token{},
	}
}

func (current *segment) AddStatement(content common.IncompleteStatement) {
	current.Program.Statements = append(current.Program.Statements, content)
}
func (current *segment) AddStatements(contents []common.IncompleteStatement) {
	current.Program.Statements = append(current.Program.Statements, contents...)
}

func (current *segment) hasStatements() bool {
	return len(current.Program.Statements) > 0
}

func (current *segment) hasInvalidStatements() bool {

	if !current.hasStatements() {
		return false
	}

	return false
}

func (current *segment) firstStatement() common.IncompleteStatement {
	if len(current.Program.Statements) > 0 {
		return current.Program.Statements[0]
	}
	return nil
}

func (current *segment) lastStatement() common.IncompleteStatement {
	if len(current.Program.Statements) > 0 {
		return current.Program.Statements[len(current.Program.Statements)-1]
	}
	return nil
}

func (current *segment) mergeExpressions(right common.IncompleteExpression) {
	current.mergeBinaryOr(current.lastStatement().(*common.IncompleteExpression).Root, right.Root)
}

func (current *segment) mergeBinaryOr(leftNode, rightNode *common.IncompleteBinaryOr) {

	// LeftNode with left and without right
	if leftNode.IsLeftComplete() && !leftNode.IsRightComplete() {
		leftNode.Right = rightNode
	}

	if leftNode.IsLeftComplete() && leftNode.IsRightComplete() && rightNode.IsLeftComplete() && !rightNode.IsRightComplete() {
		current.mergeBinaryAnd(leftNode.Left, rightNode.Left)
	}

	// both size complete
	if leftNode.IsComplete() && rightNode.IsComplete() {
		current.mergeBinaryOr(leftNode.Right, rightNode)
	}

}

func (current *segment) mergeBinaryAnd(leftNode, rightNode *common.IncompleteBinaryAnd) {

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

func (current *segment) mergeEquality(leftNode, rightNode *common.IncompleteEquality) {

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

func (current *segment) mergeComparison(leftNode, rightNode *common.IncompleteComparison) {

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

func (current *segment) mergeTerm(leftNode, rightNode *common.IncompleteTerm) {

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

func (current *segment) mergeFactor(leftNode, rightNode *common.IncompleteFactor) {

	if leftNode.IsLeftComplete() && leftNode.IsOperatorComplete() && !leftNode.IsRightComplete() {
		leftNode.Right = rightNode
	}

	// TODO: See all the posible cases
	if leftNode.IsLeftComplete() && leftNode.IsOperatorComplete() && !leftNode.IsRightComplete() && !rightNode.IsLeftComplete() && rightNode.IsRightComplete() && rightNode.IsOperatorComplete() {
		current.mergeUnary(*leftNode.Left)
	}

	if leftNode.IsComplete() && rightNode.IsComplete() {
		current.mergeFactor(leftNode.Right, rightNode)
	}
}

func (current *segment) mergeUnary(rightNode common.IncompleteUnary) {
	switch rightNode := rightNode.(type) {
	case *common.IncompleteUnaryWithOperator:
		current.mergeUnaryWithOperator(rightNode)
	case *common.IncompletePrimary:
		current.mergePrimary(rightNode)
	}
}

func (current *segment) mergeUnaryWithOperator(rightNode *common.IncompleteUnaryWithOperator) {
}

func (current *segment) mergePrimary(leftNode *common.IncompletePrimary) {

}

func (current *segment) Merge(other segment) segment {
	defer func() {
		current.Tokens = append(current.Tokens, other.Tokens...)
	}()

	if !current.hasStatements() && other.hasStatements() {
		return other
	}

	if !other.hasStatements() && current.hasStatements() {
		return *current
	}

	if current.CouldMergeEnd && other.CouldMergeStart {

		switch current.lastStatement().(type) {
		case *common.IncompleteExpression:
			switch rightExpr := other.firstStatement().(type) {
			case *common.IncompleteExpression:
				current.mergeExpressions(*rightExpr)
				//current.Program.Statements[len(current.Program.Statements)-1] = &left
				current.CouldMergeEnd = other.CouldMergeEnd
				return *current
			}
		}
	}

	current.AddStatements(other.Program.Statements)

	return *current
}
