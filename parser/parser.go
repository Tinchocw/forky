package parser

import (
	"fmt"
	"slices"

	"github.com/Tinchocw/Interprete-concurrente/common"
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

func (p *Parser) previous() common.Token {
	return p.tokens[p.current-1]
}

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
	var statements []common.IncompleteStatement
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return segment{}, err
		}
		statements = append(statements, stmt)
	}

	return segment{
		Program:         common.PartialProgram{Statements: statements},
		Tokens:          p.tokens,
		CouldMergeStart: p.CouldMergeStart,
		CouldMergeEnd:   p.CouldMergeEnd,
	}, nil
}

func (p *Parser) blockStatement() (common.IncompleteBlockStatement, error) {

	if !p.match(common.OPEN_BRACKET) {
		return common.IncompleteBlockStatement{}, fmt.Errorf("expected '{' at the beginning of block")
	}

	var statements []common.IncompleteStatement
	for !p.isAtEnd() && !p.check(common.CLOSE_BRACKET) {
		stmt, err := p.statement()
		if err != nil {
			return common.IncompleteBlockStatement{}, err
		}
		statements = append(statements, stmt)
	}

	if !p.match(common.CLOSE_BRACKET) {
		return common.IncompleteBlockStatement{}, fmt.Errorf("expected '}' at the end of block")
	}

	return common.IncompleteBlockStatement{Statements: statements}, nil
}

