package scanner

import (
	"fmt"
	"unicode/utf8"

	"github.com/Tinchocw/forky/common"
)

type scanner struct {
	content       string
	tokens        []common.Token
	index         int
	canMergeStart bool
	canMergeEnd   bool
}

func createScanner(content string) scanner {
	return scanner{
		content:       content,
		tokens:        []common.Token{},
		index:         0,
		canMergeStart: true,
		canMergeEnd:   true,
	}
}

func (s *scanner) isAtStart() bool {
	return s.index == 1 // Is at start if index is 1, as we have already advanced once
}

func (s *scanner) isAtEnd() bool {
	return s.index >= len(s.content)
}

func (s *scanner) lookahead() rune {
	if s.isAtEnd() {
		panic("No more characters to lookahead")
	}
	r, _ := utf8.DecodeRuneInString(s.content[s.index:])
	return r
}

func (s *scanner) advance() rune {
	r, size := utf8.DecodeRuneInString(s.content[s.index:])
	s.index += size
	return r
}

// If the next characters match the expected string, consume them and return true.
// Otherwise, return false.
func (s *scanner) match(expected string) bool {
	if len(s.content)-s.index < len(expected) {
		return false
	}

	if s.content[s.index:s.index+len(expected)] == expected {
		s.index += len(expected)
		return true
	}

	return false
}

func (s *scanner) matchRune(expected rune) bool {
	return s.match(string(expected))
}

func (s *scanner) addToken(typ common.TokenType) { s.tokens = append(s.tokens, common.NewToken(typ)) }

func (s *scanner) addTokenWithValue(typ common.TokenType, value string) {
	s.tokens = append(s.tokens, common.NewTokenWithValue(typ, value))
}

func (s *scanner) clearTokens() {
	s.tokens = []common.Token{}
}

func (s *scanner) consumeWhile(condition func(rune) bool) {
	for !s.isAtEnd() && condition(s.lookahead()) {
		s.advance()
	}
}

func isAllDigits(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !(r >= '0' && r <= '9') { // inline digit check
			return false
		}
	}
	return true
}

func (s *scanner) scan() (segment, error) {
	for !s.isAtEnd() {
		r := s.advance()

		if common.IsWhitespace(r) {
			if s.isAtStart() {
				s.canMergeStart = false
			}

			if s.isAtEnd() {
				s.canMergeEnd = false
			}

			continue
		}

		switch r {
		case common.PLUS_SYMBOL:
			s.addToken(common.PLUS)

		case common.MINUS_SYMBOL:
			s.addToken(common.MINUS)

		case common.ASTERISK_SYMBOL:
			s.addToken(common.ASTERISK)

		case common.SLASH_SYMBOL:
			s.addToken(common.SLASH)

		case common.COMMA_SYMBOL:
			s.addToken(common.COMMA)

		case common.COLON_SYMBOL:
			s.addToken(common.COLON)

		case common.SEMICOLON_SYMBOL:
			s.addToken(common.SEMICOLON)

		case common.OPEN_PARENTHESIS_SYMBOL:
			s.addToken(common.OPEN_PARENTHESIS)

		case common.CLOSE_PARENTHESIS_SYMBOL:
			s.addToken(common.CLOSE_PARENTHESIS)

		case common.OPEN_BRACES_SYMBOL:
			s.addToken(common.OPEN_BRACES)

		case common.CLOSE_BRACES_SYMBOL:
			s.addToken(common.CLOSE_BRACES)

		case common.OPEN_BRACKET_SYMBOL:
			s.addToken(common.OPEN_BRACKET)

		case common.CLOSE_BRACKET_SYMBOL:
			s.addToken(common.CLOSE_BRACKET)

		case common.EQUAL_SYMBOL:
			if s.matchRune(common.EQUAL_SYMBOL) {
				s.addToken(common.EQUAL_EQUAL)
			} else {
				s.addToken(common.EQUAL)
			}

		case common.LESS_SYMBOL:
			if s.matchRune(common.EQUAL_SYMBOL) {
				s.addToken(common.LESS_EQUAL)
			} else {
				s.addToken(common.LESS)
			}

		case common.GREATER_SYMBOL:
			if s.matchRune(common.EQUAL_SYMBOL) {
				s.addToken(common.GREATER_EQUAL)
			} else {
				s.addToken(common.GREATER)
			}

		case common.BANG_SYMBOL:
			if s.matchRune(common.EQUAL_SYMBOL) {
				s.addToken(common.BANG_EQUAL)
			} else {
				s.addToken(common.BANG)
			}

		case common.START_QUOTE_SYMBOL:
			start := s.index

			s.consumeWhile(func(r rune) bool { return r != common.END_QUOTE_SYMBOL })
			literalStr := s.content[start:s.index]

			if !s.isAtEnd() {
				s.advance()
				s.addTokenWithValue(common.LITERAL, literalStr)
			} else {
				s.addTokenWithValue(common.STARTED_LITERAL, literalStr)
			}

		case common.END_QUOTE_SYMBOL:
			s.clearTokens()
			s.addTokenWithValue(common.ENDED_LITERAL, s.content[:s.index-1])
			s.canMergeStart = true

		default:
			start := s.index - utf8.RuneLen(r)
			if common.IsAlphanumeric(r) {
				s.consumeWhile(common.IsAlphanumeric)
				lexeme := s.content[start:s.index]

				if tokenType, ok := common.KEYWORDS[lexeme]; ok {
					s.addTokenWithValue(tokenType, lexeme)
					continue
				}

				if isAllDigits(lexeme) {
					s.addTokenWithValue(common.NUMBER, lexeme)
					continue
				}

				s.addTokenWithValue(common.IDENTIFIER, lexeme)
			} else {
				return segment{}, fmt.Errorf("unexpected character: %c (%d)", r, r)
			}
		}
	}

	return segment{
		CouldMergeStart: s.canMergeStart,
		CouldMergeEnd:   s.canMergeEnd,
		Tokens:          s.tokens,
		Content:         s.content,
	}, nil
}
