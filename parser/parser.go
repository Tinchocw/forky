package parser

import (
	"fmt"
	"slices"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
	primaryExpression "github.com/Tinchocw/Interprete-concurrente/common/expression/primary"
	"github.com/Tinchocw/Interprete-concurrente/common/statement"
)

type Parser struct {
	tokens          []common.Token
	current         int
	CouldMergeStart bool
	CouldMergeEnd   bool
	debug           bool
}

func NewParser(tokens []common.Token, debug bool) *Parser {
	return &Parser{tokens: tokens, current: 0, CouldMergeStart: true, CouldMergeEnd: true, debug: debug}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) peek() common.Token {
	return p.tokens[p.current]
}

// func (p *Parser) previous() common.Token {
// 	return p.tokens[p.current-1]
// }

func (p *Parser) check(posible_types ...common.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return slices.Contains(posible_types, p.peek().Typ)
}

func (p *Parser) checkAll(types ...common.TokenType) bool {
	if p.current+len(types) > len(p.tokens) {
		return false
	}
	for i, t := range types {
		if p.tokens[p.current+i].Typ != t {
			return false
		}
	}
	return true
}

func (p *Parser) advance() common.Token {
	if p.isAtEnd() {
		panic("No more tokens")
	}
	token := p.peek()
	p.current++
	return token
}

func (p *Parser) match(token common.TokenType) bool {
	if p.check(token) {
		p.advance()
		return true
	}
	return false
}

// return true if all the tokens match sequentially if not, return false and do not consume any token
func (p *Parser) matchs(token_types ...common.TokenType) bool {
	if p.checkAll(token_types...) {
		for range token_types {
			p.advance()
		}
		return true
	}
	return false
}

// STATEMENTS

func (p *Parser) program() (statement.Program, error) {
	var statements []statement.Statement
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
		}
		statements = append(statements, stmt)
	}

	return statement.Program{Statements: statements}, nil
}

func (p *Parser) blockStatement() (*statement.BlockStatement, error) {

	if !p.match(common.OPEN_BRACES) {
		return &statement.BlockStatement{}, fmt.Errorf("expected '{' at the beginning of block")
	}

	var statements []statement.Statement
	for !p.isAtEnd() && !p.check(common.CLOSE_BRACES) {
		stmt, err := p.statement()
		if err != nil {
			return &statement.BlockStatement{}, err
		}
		statements = append(statements, stmt)
	}

	if !p.match(common.CLOSE_BRACES) {
		return &statement.BlockStatement{}, fmt.Errorf("expected '}' at the end of block")
	}

	return &statement.BlockStatement{Statements: statements}, nil
}

func (p *Parser) statement() (statement.Statement, error) {

	if p.isAtEnd() {
		return nil, fmt.Errorf("unexpected end of input")
	}

	if p.debug {
		fmt.Printf("[DEBUG] Parsing statement at token %d: %s\n", p.current, p.peek().String())
	}

	switch token := p.peek(); token.Typ {
	case common.RETURN:
		return p.returnStatement()
	case common.PRINT:
		return p.printStatement()
	case common.IF:
		return p.ifStatement()
	case common.BREAK:
		return p.breakStatement()
	case common.FUNC:
		return p.funcStatement()
	case common.VAR:
		return p.varStatement()
	case common.WHILE:
		return p.whileStatement()
	case common.OPEN_BRACES:
		return p.blockStatement()

	default:
		if p.checkAll(common.IDENTIFIER, common.EQUAL) {
			return p.assignmentStatement()
		}

		return p.expressionStatement()
	}
}

func (p *Parser) printStatement() (*statement.PrintStatement, error) {
	if !p.match(common.PRINT) {
		return &statement.PrintStatement{}, fmt.Errorf("expected 'print'")
	}

	if !p.match(common.OPEN_PARENTHESIS) {
		return &statement.PrintStatement{}, fmt.Errorf("expected '(' after 'print'")
	}

	expr, err := p.expression()
	if err != nil {
		return &statement.PrintStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return &statement.PrintStatement{}, fmt.Errorf("expected ')' after expression")
	}

	if !p.match(common.SEMICOLON) {
		return &statement.PrintStatement{}, fmt.Errorf("expected ';' after print statement")
	}

	return &statement.PrintStatement{Value: expr}, nil
}

