package parser

import (
	"fmt"

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

func (current *segment) mergeExpressions(leftAst, rightAst *common.Expression) error {
	leftDepth := leftAst.ExpresionDepth(true)
	rightDepth := rightAst.ExpresionDepth(false)

	if leftDepth > rightDepth {
		subExpr := leftAst.GetSubExpression(leftDepth-rightDepth, true)
		return current.merge(subExpr, rightAst)
	} else if rightDepth > leftDepth {
		subExpr := rightAst.GetSubExpression(rightDepth-leftDepth, false)
		return current.merge(leftAst, subExpr)
	} else {
		return current.merge(leftAst.Root, rightAst.Root)
	}

}

// 1. izquierda
// 2. izquierda + operador
// 3. operador + derecha
// 4. izquierda + operador + derecha
// 5. operador

func (current *segment) merge(leftAst, rightAst common.MergableNode) error {

	// IZQ = 3 o 4 | DER = *.
	if leftAst.HasRight() {
		return current.merge(leftAst.GetRight(), rightAst)
	}

	// IZQ = 1 | DER = *
	if leftAst.HasLeft() && !leftAst.HasOperator() {
		// DER = 1 o 2 o 4
		if rightAst.HasLeft() {
			err := current.merge(leftAst.GetLeft(), rightAst.GetLeft())
			if err != nil {
				return err
			}
		}

		// DER = 3 o 4 o 5
		if rightAst.HasOperator() {
			leftAst.SetOperator(rightAst.GetOperator())

			if rightAst.HasRight() {
				leftAst.SetRight(rightAst.GetRight())
			}
		}

		return nil
	}

	// IZQ = 2 o 5 | DER = 1 o 2 o 4
	if leftAst.HasOperator() && rightAst.HasLeft() {
		leftAst.SetRight(rightAst)
		return nil
	}

	return fmt.Errorf("cannot merge expressions")
}

func (current *segment) Merge(other segment) error {
	defer func() {
		current.Tokens = append(current.Tokens, other.Tokens...)
	}()

	if current.CouldMergeEnd && other.CouldMergeStart {
		switch leftExpr := current.lastStatement().(type) {
		case common.ExpressionStatement:
			switch rightExpr := other.firstStatement().(type) {
			case common.ExpressionStatement:

				err := current.mergeExpressions(leftExpr.Expression, rightExpr.Expression)
				if err != nil {
					return err
				}
				current.CouldMergeEnd = other.CouldMergeEnd

				return nil
			}

		}

	}

	return nil
}
