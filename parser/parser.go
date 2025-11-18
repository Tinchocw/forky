package parser

import (
	"fmt"
	"slices"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
	"github.com/Tinchocw/forky/common/statement"
	"github.com/Tinchocw/forky/common/statement/assignment"
	"github.com/Tinchocw/forky/common/statement/block"
	"github.com/Tinchocw/forky/common/statement/declaration"
	"github.com/Tinchocw/forky/common/statement/extra"
	"github.com/Tinchocw/forky/common/statement/flow"
	"github.com/Tinchocw/forky/common/statement/function"
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
	statements := []statement.Statement{}
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return statement.Program{}, err
		}
		statements = append(statements, stmt)
	}

	return statement.Program{Statements: statements}, nil
}

func (p *Parser) blockStatement() (*block.BlockStatement, error) {

	if !p.match(common.OPEN_BRACES) {
		return nil, fmt.Errorf("expected '{' at the beginning of block")
	}

	statements := []statement.Statement{}
	for !p.isAtEnd() && !p.check(common.CLOSE_BRACES) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	if !p.match(common.CLOSE_BRACES) {
		return nil, fmt.Errorf("expected '}' at the end of block")
	}

	return &block.BlockStatement{Statements: statements}, nil
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
	case common.FORK:
		return p.forkStatement()
	case common.IF:
		return p.ifStatement()
	case common.BREAK:
		return p.breakStatement()
	case common.FUNC:
		return p.funcStatement()
	case common.VAR:
		return p.declarationStatement()
	case common.SET:
		return p.assignmentStatement()
	case common.WHILE:
		return p.whileStatement()
	case common.OPEN_BRACES:
		return p.blockStatement()
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) printStatement() (*extra.PrintStatement, error) {
	if !p.match(common.PRINT) {
		return nil, fmt.Errorf("expected 'print'")
	}

	if !p.match(common.OPEN_PARENTHESIS) {
		return nil, fmt.Errorf("expected '(' after 'print'")
	}

	if p.match(common.CLOSE_PARENTHESIS) {
		if !p.match(common.SEMICOLON) {
			return nil, fmt.Errorf("expected ';' after print statement")
		}

		return &extra.PrintStatement{}, nil
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return nil, fmt.Errorf("expected ')' after expression")
	}

	if !p.match(common.SEMICOLON) {
		return nil, fmt.Errorf("expected ';' after print statement")
	}

	return &extra.PrintStatement{Value: expr}, nil
}

func (p *Parser) forkStatement() (extra.ForkStatement, error) {
	if !p.match(common.FORK) {
		return nil, fmt.Errorf("expected 'fork'")
	}

	if p.check(common.OPEN_BRACES) {
		return p.forkBlockStatement()
	} else {
		return p.forkArrayStatement()
	}
}

func (p *Parser) forkBlockStatement() (*extra.ForkBlockStatement, error) {
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &extra.ForkBlockStatement{Block: body}, nil
}

func (p *Parser) forkArrayStatement() (*extra.ForkArrayStatement, error) {
	array, err := p.expression()
	if err != nil {
		return nil, err
	}

	// fork array index,elem {

	var elemName *string
	var indexName *string

	if p.check(common.IDENTIFIER) {
		firstToken := p.advance()
		elemName = &firstToken.Value

		if p.match(common.COMMA) {
			if !p.check(common.IDENTIFIER) {
				return nil, fmt.Errorf("expected identifier after ',' in fork array statement")
			}
			secondToken := p.advance()

			indexName = &firstToken.Value
			elemName = &secondToken.Value
		}
	}

	block, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &extra.ForkArrayStatement{Array: array, ElemName: elemName, IndexName: indexName, Block: block}, nil
}

