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
	tokens  []common.Token
	current int
	// lastStatement   common.Statement
	// lastExpression  *expression.Expression
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

func (p *Parser) program() (segment, error) {
	var statements []statement.Statement
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return segment{}, err
		}
		statements = append(statements, stmt)
	}

	return segment{
		Program:         statement.Program{Statements: statements},
		Tokens:          p.tokens,
		CouldMergeStart: p.CouldMergeStart,
		CouldMergeEnd:   p.CouldMergeEnd,
	}, nil
}

func (p *Parser) blockStatement() (*statement.BlockStatement, error) {

	if !p.match(common.OPEN_BRACKET) {
		return &statement.BlockStatement{}, fmt.Errorf("expected '{' at the beginning of block")
	}

	var statements []statement.Statement
	for !p.isAtEnd() && !p.check(common.CLOSE_BRACKET) {
		stmt, err := p.statement()
		if err != nil {
			return &statement.BlockStatement{}, err
		}
		statements = append(statements, stmt)
	}

	if !p.match(common.CLOSE_BRACKET) {
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
	case common.OPEN_BRACKET:
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
	root, _, err := p.binaryOr()
	if err != nil {
		return &expression.Expression{}, err
	}

	expression := &expression.Expression{Root: root}

	for p.match(common.CLOSE_PARENTHESIS) {
		expression = createGroupingExpression(expression, false, true)
	}

	return expression, nil
}

func (p *Parser) binaryOr() (*expression.BinaryOr, bool, error) {
	if p.isAtEnd() {
		return nil, false, nil
	}

	if p.check(common.OR) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.BinaryOr{Left: nil, Operator: &operator, Right: nil}, true, nil
		}
		right, startWithOperator, err := p.binaryOr()
		if err != nil {
			return nil, false, err
		}
		if startWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}

		return &expression.BinaryOr{Left: nil, Operator: &operator, Right: right}, true, nil
	}

	left, startWithOperator, err := p.binaryAnd()
	if err != nil {
		return nil, false, err
	}

	if p.check(common.OR) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.BinaryOr{Left: left, Operator: &operator, Right: nil}, startWithOperator, nil
		}
		right, startWithOperator, err := p.binaryOr()
		if err != nil {
			return nil, false, err
		}
		if startWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}

		return &expression.BinaryOr{Left: left, Operator: &operator, Right: right}, startWithOperator, nil
	}

	return &expression.BinaryOr{Left: left, Right: nil}, startWithOperator, nil
}

func (p *Parser) binaryAnd() (*expression.BinaryAnd, bool, error) {

	if p.check(common.AND) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.BinaryAnd{Left: nil, Operator: &operator, Right: nil}, true, nil
		}
		right, startWithOperator, err := p.binaryAnd()
		if err != nil {
			return nil, false, err
		}
		if startWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}
		return &expression.BinaryAnd{Left: nil, Operator: &operator, Right: right}, true, nil
	}

	left, startWithOperator, err := p.equality()
	if err != nil {
		return nil, false, err
	}

	if p.check(common.AND) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.BinaryAnd{Left: left, Operator: &operator, Right: nil}, startWithOperator, nil
		}
		right, rightStartWithOperator, err := p.binaryAnd()
		if err != nil {
			return nil, false, err
		}
		if rightStartWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}

		return &expression.BinaryAnd{Left: left, Operator: &operator, Right: right}, startWithOperator, nil
	}

	return &expression.BinaryAnd{Left: left, Operator: nil, Right: nil}, startWithOperator, nil
}

func (p *Parser) equality() (*expression.Equality, bool, error) {

	if p.check(common.BANG_EQUAL, common.EQUAL_EQUAL) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.Equality{Left: nil, Operator: &operator, Right: nil}, true, nil
		}
		right, startWithOperator, err := p.equality()
		if err != nil {
			return nil, false, err
		}
		if startWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}
		return &expression.Equality{Left: nil, Operator: &operator, Right: right}, true, nil
	}

	left, startWithOperator, err := p.comparison()
	if err != nil {
		return nil, false, err
	}

	if p.check(common.BANG_EQUAL, common.EQUAL_EQUAL) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.Equality{Left: left, Operator: &operator, Right: nil}, startWithOperator, nil
		}
		right, rightStartWithOperator, err := p.equality()
		if err != nil {
			return nil, false, err
		}
		if rightStartWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}

		return &expression.Equality{Left: left, Operator: &operator, Right: right}, startWithOperator, nil
	}

	return &expression.Equality{Left: left, Operator: nil, Right: nil}, startWithOperator, nil
}

