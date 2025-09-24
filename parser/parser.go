package parser

import (
	"fmt"
	"slices"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Parser struct {
	tokens  []common.Token
	current int
	debug   bool
}

func NewParser(tokens []common.Token, debug bool) *Parser {
	return &Parser{tokens: tokens, current: 0, debug: debug}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) peek() common.Token {
	return p.tokens[p.current]
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

func (p *Parser) program() (common.Program, error) {
	var statements []common.Statement
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return common.Program{}, err
		}
		statements = append(statements, stmt)
	}
	return common.Program{Statements: statements}, nil
}

func (p *Parser) blockStatement() (common.BlockStatement, error) {

	if !p.match(common.OPEN_BRACKET) {
		return common.BlockStatement{}, fmt.Errorf("expected '{' at the beginning of block")
	}

	var statements []common.Statement
	for !p.isAtEnd() && !p.check(common.CLOSE_BRACKET) {
		stmt, err := p.statement()
		if err != nil {
			return common.BlockStatement{}, err
		}
		statements = append(statements, stmt)
	}

	if !p.match(common.CLOSE_BRACKET) {
		return common.BlockStatement{}, fmt.Errorf("expected '}' at the end of block")
	}

	return common.BlockStatement{Statements: statements}, nil
}

func (p *Parser) statement() (common.Statement, error) {

	if p.isAtEnd() {
		return nil, fmt.Errorf("unexpected end of input")
	}

	if p.debug {
		fmt.Printf("[DEBUG] Parsing statement at token %d: %s\n", p.current, p.peek().String())
	}

	if p.check(common.OPEN_BRACKET) {
		return p.blockStatement()
	}

	if p.checkAll(common.IDENTIFIER, common.EQUAL) {
		return p.assignmentStatement()
	}

	if p.checkAll(common.IDENTIFIER, common.OPEN_PARENTHESIS) {
		return p.functionCallStatement()
	}

	token := p.advance()

	switch token.Typ {
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

	default:
		return nil, fmt.Errorf("unexpected token: %v", token.String())
	}
}

func (p *Parser) printStatement() (common.PrintStatement, error) {
	if !p.match(common.OPEN_PARENTHESIS) {
		return common.PrintStatement{}, fmt.Errorf("expected '(' after 'print'")
	}

	expr, err := p.expression()
	if err != nil {
		return common.PrintStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return common.PrintStatement{}, fmt.Errorf("expected ')' after expression")
	}

	if !p.match(common.SEMICOLON) {
		return common.PrintStatement{}, fmt.Errorf("expected ';' after print statement")
	}

	return common.PrintStatement{Value: expr}, nil
}

