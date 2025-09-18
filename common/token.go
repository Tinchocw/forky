package common

import "fmt"

type Token struct {
	Typ   TokenType
	Value string
}

func NewToken(typ TokenType) Token                        { return Token{Typ: typ, Value: ""} }
func NewTokenWithValue(typ TokenType, value string) Token { return Token{Typ: typ, Value: value} }

// valueCarrier indicates whether this token type should print and preserve its value.
func (tt TokenType) valueCarrier() bool {
	switch tt {
	case NUMBER, LITERAL, IDENTIFIER, STARTED_LITERAL, ENDED_LITERAL:
		return true
	default:
		return false
	}
}

// String implements fmt.Stringer for Token to allow pretty printing.
// Format: TOKEN_TYPE(value). For literals we keep the original delimiters (e.g. "text').
func (t Token) String() string {
	if t.Typ.valueCarrier() {
		return fmt.Sprintf("%s(\"%s\")", t.Typ.String(), t.Value)
	}
	return t.Typ.String()
}

type TokenType int

var tokenTypeNames = [...]string{
	PLUS:              "PLUS",
	MINUS:             "MINUS",
	ASTERISK:          "ASTERISK",
	SLASH:             "SLASH",
	EQUAL:             "EQUAL",
	BANG:              "BANG",
	LESS:              "LESS",
	GREATER:           "GREATER",
	OPEN_PARENTHESIS:  "OPEN_PARENTHESIS",
	CLOSE_PARENTHESIS: "CLOSE_PARENTHESIS",
	COMMA:             "COMMA",
	COLON:             "COLON",
	SEMICOLON:         "SEMICOLON",
	EQUAL_EQUAL:       "EQUAL_EQUAL",
	BANG_EQUAL:        "BANG_EQUAL",
	LESS_EQUAL:        "LESS_EQUAL",
	GREATER_EQUAL:     "GREATER_EQUAL",
	NUMBER:            "NUMBER",
	LITERAL:           "LITERAL",
	TRUE:              "TRUE",
	FALSE:             "FALSE",
	NONE:              "NONE",
	IF:                "IF",
	ELSE:              "ELSE",
	WHILE:             "WHILE",
	RETURN:            "RETURN",
	CONTINUE:          "CONTINUE",
	BREAK:             "BREAK",
	IDENTIFIER:        "IDENTIFIER",
	FUNC:              "FUNC",
	VAR:               "VAR",
	STARTED_LITERAL:   "STARTED_LITERAL",
	ENDED_LITERAL:     "ENDED_LITERAL",
}

func (t TokenType) String() string {
	if int(t) >= 0 && int(t) < len(tokenTypeNames) && tokenTypeNames[t] != "" {
		return tokenTypeNames[t]
	}
	return "UNKNOWN"
}

const (
	// SINGLE CHARACTER TOKENS
	PLUS TokenType = iota
	MINUS
	ASTERISK
	SLASH
	EQUAL
	BANG
	LESS
	GREATER
	OPEN_PARENTHESIS
	CLOSE_PARENTHESIS
	COMMA
	COLON
	SEMICOLON

	// MULTI CHARACTER TOKENS
	EQUAL_EQUAL
	BANG_EQUAL
	LESS_EQUAL
	GREATER_EQUAL

	// LITERALS
	NUMBER
	LITERAL
	TRUE
	FALSE
	NONE

	// CONTROL FLOW
	IF
	ELSE
	WHILE
	RETURN
	CONTINUE
	BREAK

	// IDENTIFIERS
	IDENTIFIER
	FUNC
	VAR

	// SPECIAL TOKENS PRE MERGE
	STARTED_LITERAL
	ENDED_LITERAL
)

// SYMBOLS
const (
	PLUS_SYMBOL              = '+'
	MINUS_SYMBOL             = '-'
	ASTERISK_SYMBOL          = '*'
	SLASH_SYMBOL             = '/'
	EQUAL_SYMBOL             = '='
	BANG_SYMBOL              = '!'
	LESS_SYMBOL              = '<'
	GREATER_SYMBOL           = '>'
	START_QUOTE_SYMBOL       = '"'
	END_QUOTE_SYMBOL         = '\''
	OPEN_PARENTHESIS_SYMBOL  = '('
	CLOSE_PARENTHESIS_SYMBOL = ')'
	COMMA_SYMBOL             = ','
	COLON_SYMBOL             = ':'
	SEMICOLON_SYMBOL         = ';'
)

// Keywords
const (
	TRUE_KEYWORD     = "true"
	FALSE_KEYWORD    = "false"
	NONE_KEYWORD     = "none"
	IF_KEYWORD       = "if"
	ELSE_KEYWORD     = "else"
	WHILE_KEYWORD    = "while"
	FUNC_KEYWORD     = "func"
	RETURN_KEYWORD   = "return"
	VAR_KEYWORD      = "var"
	CONTINUE_KEYWORD = "continue"
	BREAK_KEYWORD    = "break"
)

var KEYWORDS = map[string]TokenType{
	TRUE_KEYWORD:   TRUE,
	FALSE_KEYWORD:  FALSE,
	NONE_KEYWORD:   NONE,
	IF_KEYWORD:     IF,
	ELSE_KEYWORD:   ELSE,
	WHILE_KEYWORD:  WHILE,
	FUNC_KEYWORD:   FUNC,
	RETURN_KEYWORD: RETURN,
	VAR_KEYWORD:    VAR,
}

var KEYWORDS_VALUES = map[TokenType]string{
	TRUE:   TRUE_KEYWORD,
	FALSE:  FALSE_KEYWORD,
	NONE:   NONE_KEYWORD,
	IF:     IF_KEYWORD,
	ELSE:   ELSE_KEYWORD,
	WHILE:  WHILE_KEYWORD,
	FUNC:   FUNC_KEYWORD,
	RETURN: RETURN_KEYWORD,
	VAR:    VAR_KEYWORD,
}

func IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func IsAlphanumeric(r rune) bool {
	return IsLetter(r) || IsNumber(r)
}

func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}