func (p *Parser) comparison() (*expression.Comparison, bool, error) {

	if p.check(common.GREATER, common.GREATER_EQUAL, common.LESS, common.LESS_EQUAL) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.Comparison{Left: nil, Operator: &operator, Right: nil}, true, nil
		}
		right, startWithOperator, err := p.comparison()
		if err != nil {
			return nil, false, err
		}
		if startWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}
		return &expression.Comparison{Left: nil, Operator: &operator, Right: right}, true, nil
	}

	left, startWithOperator, err := p.term()
	if err != nil {
		return nil, false, err
	}

	if p.check(common.GREATER, common.GREATER_EQUAL, common.LESS, common.LESS_EQUAL) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.Comparison{Left: left, Operator: &operator, Right: nil}, startWithOperator, nil
		}
		right, rightStartWithOperator, err := p.comparison()
		if err != nil {
			return nil, false, err
		}
		if rightStartWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}

		return &expression.Comparison{Left: left, Operator: &operator, Right: right}, startWithOperator, nil
	}

	return &expression.Comparison{Left: left, Operator: nil, Right: nil}, startWithOperator, nil
}

func (p *Parser) term() (*expression.Term, bool, error) {

	if p.check(common.MINUS, common.PLUS) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.Term{Left: nil, Operator: &operator, Right: nil}, true, nil
		}

		right, startWithOperator, err := p.term()
		if err != nil {
			return nil, false, err
		}

		if startWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}

		return &expression.Term{Left: nil, Operator: &operator, Right: right}, true, nil
	}

	left, startWithOperator, err := p.factor()
	if err != nil {
		return nil, false, err
	}

	if p.check(common.MINUS, common.PLUS) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.Term{Left: left, Operator: &operator, Right: nil}, startWithOperator, nil
		}

		right, rightStartWithOperator, err := p.term()
		if err != nil {
			return nil, false, err
		}
		if rightStartWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}

		return &expression.Term{Left: left, Operator: &operator, Right: right}, startWithOperator, nil
	}
	return &expression.Term{Left: left, Operator: nil, Right: nil}, startWithOperator, nil
}

func (p *Parser) factor() (*expression.Factor, bool, error) {

	if p.check(common.SLASH, common.ASTERISK) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.Factor{Left: nil, Operator: &operator, Right: nil}, true, nil
		}
		right, startWithOperator, err := p.factor()
		if err != nil {
			return nil, false, err
		}
		if startWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}
		return &expression.Factor{Left: nil, Operator: &operator, Right: right}, true, nil
	}

	left, err := p.unary()
	if err != nil {
		return nil, false, err
	}

	if p.check(common.SLASH, common.ASTERISK) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.Factor{Left: left, Operator: &operator, Right: nil}, false, nil
		}
		right, rightStartWithOperator, err := p.factor()
		if err != nil {
			return nil, false, err
		}
		if rightStartWithOperator {
			return nil, false, fmt.Errorf("unexpected operator after %s", operator.String())
		}

		return &expression.Factor{Left: left, Operator: &operator, Right: right}, false, nil
	}

	return &expression.Factor{Left: left, Operator: nil, Right: nil}, false, nil
}

func (p *Parser) unary() (expression.Unary, error) {
	if p.check(common.BANG, common.TILDE) {
		operator := p.advance()
		if p.isAtEnd() {
			return &expression.UnaryWithOperator{Operator: &operator, Right: nil}, nil
		}
		right, err := p.unary()
		if err != nil {
			return &expression.UnaryWithOperator{}, err
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

			return &primaryExpression.Primary{Value: &primaryExpression.Call{Callee: &identifier.Value, Arguments: args}}, nil
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

// UTILS

func createGroupingExpression(expr *expression.Expression, startingParenthesis, closingParenthesis bool) *expression.Expression {
	return &expression.Expression{Root: &expression.BinaryOr{
		Left: &expression.BinaryAnd{
			Left: &expression.Equality{
				Left: &expression.Comparison{
					Left: &expression.Term{
						Left: &expression.Factor{
							Left: &primaryExpression.Primary{
								Value: &primaryExpression.GroupingExpression{
									Expression:          expr,
									StartingParenthesis: startingParenthesis,
									ClosingParenthesis:  closingParenthesis,
								},
							},
						},
					},
				},
			},
		},
	},
	}
}

// END UTILS

func (p *Parser) parse() (segment, error) {
	if p.debug {
		fmt.Printf("[DEBUG] Starting parsing with %d tokens\n", len(p.tokens))
	}

	sg, err := p.program()

	if p.debug {
		if err != nil {
			fmt.Printf("[DEBUG] Parsing failed: %v\n", err)
		} else {
			fmt.Printf("[DEBUG] Parsing completed successfully with %d statements\n", len(sg.Program.Statements))
		}
	}

	return sg, err
}
