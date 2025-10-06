package parser

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
	"github.com/Tinchocw/Interprete-concurrente/common/statement"
)

type segment struct {
	CouldMergeStart bool
	CouldMergeEnd   bool
	Program         statement.Program
	Tokens          []common.Token
}

func NewSegment() segment {
	return segment{
		CouldMergeStart: true,
		CouldMergeEnd:   true,
		Program:         statement.Program{Statements: []statement.Statement{}},
		Tokens:          []common.Token{},
	}
}

func (current *segment) AddStatement(content statement.Statement) {
	current.Program.Statements = append(current.Program.Statements, content)
}
func (current *segment) AddStatements(contents []statement.Statement) {
	current.Program.Statements = append(current.Program.Statements, contents...)
}

func (current *segment) firstStatement() statement.Statement {
	if len(current.Program.Statements) > 0 {
		return current.Program.Statements[0]
	}
	return nil
}

func (current *segment) lastStatement() statement.Statement {
	if len(current.Program.Statements) > 0 {
		return current.Program.Statements[len(current.Program.Statements)-1]
	}
	return nil
}

// ( 2 + 3  | )

func (current *segment) mergeExpressions(leftAst, rightAst *expression.Expression) error {
	leftDepth := leftAst.ExpresionDepth(expression.Right)
	rightDepth := rightAst.ExpresionDepth(expression.Left)

	if leftDepth > rightDepth {
		subExpr := leftAst.GetSubExpression(leftDepth-rightDepth, expression.Right)
		return current.merge(subExpr, rightAst)
	} else if rightDepth > leftDepth {
		subExpr := rightAst.GetSubExpression(rightDepth-leftDepth, expression.Left)
		return current.merge(leftAst, subExpr)
	} else {
		return current.merge(leftAst, rightAst)
	}
}

// 1. izquierda
// 2. izquierda + operador
// 3. operador + derecha
// 4. izquierda + operador + derecha
// 5. operador

// grouping ( 2 + 3  | 1 )

// 2 + 1 ) | + 1 )
// ( | 1 + 1  | )

func (current *segment) merge(leftAst, rightAst expression.MergableNode) error {

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

	// Merge empty expression
	if rightAst.HasLeft() {
		leftAst.SetLeft(rightAst.GetLeft())
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
		case *statement.ExpressionStatement:
			switch rightExpr := other.firstStatement().(type) {
			case *statement.ExpressionStatement:
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
