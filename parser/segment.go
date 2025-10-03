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

func (current *segment) mergeExpressions(leftAst, rightAst *common.Expression) {
	current.merge(leftAst.Root, rightAst.Root)
}

func (current *segment) merge(leftAst, rightAst common.MergableNode) {

	if leftAst.IsLeftComplete() && !leftAst.IsRightComplete() && !leftAst.IsOperatorComplete() &&
		rightAst.IsLeftComplete() && !rightAst.IsRightComplete() && !rightAst.IsOperatorComplete() {
		current.merge(leftAst.GetLeft(), rightAst.GetLeft())
	}

	if leftAst.IsOperatorComplete() && !leftAst.IsRightComplete() && (rightAst.IsComplete() || rightAst.IsLeftComplete() && !rightAst.IsRightComplete() && !rightAst.IsOperatorComplete()) {
		leftAst.SetRight(rightAst)
		return
	}

	if leftAst.IsLeftComplete() && !leftAst.IsRightComplete() && !leftAst.IsOperatorComplete() && !rightAst.IsLeftComplete() && rightAst.IsOperatorComplete() {
		leftAst.SetOperator(rightAst.GetOperator())
		leftAst.SetRight(rightAst.GetRight())
		return
	}

	if leftAst.IsComplete() {
		current.merge(leftAst.GetRight(), rightAst)
		return
	}

	if leftAst.IsLeftComplete() && !leftAst.IsRightComplete() && !leftAst.IsOperatorComplete() && rightAst.IsComplete() {
		current.merge(leftAst.GetLeft(), rightAst.GetLeft())
		leftAst.SetOperator(rightAst.GetOperator())
		leftAst.SetRight(rightAst.GetRight())
		return
	}

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