func (p *Parser) ifStatement() (*flow.IfStatement, error) {
	if !p.match(common.IF) {
		return nil, fmt.Errorf("expected 'if'")
	}

	if !p.match(common.OPEN_PARENTHESIS) {
		return nil, fmt.Errorf("expected '(' after 'if'")
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return nil, fmt.Errorf("expected ')' after if condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	ifStatement := &flow.IfStatement{Condition: condition, Body: body}

	if p.matchs(common.ELSE, common.IF) {
		elseIf, err := p.elseIfStatement()
		if err != nil {
			return nil, err
		}

		ifStatement.ElseIf = elseIf
	}

	if p.match(common.ELSE) {
		elseBody, err := p.blockStatement()
		if err != nil {
			return nil, err
		}
		ifStatement.Else = &flow.ElseStatement{Body: elseBody}
	}

	return ifStatement, nil
}

func (p *Parser) elseIfStatement() (*flow.ElseIfStatement, error) {
	if !p.match(common.OPEN_PARENTHESIS) {
		return nil, fmt.Errorf("expected '(' after 'else if'")
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return nil, fmt.Errorf("expected ')' after else if condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	if p.matchs(common.ELSE, common.IF) {
		elseIf, err := p.elseIfStatement()
		if err != nil {
			return nil, err
		}
		return &flow.ElseIfStatement{Condition: condition, Body: body, ElseIf: elseIf}, nil
	}

	return &flow.ElseIfStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) breakStatement() (*flow.BreakStatement, error) {
	if !p.match(common.BREAK) {
		return nil, fmt.Errorf("expected 'break'")
	}
	if !p.match(common.SEMICOLON) {
		return nil, fmt.Errorf("expected ';' after 'break'")
	}
	return &flow.BreakStatement{}, nil
}

func (p *Parser) returnStatement() (*function.ReturnStatement, error) {
	if !p.match(common.RETURN) {
		return nil, fmt.Errorf("expected 'return'")
	}

	if p.match(common.SEMICOLON) {
		return &function.ReturnStatement{}, nil
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.match(common.SEMICOLON) {
		return nil, fmt.Errorf("expected ';' after 'return'")
	}
	return &function.ReturnStatement{Value: expr}, nil
}

func (p *Parser) funcStatement() (*function.FunctionDef, error) {
	if !p.match(common.FUNC) {
		return nil, fmt.Errorf("expected 'func'")
	}

	if !p.check(common.IDENTIFIER) {
		return nil, fmt.Errorf("expected function name after 'func'")
	}
	name := p.advance()

	if !p.match(common.OPEN_PARENTHESIS) {
		return nil, fmt.Errorf("expected '(' after function name")
	}

	parameters := []string{}

	if !p.match(common.CLOSE_PARENTHESIS) {
		for {
			if !p.check(common.IDENTIFIER) {
				return nil, fmt.Errorf("expected parameter name")
			}
			parameters = append(parameters, p.advance().Value)

			if p.match(common.CLOSE_PARENTHESIS) {
				break
			}

			if !p.match(common.COMMA) {
				return nil, fmt.Errorf("expected ',' or ')' after parameter")
			}
		}
	}

	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &function.FunctionDef{Name: &name.Value, Parameters: parameters, Body: body}, nil
}

func (p *Parser) whileStatement() (*flow.WhileStatement, error) {
	if !p.match(common.WHILE) {
		return nil, fmt.Errorf("expected 'while' at the beginning of while statement")
	}

	if !p.match(common.OPEN_PARENTHESIS) {
		return nil, fmt.Errorf("expected '(' after 'while'")
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(common.CLOSE_PARENTHESIS) {
		return nil, fmt.Errorf("expected ')' after while condition")
	}

	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &flow.WhileStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) assignmentStatement() (assignment.Assignment, error) {
	if !p.match(common.SET) {
		return nil, fmt.Errorf("expected 'set' at the beginning of assignment")
	}

	if !p.check(common.IDENTIFIER) {
		return nil, fmt.Errorf("expected variable name")
	}

	name := p.advance()
	var assigment assignment.Assignment
	var err error

	if p.check(common.OPEN_BRACKET) {
		assigment, err = p.arrayAssignmentStatement(name)
	} else {
		assigment, err = p.varAssigmentStatement(name)
	}

	return assigment, err
}

func (p *Parser) varAssigmentStatement(name common.Token) (assignment.Assignment, error) {
	if !p.match(common.EQUAL) {
		return nil, fmt.Errorf("expected '=' after variable name")
	}

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(common.SEMICOLON) {
		return nil, fmt.Errorf("expected ';' after assignment")
	}

	return &assignment.VarAssignment{Name: name.Value, Value: value}, nil
}

func (p *Parser) arrayAssignmentStatement(name common.Token) (assignment.Assignment, error) {
	indexes := []*expression.ExpressionNode{}

	for p.match(common.OPEN_BRACKET) {
		index, err := p.expression()
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, index)
		if !p.match(common.CLOSE_BRACKET) {
			return nil, fmt.Errorf("expected ']' after index expression")
		}
	}

	if !p.match(common.EQUAL) {
		return nil, fmt.Errorf("expected '=' after array name and indexes")
	}

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(common.SEMICOLON) {
		return nil, fmt.Errorf("expected ';' after assignment")
	}

	return &assignment.ArrayAssignment{Name: name.Value, Indexes: indexes, Value: value}, nil
}

func (p *Parser) expressionStatement() (*statement.ExpressionStatement, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(common.SEMICOLON) {
		return nil, fmt.Errorf("expected ';' after expression")
	}

	return &statement.ExpressionStatement{Expression: expr}, nil
}

func (p *Parser) declarationStatement() (declaration.DeclarationStatement, error) {
	if !p.match(common.VAR) {
		return nil, fmt.Errorf("expected 'var' at the beginning of a declaration")
	}

	if !p.check(common.IDENTIFIER) {
		return nil, fmt.Errorf("expected variable name after '%s'", common.VAR_KEYWORD)
	}

	name := p.advance()
	var declaration declaration.DeclarationStatement
	var err error

	if p.check(common.OPEN_BRACKET) {
		declaration, err = p.arrayDeclarationStatement(name)
	} else {
		declaration, err = p.varDeclarationStatement(name)
	}

	if err != nil {
		return nil, err
	}

	return declaration, nil
}

func (p *Parser) varDeclarationStatement(name common.Token) (*declaration.VarDeclaration, error) {
	if !p.match(common.EQUAL) {
		if !p.match(common.SEMICOLON) {
			return nil, fmt.Errorf("expected '=' or ';' after variable name")
		}
		return &declaration.VarDeclaration{Name: name.Value}, nil

	}

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(common.SEMICOLON) {
		return nil, fmt.Errorf("expected ';' after variable declaration")
	}

	return &declaration.VarDeclaration{Name: name.Value, Value: value}, nil
}

func (p *Parser) arrayDeclarationStatement(name common.Token) (*declaration.ArrayDeclaration, error) {
	lengths := []*expression.ExpressionNode{}

	for p.match(common.OPEN_BRACKET) {
		length, err := p.expression()
		if err != nil {
			return nil, err
		}
		lengths = append(lengths, length)

		if !p.match(common.CLOSE_BRACKET) {
			return nil, fmt.Errorf("expected ']' after size expression")
		}
	}

	if !p.match(common.EQUAL) {
		if !p.match(common.SEMICOLON) {
			return nil, fmt.Errorf("expected '=' or ';' after variable name")
		}

		return &declaration.ArrayDeclaration{Name: name.Value, Lengths: lengths}, nil
	}

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(common.SEMICOLON) {
		return nil, fmt.Errorf("expected ';' after variable declaration")
	}

	return &declaration.ArrayDeclaration{Name: name.Value, Lengths: lengths, Value: value}, nil
}

// EXPRESIONES

func (p *Parser) expression() (*expression.ExpressionNode, error) {
	root, err := p.logicalOr()
	if err != nil {
		return nil, err
	}
	return &expression.ExpressionNode{Root: root}, nil
}

func (p *Parser) logicalOr() (*expression.LogicalOrNode, error) {
	left, err := p.logicalAnd()
	if err != nil {
		return nil, err
	}

	bor := &expression.LogicalOrNode{Left: left}
	first := true

	for p.check(common.OR) {
		if !first {
			bor = &expression.LogicalOrNode{Left: bor}
		}

		operator := p.advance()
		right, err := p.logicalAnd()
		if err != nil {
			return nil, err
		}
		bor.Operator = &operator
		bor.Right = right

		first = false
	}

	return bor, nil
}

func (p *Parser) logicalAnd() (*expression.LogicalAndNode, error) {
	left, err := p.equality()
	if err != nil {
		return nil, err
	}

	band := &expression.LogicalAndNode{Left: left}
	first := true

	for p.check(common.AND) {
		if !first {
			band = &expression.LogicalAndNode{Left: band}
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

func (p *Parser) equality() (*expression.EqualityNode, error) {
	left, err := p.comparison()
	if err != nil {
		return nil, err
	}

	eq := &expression.EqualityNode{Left: left}
	first := true

	for p.check(common.BANG_EQUAL, common.EQUAL_EQUAL) {
		if !first {
			eq = &expression.EqualityNode{Left: eq}
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

func (p *Parser) comparison() (*expression.ComparisonNode, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}

	comp := &expression.ComparisonNode{Left: left}
	first := true

	for p.check(common.GREATER, common.GREATER_EQUAL, common.LESS, common.LESS_EQUAL) {
		if !first {
			comp = &expression.ComparisonNode{Left: comp}
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

func (p *Parser) term() (*expression.TermNode, error) {
	left, err := p.factor()
	if err != nil {
		return nil, err
	}

	t := &expression.TermNode{Left: left}
	first := true

	for p.check(common.MINUS, common.PLUS) {
		if !first {
			t = &expression.TermNode{Left: t}
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

func (p *Parser) factor() (*expression.FactorNode, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}

	f := &expression.FactorNode{Left: left}
	first := true

	for p.check(common.SLASH, common.ASTERISK) {
		if !first {
			f = &expression.FactorNode{Left: f}
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

func (p *Parser) unary() (*expression.UnaryNode, error) {
	u := &expression.UnaryNode{}

	if p.check(common.BANG, common.TILDE) {
		operator := p.advance()
		u.Operator = &operator

		if p.check(common.BANG, common.TILDE) {
			rigth, err := p.unary()
			if err != nil {
				return nil, err
			}
			u.Right = rigth
			return u, nil
		}
	}

	right, err := p.arrayAccess()
	if err != nil {
		return nil, err
	}
	u.Right = right
	return u, nil
}

func (p *Parser) arrayAccess() (*expression.ArrayAccessNode, error) {
	left, err := p.functionCall()
	if err != nil {
		return nil, err
	}

	aa := &expression.ArrayAccessNode{Left: left}
	first := true

	for p.match(common.OPEN_BRACKET) {
		if !first {
			aa = &expression.ArrayAccessNode{Left: aa}
		}

		indexExpr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(common.CLOSE_BRACKET) {
			return nil, fmt.Errorf("expected ']' after index expression")
		}

		aa.Index = indexExpr
		first = false
	}

	return aa, nil
}

func (p *Parser) functionCall() (*expression.FunctionCallNode, error) {
	left, err := p.primary()
	if err != nil {
		return nil, err
	}

	fc := &expression.FunctionCallNode{Callee: left}
	first := true

	for p.match(common.OPEN_PARENTHESIS) {
		if !first {
			fc = &expression.FunctionCallNode{Callee: fc}
		}
		args := []*expression.ExpressionNode{}

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

		fc.Arguments = args
		first = false
	}

	return fc, nil
}

func (p *Parser) primary() (expression.Primary, error) {
	if p.isAtEnd() {
		return nil, fmt.Errorf("unexpected end of input")
	}

	if p.check(common.FALSE, common.TRUE, common.NONE, common.NUMBER, common.LITERAL) {
		token := p.advance()
		return &expression.TokenLiteralNode{Token: &token}, nil
	}

	if p.match(common.OPEN_PARENTHESIS) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(common.CLOSE_PARENTHESIS) {
			return nil, fmt.Errorf("expected ')' after expression")
		}

		return &expression.GroupingExpressionNode{Expression: expr}, nil
	}

	if p.match(common.OPEN_BRACKET) {
		elements := []*expression.ExpressionNode{}
		for !p.match(common.CLOSE_BRACKET) && !p.isAtEnd() {
			element, err := p.expression()
			if err != nil {
				return nil, err
			}
			elements = append(elements, element)
			if !p.match(common.COMMA) {
				if p.match(common.CLOSE_BRACKET) {
					break
				}

				return nil, fmt.Errorf("expected ',' or ']' after array element")
			}
		}

		return &expression.ArrayLiteralNode{Elements: elements}, nil
	}

	if p.check(common.IDENTIFIER) {
		identifier := p.advance()
		return &expression.TokenLiteralNode{Token: &identifier}, nil
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
