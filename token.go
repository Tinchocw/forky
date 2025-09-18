package main

import "fmt"

type Token struct {
	typ   TokenType
	value string
}

func NewToken(typ TokenType) Token {
	return Token{typ: typ, value: ""}
}

func NewTokenWithValue(typ TokenType, value string) Token {
	return Token{typ: typ, value: value}
}

// String implements fmt.Stringer for Token to allow pretty printing.
// Format: TOKEN_TYPE(value). For literals we keep the original delimiters (e.g. "text').
func (t Token) String() string {
	switch t.typ {
	case NUMBER, LITERAL, IDENTIFIER, STARTED_LITERAL, ENDED_LITERAL:
		return fmt.Sprintf("%s(\"%s\")", t.typ.String(), t.value)
	default:
		// Keywords TRUE, FALSE, NONE and operators/comparators: just the type name
		return t.typ.String()
	}
}

type TokenType int

func (t TokenType) String() string {
	switch t {
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case EQUAL:
		return "EQUAL"
	case MARK:
		return "MARK"
	case LESS:
		return "LESS"
	case GREATER:
		return "GREATER"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case MARK_EQUAL:
		return "MARK_EQUAL"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case NUMBER:
		return "NUMBER"
	case LITERAL:
		return "LITERAL"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STARTED_LITERAL:
		return "STARTED_LITERAL"
	case ENDED_LITERAL:
		return "ENDED_LITERAL"
	default:
		return "UNKNOWN"
	}
}

const (
	// SINGLE CHARACTER TOKENS
	PLUS TokenType = iota
	MINUS
	ASTERISK
	SLASH
	EQUAL
	MARK
	LESS
	GREATER

	// MULTI CHARACTER TOKENS
	EQUAL_EQUAL
	MARK_EQUAL
	LESS_EQUAL
	GREATER_EQUAL

	// LITERALS
	NUMBER
	LITERAL
	TRUE
	FALSE
	NONE

	// OTHERS
	IDENTIFIER

	// SPECIAL TOKENS PRE MERGE
	STARTED_LITERAL
	ENDED_LITERAL
)

// SYMBOLS
const (
	PLUS_SYMBOL     = '+'
	MINUS_SYMBOL    = '-'
	ASTERISK_SYMBOL = '*'
	SLASH_SYMBOL    = '/'
	EQUAL_SYMBOL    = '='
	MARK_SYMBOL     = '!'
	LESS_SYMBOL     = '<'
	GREATER_SYMBOL  = '>'
	START_QUOTE     = '"'
	END_QUOTE       = '\''
)

// Keywords
const (
	TRUE_KEYWORD  = "true"
	FALSE_KEYWORD = "false"
	NONE_KEYWORD  = "none"
)

var KEYWORDS = map[string]TokenType{
	TRUE_KEYWORD:  TRUE,
	FALSE_KEYWORD: FALSE,
	NONE_KEYWORD:  NONE,
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isAlphanumeric(r rune) bool {
	return isLetter(r) || isNumber(r)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}