func (p *Parser) ifStatement() (common.IfStatement, error) {
	if !p.match(common.OPEN_PARENTHESIS) {
		return common.IfStatement{}, fmt.Errorf("expected '(' after 'if'")
	}

	condition, err := p.expression()
	if err != nil {
		return common.IfStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return common.IfStatement{}, fmt.Errorf("expected ')' after if condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return common.IfStatement{}, err
	}

	ifStatement := common.IfStatement{Condition: condition, Body: body}

	if p.matchs(common.ELSE, common.IF) {
		elseIf, err := p.elseIfStatement()
		if err != nil {
			return common.IfStatement{}, err
		}

		ifStatement.ElseIf = &elseIf
	}

	if p.match(common.ELSE) {
		elseBody, err := p.blockStatement()
		if err != nil {
			return common.IfStatement{}, err
		}
		ifStatement.Else = &common.ElseStatement{Body: elseBody}
	}

	return ifStatement, nil
}

func (p *Parser) elseIfStatement() (common.ElseIfStatement, error) {
	if !p.match(common.OPEN_PARENTHESIS) {
		return common.ElseIfStatement{}, fmt.Errorf("expected '(' after 'else if'")
	}

	condition, err := p.expression()
	if err != nil {
		return common.ElseIfStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return common.ElseIfStatement{}, fmt.Errorf("expected ')' after else if condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return common.ElseIfStatement{}, err
	}

	if p.matchs(common.ELSE, common.IF) {
		elseIf, err := p.elseIfStatement()
		if err != nil {
			return common.ElseIfStatement{}, err
		}
		return common.ElseIfStatement{Condition: condition, Body: body, ElseIf: &elseIf}, nil
	}

	return common.ElseIfStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) breakStatement() (common.BreakStatement, error) {
	if !p.match(common.SEMICOLON) {
		return common.BreakStatement{}, fmt.Errorf("expected ';' after 'break'")
	}
	return common.BreakStatement{}, nil
}

func (p *Parser) returnStatement() (common.ReturnStatement, error) {

	expr, err := p.expression()
	if err != nil {
		return common.ReturnStatement{}, err
	}
	if !p.match(common.SEMICOLON) {
		return common.ReturnStatement{}, fmt.Errorf("expected ';' after 'return'")
	}
	return common.ReturnStatement{Value: expr}, nil
}

func (p *Parser) funcStatement() (common.FunctionDef, error) {
	if !p.check(common.IDENTIFIER) {
		return common.FunctionDef{}, fmt.Errorf("expected function name after 'func'")
	}
	name := p.advance()

	if !p.match(common.OPEN_PARENTHESIS) {
		return common.FunctionDef{}, fmt.Errorf("expected '(' after function name")
	}

	var parameters []string

	if !p.match(common.CLOSE_PARENTHESIS) {
		for {
			if !p.check(common.IDENTIFIER) {
				return common.FunctionDef{}, fmt.Errorf("expected parameter name")
			}
			parameters = append(parameters, p.advance().Value)

			if p.match(common.CLOSE_PARENTHESIS) {
				break
			}

			if !p.match(common.COMMA) {
				return common.FunctionDef{}, fmt.Errorf("expected ',' or ')' after parameter")
			}
		}
	}

	body, err := p.blockStatement()
	if err != nil {
		return common.FunctionDef{}, err
	}

	return common.FunctionDef{Name: name.Value, Parameters: parameters, Body: body}, nil
}

func (p *Parser) whileStatement() (common.WhileStatement, error) {
	if !p.match(common.OPEN_PARENTHESIS) {
		return common.WhileStatement{}, fmt.Errorf("expected '(' after 'while'")
	}

	condition, err := p.expression()
	if err != nil {
		return common.WhileStatement{}, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return common.WhileStatement{}, fmt.Errorf("expected ')' after while condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return common.WhileStatement{}, err
	}

	return common.WhileStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) assignmentStatement() (common.Assignment, error) {

	if !p.check(common.IDENTIFIER) {
		return common.Assignment{}, fmt.Errorf("expected variable name")
	}

	name := p.advance().Value

	if !p.match(common.EQUAL) {
		return common.Assignment{}, fmt.Errorf("expected '=' after variable name")
	}

	value, err := p.expression()
	if err != nil {
		return common.Assignment{}, err
	}

	if !p.match(common.SEMICOLON) {
		return common.Assignment{}, fmt.Errorf("expected ';' after assignment")
	}

	return common.Assignment{Name: name, Value: value}, nil
}

func (p *Parser) functionCallStatement() (common.CallStatement, error) {
	if !p.check(common.IDENTIFIER) {
		return common.CallStatement{}, fmt.Errorf("expected function name")
	}

	name := p.advance().Value

	if !p.match(common.OPEN_PARENTHESIS) {
		return common.CallStatement{}, fmt.Errorf("expected '(' after function name")
	}

	var args []common.Expression
	if !p.match(common.CLOSE_PARENTHESIS) {
		for {
			arg, err := p.expression()
			if err != nil {
				return common.CallStatement{}, err
			}
			args = append(args, arg)

			if p.match(common.CLOSE_PARENTHESIS) {
				break
			}

			if !p.match(common.COMMA) {
				return common.CallStatement{}, fmt.Errorf("expected ',' or ')' after argument")
			}
		}
	}

	if !p.match(common.SEMICOLON) {
		return common.CallStatement{}, fmt.Errorf("expected ';' after function call")
	}

	return common.CallStatement{Call: common.Call{Callee: name, Arguments: args}}, nil
}

func (p *Parser) varStatement() (common.VarDeclaration, error) {
	if !p.check(common.IDENTIFIER) {
		return common.VarDeclaration{}, fmt.Errorf("expected variable name after '%s'", common.VAR_KEYWORD)
	}

	name := p.advance()

	if !p.match(common.EQUAL) {
		return common.VarDeclaration{}, fmt.Errorf("expected '=' after variable name")
	}

	value, err := p.expression()
	if err != nil {
		return common.VarDeclaration{}, err
	}

	if !p.match(common.SEMICOLON) {
		return common.VarDeclaration{}, fmt.Errorf("expected ';' after variable declaration")
	}

	return common.VarDeclaration{Name: name.Value, Value: value}, nil
}

// EXPRESIONES

func (p *Parser) expression() (common.Expression, error) {
	root, err := p.binaryOr()
	if err != nil {
		return common.Expression{}, err
	}
	return common.Expression{Root: root}, nil
}

func (p *Parser) binaryOr() (common.BinaryOr, error) {
	left, err := p.binaryAnd()
	if err != nil {
		return common.BinaryOr{}, err
	}

	if p.match(common.OR) {
		right, err := p.binaryOr()
		if err != nil {
			return common.BinaryOr{}, err
		}

		return common.BinaryOr{Left: left, Right: &right}, nil
	}

	return common.BinaryOr{Left: left, Right: nil}, nil
}

func (p *Parser) binaryAnd() (common.BinaryAnd, error) {
	left, err := p.equality()
	if err != nil {
		return common.BinaryAnd{}, err
	}

	if p.match(common.AND) {
		right, err := p.binaryAnd()
		if err != nil {
			return common.BinaryAnd{}, err
		}

		return common.BinaryAnd{Left: left, Right: &right}, nil
	}

	return common.BinaryAnd{Left: left, Right: nil}, nil
}

func (p *Parser) equality() (common.Equality, error) {
	left, err := p.comparison()
	if err != nil {
		return common.Equality{}, err
	}

	if p.check(common.BANG_EQUAL, common.EQUAL_EQUAL) {
		operator := p.advance()
		right, err := p.equality()
		if err != nil {
			return common.Equality{}, err
		}

		return common.Equality{Left: left, Operator: &operator, Right: &right}, nil
	}

	return common.Equality{Left: left, Operator: nil, Right: nil}, nil
}

func (p *Parser) comparison() (common.Comparison, error) {
	left, err := p.term()
	if err != nil {
		return common.Comparison{}, err
	}

	if p.check(common.GREATER, common.GREATER_EQUAL, common.LESS, common.LESS_EQUAL) {
		operator := p.advance()
		right, err := p.comparison()
		if err != nil {
			return common.Comparison{}, err
		}

		return common.Comparison{Left: left, Operator: &operator, Right: &right}, nil
	}

	return common.Comparison{Left: left, Operator: nil, Right: nil}, nil
}

func (p *Parser) term() (common.Term, error) {
	left, err := p.factor()
	if err != nil {
		return common.Term{}, err
	}

	if p.check(common.MINUS, common.PLUS) {
		operator := p.advance()
		right, err := p.term()
		if err != nil {
			return common.Term{}, err
		}

		return common.Term{Left: left, Operator: &operator, Right: &right}, nil
	}
	return common.Term{Left: left, Operator: nil, Right: nil}, nil
}

func (p *Parser) factor() (common.Factor, error) {
	left, err := p.unary()
	if err != nil {
		return common.Factor{}, err
	}

	if p.check(common.SLASH, common.ASTERISK) {
		operator := p.advance()
		right, err := p.factor()
		if err != nil {
			return common.Factor{}, err
		}

		return common.Factor{Left: left, Operator: &operator, Right: &right}, nil
	}

	return common.Factor{Left: left, Operator: nil, Right: nil}, nil
}

func (p *Parser) unary() (common.Unary, error) {
	if p.check(common.BANG, common.MINUS) {
		operator := p.advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return common.UnaryWithOperator{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (common.Unary, error) {
	if p.isAtEnd() {
		return common.Primary{}, fmt.Errorf("unexpected end of input")
	}

	if p.check(common.FALSE, common.TRUE, common.NONE, common.NUMBER, common.LITERAL) {
		token := p.advance()
		return common.Primary{Value: token}, nil
	}

	if p.match(common.OPEN_PARENTHESIS) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(common.CLOSE_PARENTHESIS) {
			return nil, fmt.Errorf("expected ')' after expression")
		}
		return common.Primary{Value: common.GroupingExpression{Expression: expr}}, nil
	}

	if p.check(common.IDENTIFIER) {
		identifier := p.advance()

		if p.match(common.OPEN_PARENTHESIS) {
			var args []common.Expression
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

			return common.Primary{Value: common.Call{Callee: identifier.Value, Arguments: args}}, nil
		} else {
			return common.Primary{Value: identifier}, nil
		}
	}

	return nil, fmt.Errorf("unexpected token: %v", p.peek().String())
}

func (p *Parser) Parse() (common.Program, error) {
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