func (p *Parser) statement() (common.IncompleteStatement, error) {

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

func (p *Parser) printStatement() (common.IncompletePrintStatement, error) {
	if !p.match(common.PRINT) {
		return common.IncompletePrintStatement{}, fmt.Errorf("expected 'print'")
	}

	if !p.match(common.OPEN_PARENTHESIS) {
		return common.IncompletePrintStatement{}, fmt.Errorf("expected '(' after 'print'")
	}

	expr, err := p.expression()
	if err != nil {
		return common.IncompletePrintStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return common.IncompletePrintStatement{}, fmt.Errorf("expected ')' after expression")
	}

	if !p.match(common.SEMICOLON) {
		return common.IncompletePrintStatement{}, fmt.Errorf("expected ';' after print statement")
	}

	return common.IncompletePrintStatement{Value: &expr}, nil
}

func (p *Parser) ifStatement() (common.IncompleteIfStatement, error) {
	if !p.match(common.IF) {
		return common.IncompleteIfStatement{}, fmt.Errorf("expected 'if'")
	}

	if !p.match(common.OPEN_PARENTHESIS) {
		return common.IncompleteIfStatement{}, fmt.Errorf("expected '(' after 'if'")
	}

	condition, err := p.expression()
	if err != nil {
		return common.IncompleteIfStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return common.IncompleteIfStatement{}, fmt.Errorf("expected ')' after if condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return common.IncompleteIfStatement{}, err
	}

	ifStatement := common.IncompleteIfStatement{Condition: &condition, Body: &body}

	if p.matchs(common.ELSE, common.IF) {
		elseIf, err := p.elseIfStatement()
		if err != nil {
			return common.IncompleteIfStatement{}, err
		}

		ifStatement.ElseIf = &elseIf
	}

	if p.match(common.ELSE) {
		elseBody, err := p.blockStatement()
		if err != nil {
			return common.IncompleteIfStatement{}, err
		}
		ifStatement.Else = &common.IncompleteElseStatement{Body: &elseBody}
	}

	return ifStatement, nil
}

func (p *Parser) elseIfStatement() (common.IncompleteElseIfStatement, error) {
	if !p.match(common.OPEN_PARENTHESIS) {
		return common.IncompleteElseIfStatement{}, fmt.Errorf("expected '(' after 'else if'")
	}

	condition, err := p.expression()
	if err != nil {
		return common.IncompleteElseIfStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return common.IncompleteElseIfStatement{}, fmt.Errorf("expected ')' after else if condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return common.IncompleteElseIfStatement{}, err
	}

	if p.matchs(common.ELSE, common.IF) {
		elseIf, err := p.elseIfStatement()
		if err != nil {
			return common.IncompleteElseIfStatement{}, err
		}
		return common.IncompleteElseIfStatement{Condition: &condition, Body: &body, ElseIf: &elseIf}, nil
	}

	return common.IncompleteElseIfStatement{Condition: &condition, Body: &body}, nil
}

func (p *Parser) breakStatement() (common.IncompleteBreakStatement, error) {
	if !p.match(common.BREAK) {
		return common.IncompleteBreakStatement{}, fmt.Errorf("expected 'break'")
	}
	if !p.match(common.SEMICOLON) {
		return common.IncompleteBreakStatement{}, fmt.Errorf("expected ';' after 'break'")
	}
	return common.IncompleteBreakStatement{}, nil
}

func (p *Parser) returnStatement() (common.IncompleteReturnStatement, error) {
	if !p.match(common.RETURN) {
		return common.IncompleteReturnStatement{}, fmt.Errorf("expected 'return'")
	}

	expr, err := p.expression()
	if err != nil {
		return common.IncompleteReturnStatement{}, err
	}
	if !p.match(common.SEMICOLON) {
		return common.IncompleteReturnStatement{}, fmt.Errorf("expected ';' after 'return'")
	}
	return common.IncompleteReturnStatement{Value: &expr}, nil
}

func (p *Parser) funcStatement() (common.IncompleteFunctionDef, error) {
	if !p.match(common.FUNC) {
		return common.IncompleteFunctionDef{}, fmt.Errorf("expected 'func'")
	}

	if !p.check(common.IDENTIFIER) {
		return common.IncompleteFunctionDef{}, fmt.Errorf("expected function name after 'func'")
	}
	name := p.advance()

	if !p.match(common.OPEN_PARENTHESIS) {
		return common.IncompleteFunctionDef{}, fmt.Errorf("expected '(' after function name")
	}

	var parameters []string

	if !p.match(common.CLOSE_PARENTHESIS) {
		for {
			if !p.check(common.IDENTIFIER) {
				return common.IncompleteFunctionDef{}, fmt.Errorf("expected parameter name")
			}
			parameters = append(parameters, p.advance().Value)

			if p.match(common.CLOSE_PARENTHESIS) {
				break
			}

			if !p.match(common.COMMA) {
				return common.IncompleteFunctionDef{}, fmt.Errorf("expected ',' or ')' after parameter")
			}
		}
	}

	body, err := p.blockStatement()
	if err != nil {
		return common.IncompleteFunctionDef{}, err
	}

	return common.IncompleteFunctionDef{Name: &name.Value, Parameters: parameters, Body: &body}, nil
}

func (p *Parser) whileStatement() (common.IncompleteWhileStatement, error) {
	if !p.match(common.WHILE) {
		return common.IncompleteWhileStatement{}, fmt.Errorf("expected 'while' at the beginning of while statement")
	}

	if !p.match(common.OPEN_PARENTHESIS) {
		return common.IncompleteWhileStatement{}, fmt.Errorf("expected '(' after 'while'")
	}

	condition, err := p.expression()
	if err != nil {
		return common.IncompleteWhileStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return common.IncompleteWhileStatement{}, fmt.Errorf("expected ')' after while condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return common.IncompleteWhileStatement{}, err
	}

	return common.IncompleteWhileStatement{Condition: &condition, Body: &body}, nil
}

func (p *Parser) assignmentStatement() (common.IncompleteAssignment, error) {
	if !p.check(common.IDENTIFIER) {
		return common.IncompleteAssignment{}, fmt.Errorf("expected variable name")
	}

	name := p.advance().Value

	if !p.match(common.EQUAL) {
		return common.IncompleteAssignment{}, fmt.Errorf("expected '=' after variable name")
	}

	value, err := p.expression()
	if err != nil {
		return common.IncompleteAssignment{}, err
	}

	if !p.match(common.SEMICOLON) {
		return common.IncompleteAssignment{}, fmt.Errorf("expected ';' after assignment")
	}

	return common.IncompleteAssignment{Name: &name, Value: &value}, nil
}

func (p *Parser) expressionStatement() (common.IncompleteExpressionStatement, error) {
	expr, err := p.expression()
	if err != nil {
		return common.IncompleteExpressionStatement{}, err
	}

	if p.match(common.SEMICOLON) {
		// It's supposed to be a valid expression statement in the right side
		p.CouldMergeEnd = false
	}

	return common.IncompleteExpressionStatement{Expression: &expr}, nil
}

func (p *Parser) varStatement() (common.IncompleteVarDeclaration, error) {
	if !p.match(common.VAR) {
		return common.IncompleteVarDeclaration{}, fmt.Errorf("expected 'var' at the beginning of variable declaration")
	}

	if !p.check(common.IDENTIFIER) {
		return common.IncompleteVarDeclaration{}, fmt.Errorf("expected variable name after '%s'", common.VAR_KEYWORD)
	}

	name := p.advance()

	if !p.match(common.EQUAL) {
		return common.IncompleteVarDeclaration{}, fmt.Errorf("expected '=' after variable name")
	}

	value, err := p.expression()
	if err != nil {
		return common.IncompleteVarDeclaration{}, err
	}

	if !p.match(common.SEMICOLON) {
		return common.IncompleteVarDeclaration{}, fmt.Errorf("expected ';' after variable declaration")
	}

	return common.IncompleteVarDeclaration{Name: &name.Value, Value: &value}, nil
}

// EXPRESIONES

func (p *Parser) expression() (common.IncompleteExpression, error) {
	root, err := p.binaryOr()
	if err != nil {
		return common.IncompleteExpression{}, err
	}
	return common.IncompleteExpression{Root: &root}, nil
}

func (p *Parser) binaryOr() (common.IncompleteBinaryOr, error) {

	if p.isAtEnd() {
		return common.IncompleteBinaryOr{}, nil
	}

	if p.match(common.OR) {
		right, err := p.binaryOr()
		if err != nil {
			return common.IncompleteBinaryOr{}, err
		}
		return common.IncompleteBinaryOr{Left: nil, Right: &right}, nil
	}

	left, err := p.binaryAnd()
	if err != nil {
		return common.IncompleteBinaryOr{}, err
	}

	if p.match(common.OR) {
		right, err := p.binaryOr()
		if err != nil {
			return common.IncompleteBinaryOr{}, err
		}

		return common.IncompleteBinaryOr{Left: &left, Right: &right}, nil
	}

	return common.IncompleteBinaryOr{Left: &left, Right: nil}, nil
}

func (p *Parser) binaryAnd() (common.IncompleteBinaryAnd, error) {

	if p.match(common.AND) {
		right, err := p.binaryAnd()
		if err != nil {
			return common.IncompleteBinaryAnd{}, err
		}
		return common.IncompleteBinaryAnd{Left: nil, Right: &right}, nil
	}

	left, err := p.equality()
	if err != nil {
		return common.IncompleteBinaryAnd{}, err
	}

	if p.match(common.AND) {
		right, err := p.binaryAnd()
		if err != nil {
			return common.IncompleteBinaryAnd{}, err
		}

		return common.IncompleteBinaryAnd{Left: &left, Right: &right}, nil
	}

	return common.IncompleteBinaryAnd{Left: &left, Right: nil}, nil
}

func (p *Parser) equality() (common.IncompleteEquality, error) {

	if p.check(common.BANG_EQUAL, common.EQUAL_EQUAL) {
		operator := p.advance()
		right, err := p.equality()
		if err != nil {
			return common.IncompleteEquality{}, err
		}
		return common.IncompleteEquality{Left: nil, Operator: &operator, Right: &right}, nil
	}

	left, err := p.comparison()
	if err != nil {
		return common.IncompleteEquality{}, err
	}

	if p.check(common.BANG_EQUAL, common.EQUAL_EQUAL) {
		operator := p.advance()
		right, err := p.equality()
		if err != nil {
			return common.IncompleteEquality{}, err
		}

		return common.IncompleteEquality{Left: &left, Operator: &operator, Right: &right}, nil
	}

	return common.IncompleteEquality{Left: &left, Operator: nil, Right: nil}, nil
}

func (p *Parser) comparison() (common.IncompleteComparison, error) {

	if p.check(common.GREATER, common.GREATER_EQUAL, common.LESS, common.LESS_EQUAL) {
		operator := p.advance()
		right, err := p.comparison()
		if err != nil {
			return common.IncompleteComparison{}, err
		}
		return common.IncompleteComparison{Left: nil, Operator: &operator, Right: &right}, nil
	}

	left, err := p.term()
	if err != nil {
		return common.IncompleteComparison{}, err
	}

	if p.check(common.GREATER, common.GREATER_EQUAL, common.LESS, common.LESS_EQUAL) {
		operator := p.advance()
		right, err := p.comparison()
		if err != nil {
			return common.IncompleteComparison{}, err
		}

		return common.IncompleteComparison{Left: &left, Operator: &operator, Right: &right}, nil
	}

	return common.IncompleteComparison{Left: &left, Operator: nil, Right: nil}, nil
}

func (p *Parser) term() (common.IncompleteTerm, error) {

	if p.check(common.MINUS, common.PLUS) {
		operator := p.advance()
		right, err := p.term()
		if err != nil {
			return common.IncompleteTerm{}, err
		}
		return common.IncompleteTerm{Left: nil, Operator: &operator, Right: &right}, nil
	}

	left, err := p.factor()
	if err != nil {
		return common.IncompleteTerm{}, err
	}

	if p.check(common.MINUS, common.PLUS) {
		operator := p.advance()
		if p.isAtEnd() {
			return common.IncompleteTerm{Left: &left, Operator: &operator, Right: nil}, nil
		}

		right, err := p.term()
		if err != nil {
			return common.IncompleteTerm{}, err
		}

		return common.IncompleteTerm{Left: &left, Operator: &operator, Right: &right}, nil
	}
	return common.IncompleteTerm{Left: &left, Operator: nil, Right: nil}, nil
}

func (p *Parser) factor() (common.IncompleteFactor, error) {

	if p.check(common.SLASH, common.ASTERISK) {
		operator := p.advance()
		right, err := p.factor()
		if err != nil {
			return common.IncompleteFactor{}, err
		}
		return common.IncompleteFactor{Left: nil, Operator: &operator, Right: &right}, nil
	}

	left, err := p.unary()
	if err != nil {
		return common.IncompleteFactor{}, err
	}

	if p.check(common.SLASH, common.ASTERISK) {
		operator := p.advance()
		if p.isAtEnd() {
			return common.IncompleteFactor{Left: &left, Operator: &operator, Right: nil}, nil
		}
		right, err := p.factor()
		if err != nil {
			return common.IncompleteFactor{}, err
		}

		return common.IncompleteFactor{Left: &left, Operator: &operator, Right: &right}, nil
	}

	return common.IncompleteFactor{Left: &left, Operator: nil, Right: nil}, nil
}

func (p *Parser) unary() (common.IncompleteUnary, error) {
	if p.check(common.BANG, common.MINUS) {
		operator := p.advance()
		if p.isAtEnd() {
			return common.IncompleteUnaryWithOperator{Operator: &operator, Right: nil}, nil
		}
		right, err := p.unary()
		if err != nil {
			return common.IncompleteUnaryWithOperator{}, err
		}
		return common.IncompleteUnaryWithOperator{Operator: &operator, Right: &right}, nil
	}
	return p.primary()
}

// var = 3
// ;

func (p *Parser) primary() (common.IncompleteUnary, error) {

	if p.isAtEnd() {
		if p.previous().Typ == common.OPEN_PARENTHESIS {
			//return common.IncompletePrimary{Value: common.IncompleteGroupingExpression{hasOpen: true, hasClose: false}}, nil
		}
		return common.IncompletePrimary{}, nil
	}

	if p.check(common.FALSE, common.TRUE, common.NONE, common.NUMBER, common.LITERAL) {
		token := p.advance()
		return common.IncompletePrimary{Value: token}, nil
	}

	if p.match(common.OPEN_PARENTHESIS) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(common.CLOSE_PARENTHESIS) {
			return nil, fmt.Errorf("expected ')' after expression")
		}
		return common.IncompletePrimary{Value: common.IncompleteGroupingExpression{Expression: &expr}}, nil
	}

	if p.check(common.IDENTIFIER) {
		identifier := p.advance()

		if p.match(common.OPEN_PARENTHESIS) {
			var args []common.IncompleteExpression
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

			return common.IncompletePrimary{Value: common.IncompleteCall{Callee: &identifier.Value, Arguments: &args}}, nil
		} else {
			return common.IncompletePrimary{Value: identifier}, nil
		}
	}

	return nil, fmt.Errorf("unexpected token: %v", p.peek().String())
}

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