func (p *Parser) ifStatement() (*statement.IfStatement, error) {
	if !p.match(common.IF) {
		return &statement.IfStatement{}, fmt.Errorf("expected 'if'")
	}

	if !p.match(common.OPEN_PARENTHESIS) {
		return &statement.IfStatement{}, fmt.Errorf("expected '(' after 'if'")
	}

	condition, err := p.expression()
	if err != nil {
		return &statement.IfStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return &statement.IfStatement{}, fmt.Errorf("expected ')' after if condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return &statement.IfStatement{}, err
	}

	ifStatement := &statement.IfStatement{Condition: condition, Body: body}

	if p.matchs(common.ELSE, common.IF) {
		elseIf, err := p.elseIfStatement()
		if err != nil {
			return &statement.IfStatement{}, err
		}

		ifStatement.ElseIf = elseIf
	}

	if p.match(common.ELSE) {
		elseBody, err := p.blockStatement()
		if err != nil {
			return &statement.IfStatement{}, err
		}
		ifStatement.Else = &statement.ElseStatement{Body: elseBody}
	}

	return ifStatement, nil
}

func (p *Parser) elseIfStatement() (*statement.ElseIfStatement, error) {
	if !p.match(common.OPEN_PARENTHESIS) {
		return &statement.ElseIfStatement{}, fmt.Errorf("expected '(' after 'else if'")
	}

	condition, err := p.expression()
	if err != nil {
		return &statement.ElseIfStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return &statement.ElseIfStatement{}, fmt.Errorf("expected ')' after else if condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return &statement.ElseIfStatement{}, err
	}

	if p.matchs(common.ELSE, common.IF) {
		elseIf, err := p.elseIfStatement()
		if err != nil {
			return &statement.ElseIfStatement{}, err
		}
		return &statement.ElseIfStatement{Condition: condition, Body: body, ElseIf: elseIf}, nil
	}

	return &statement.ElseIfStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) breakStatement() (*statement.BreakStatement, error) {
	if !p.match(common.BREAK) {
		return &statement.BreakStatement{}, fmt.Errorf("expected 'break'")
	}
	if !p.match(common.SEMICOLON) {
		return &statement.BreakStatement{}, fmt.Errorf("expected ';' after 'break'")
	}
	return &statement.BreakStatement{}, nil
}

func (p *Parser) returnStatement() (*statement.ReturnStatement, error) {
	if !p.match(common.RETURN) {
		return &statement.ReturnStatement{}, fmt.Errorf("expected 'return'")
	}

	expr, err := p.expression()
	if err != nil {
		return &statement.ReturnStatement{}, err
	}
	if !p.match(common.SEMICOLON) {
		return &statement.ReturnStatement{}, fmt.Errorf("expected ';' after 'return'")
	}
	return &statement.ReturnStatement{Value: expr}, nil
}

func (p *Parser) funcStatement() (*statement.FunctionDef, error) {
	if !p.match(common.FUNC) {
		return &statement.FunctionDef{}, fmt.Errorf("expected 'func'")
	}

	if !p.check(common.IDENTIFIER) {
		return &statement.FunctionDef{}, fmt.Errorf("expected function name after 'func'")
	}
	name := p.advance()

	if !p.match(common.OPEN_PARENTHESIS) {
		return &statement.FunctionDef{}, fmt.Errorf("expected '(' after function name")
	}

	var parameters []string

	if !p.match(common.CLOSE_PARENTHESIS) {
		for {
			if !p.check(common.IDENTIFIER) {
				return &statement.FunctionDef{}, fmt.Errorf("expected parameter name")
			}
			parameters = append(parameters, p.advance().Value)

			if p.match(common.CLOSE_PARENTHESIS) {
				break
			}

			if !p.match(common.COMMA) {
				return &statement.FunctionDef{}, fmt.Errorf("expected ',' or ')' after parameter")
			}
		}
	}

	body, err := p.blockStatement()
	if err != nil {
		return &statement.FunctionDef{}, err
	}

	return &statement.FunctionDef{Name: &name.Value, Parameters: parameters, Body: body}, nil
}

func (p *Parser) whileStatement() (*statement.WhileStatement, error) {
	if !p.match(common.WHILE) {
		return &statement.WhileStatement{}, fmt.Errorf("expected 'while' at the beginning of while statement")
	}

	if !p.match(common.OPEN_PARENTHESIS) {
		return &statement.WhileStatement{}, fmt.Errorf("expected '(' after 'while'")
	}

	condition, err := p.expression()
	if err != nil {
		return &statement.WhileStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return &statement.WhileStatement{}, fmt.Errorf("expected ')' after while condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return &statement.WhileStatement{}, err
	}

	return &statement.WhileStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) assignmentStatement() (*statement.Assignment, error) {
	if !p.check(common.IDENTIFIER) {
		return &statement.Assignment{}, fmt.Errorf("expected variable name")
	}

	name := p.advance().Value

	if !p.match(common.EQUAL) {
		return &statement.Assignment{}, fmt.Errorf("expected '=' after variable name")
	}

	value, err := p.expression()
	if err != nil {
		return &statement.Assignment{}, err
	}

	if !p.match(common.SEMICOLON) {
		return &statement.Assignment{}, fmt.Errorf("expected ';' after assignment")
	}

	return &statement.Assignment{Name: &name, Value: value}, nil
}

func (p *Parser) expressionStatement() (*statement.ExpressionStatement, error) {
	expr, err := p.expression()
	if err != nil {
		return &statement.ExpressionStatement{}, err
	}

	if p.match(common.SEMICOLON) {
		// It's supposed to be a valid expression statement in the right side
		p.CouldMergeEnd = false
	}

	return &statement.ExpressionStatement{Expression: expr}, nil
}

func (p *Parser) varStatement() (*statement.VarDeclaration, error) {
	if !p.match(common.VAR) {
		return &statement.VarDeclaration{}, fmt.Errorf("expected 'var' at the beginning of variable declaration")
	}

	if !p.check(common.IDENTIFIER) {
		return &statement.VarDeclaration{}, fmt.Errorf("expected variable name after '%s'", common.VAR_KEYWORD)
	}

	name := p.advance()

	if !p.match(common.EQUAL) {
		return &statement.VarDeclaration{}, fmt.Errorf("expected '=' after variable name")
	}

	value, err := p.expression()
	if err != nil {
		return &statement.VarDeclaration{}, err
	}

	if !p.match(common.SEMICOLON) {
		return &statement.VarDeclaration{}, fmt.Errorf("expected ';' after variable declaration")
	}

	return &statement.VarDeclaration{Name: &name.Value, Value: value}, nil
}

// EXPRESIONES

func (p *Parser) expression() (*expression.Expression, error) {
	root, err := p.binaryOr()
	if err != nil {
		return &expression.Expression{}, err
	}
	return &expression.Expression{Root: root}, nil
}

func (p *Parser) binaryOr() (*expression.BinaryOr, error) {
	left, err := p.binaryAnd()
	if err != nil {
		return nil, err
	}

	bor := &expression.BinaryOr{Left: left}
	first := true

	for p.check(common.OR) {
		if !first {
			bor = &expression.BinaryOr{Left: bor}
		}

		operator := p.advance()
		right, err := p.binaryAnd()
		if err != nil {
			return nil, err
		}
		bor.Operator = &operator
		bor.Right = right

		first = false
	}

	return bor, nil
}

func (p *Parser) binaryAnd() (*expression.BinaryAnd, error) {
	left, err := p.equality()
	if err != nil {
		return nil, err
	}

	band := &expression.BinaryAnd{Left: left}
	first := true

	for p.check(common.AND) {
		if !first {
			band = &expression.BinaryAnd{Left: band}
		}

		operator := p.advance()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		band.Operator = &operator
		band.Right = right

		first = false
	}

	return band, nil
}

func (p *Parser) equality() (*expression.Equality, error) {
	left, err := p.comparison()
	if err != nil {
		return nil, err
	}

	eq := &expression.Equality{Left: left}
	first := true

	for p.check(common.BANG_EQUAL, common.EQUAL_EQUAL) {
		if !first {
			eq = &expression.Equality{Left: eq}
		}

		operator := p.advance()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		eq.Operator = &operator
		eq.Right = right

		first = false
	}

	return eq, nil
}

func (p *Parser) comparison() (*expression.Comparison, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}

	comp := &expression.Comparison{Left: left}
	first := true

	for p.check(common.GREATER, common.GREATER_EQUAL, common.LESS, common.LESS_EQUAL) {
		if !first {
			comp = &expression.Comparison{Left: comp}
		}

		operator := p.advance()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		comp.Operator = &operator
		comp.Right = right

		first = false
	}

	return comp, nil
}

func (p *Parser) term() (*expression.Term, error) {
	left, err := p.factor()
	if err != nil {
		return nil, err
	}

	t := &expression.Term{Left: left}
	first := true

	for p.check(common.MINUS, common.PLUS) {
		if !first {
			t = &expression.Term{Left: t}
		}

		operator := p.advance()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		t.Operator = &operator
		t.Right = right

		first = false
	}

	return t, nil
}

func (p *Parser) factor() (*expression.Factor, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}

	f := &expression.Factor{Left: left}
	first := true

	for p.check(common.SLASH, common.ASTERISK) {
		if !first {
			f = &expression.Factor{Left: f}
		}

		operator := p.advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		f.Operator = &operator
		f.Right = right

		first = false
	}

	return f, nil
}

func (p *Parser) unary() (expression.Unary, error) {
	if p.check(common.BANG, common.TILDE) {
		operator := p.advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &expression.UnaryWithOperator{Operator: &operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (expression.Unary, error) {
	if p.check(common.FALSE, common.TRUE, common.NONE, common.NUMBER, common.LITERAL) {
		token := p.advance()
		return &primaryExpression.Primary{Value: &primaryExpression.TokenValue{Token: &token}}, nil
	}

	if p.match(common.OPEN_PARENTHESIS) {
		if p.isAtEnd() {
			expr := &expression.Expression{}
			return &primaryExpression.Primary{Value: &primaryExpression.GroupingExpression{Expression: expr, StartingParenthesis: true, ClosingParenthesis: false}}, nil
		}
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(common.CLOSE_PARENTHESIS) {
			if p.isAtEnd() {
				return &primaryExpression.Primary{Value: &primaryExpression.GroupingExpression{Expression: expr, StartingParenthesis: true, ClosingParenthesis: false}}, nil
			}
			return nil, fmt.Errorf("expected ')' after expression")
		}

		return &primaryExpression.Primary{Value: &primaryExpression.GroupingExpression{Expression: expr, StartingParenthesis: true, ClosingParenthesis: true}}, nil
	}

	if p.match(common.OPEN_BRACES) {
		elements := []*expression.Expression{}
		for !p.match(common.CLOSE_BRACES) {
			element, err := p.expression()
			if err != nil {
				return nil, err
			}
			elements = append(elements, element)
			if !p.match(common.COMMA) {
				if p.check(common.CLOSE_BRACES) {
					break
				}

				return nil, fmt.Errorf("expected ',' or '}' after array element")
			}
		}

		return &primaryExpression.Primary{Value: &primaryExpression.ArrayLiteral{Elements: elements}}, nil
	}

	if p.check(common.IDENTIFIER) {
		identifier := p.advance()

		if p.match(common.OPEN_PARENTHESIS) {
			var args []*expression.Expression
			if !p.match(common.CLOSE_PARENTHESIS) {
				for {
					arg, err := p.expression()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)

					if p.match(common.CLOSE_PARENTHESIS) {
						break
					}

					if !p.match(common.COMMA) {
						return nil, fmt.Errorf("expected ',' or ')' after argument")
					}
				}
			}

			return &primaryExpression.Primary{Value: &primaryExpression.Call{Callee: identifier.Value, Arguments: args}}, nil
		} else if p.check(common.OPEN_BRACKET_SYMBOL) {
			indexes := []*expression.Expression{}

			for p.match(common.CLOSE_BRACKET_SYMBOL) {
				if p.isAtEnd() {
					return nil, fmt.Errorf("expected expression inside brackets")
				}

				indexExpr, err := p.expression()
				if err != nil {
					return nil, err
				}
				indexes = append(indexes, indexExpr)

				if !p.match(common.CLOSE_BRACKET_SYMBOL) {
					return nil, fmt.Errorf("expected ']' after index expression")
				}
			}

			return &primaryExpression.Primary{Value: &primaryExpression.ArrayAccess{ArrayName: identifier.Value, Indexes: indexes}}, nil

		} else {
			return &primaryExpression.Primary{Value: &primaryExpression.TokenValue{Token: &identifier}}, nil
		}
	}

	if p.match(common.CLOSE_PARENTHESIS) {
		expr := &expression.Expression{}
		return &primaryExpression.Primary{Value: &primaryExpression.GroupingExpression{Expression: expr, StartingParenthesis: false, ClosingParenthesis: true}}, nil
	}

	return nil, fmt.Errorf("unexpected token: %v", p.peek().String())
}

// END UTILS

func (p *Parser) parse() (statement.Program, error) {
	if p.debug {
		fmt.Printf("[DEBUG] Starting parsing with %d tokens\n", len(p.tokens))
	}

	program, err := p.program()

	if p.debug {
		if err != nil {
			fmt.Printf("[DEBUG] Parsing failed: %v\n", err)
		} else {
			fmt.Printf("[DEBUG] Parsing completed successfully with %d statements\n", len(program.Statements))
		}
	}

	return program, err
}
