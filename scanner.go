package main

// (No external imports required for scanner logic.)

type scanner struct {
	content       string
	tokens        []Token
	index         int
	canMergeStart bool
	canMergeEnd   bool
}

func createScanner(content string) scanner {
	return scanner{
		content:       content,
		tokens:        []Token{},
		index:         0,
		canMergeStart: true,
		canMergeEnd:   true,
	}
}

func (s *scanner) isAtStart() bool {
	return s.index == 0
}

func (s *scanner) isAtEnd() bool {
	return s.index >= len(s.content)
}

func (s *scanner) lookahead() rune {
	if s.isAtEnd() {
		panic("No more characters to lookahead")
	}
	return rune(s.content[s.index])
}

func (s *scanner) advance() rune {
	lookahead := s.lookahead()
	s.index++
	return lookahead
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

func (s *scanner) addToken(typ TokenType) {
	token := NewToken(typ)
	s.tokens = append(s.tokens, token)
}

func (s *scanner) addTokenWithValue(typ TokenType, value string) {
	token := NewTokenWithValue(typ, value)
	s.tokens = append(s.tokens, token)
}

func (s *scanner) clearTokens() {
	s.tokens = []Token{}
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
		if !isNumber(r) {
			return false
		}
	}
	return true
}

func (s *scanner) scan() error {
	for !s.isAtEnd() {
		r := s.advance()

		if isWhitespace(r) {
			if s.isAtStart() {
				s.canMergeStart = false
			}

			if s.isAtEnd() {
				s.canMergeEnd = false
			}
		}

		switch r {
		case PLUS_SYMBOL:
			s.addToken(PLUS)

		case MINUS_SYMBOL:
			s.addToken(MINUS)

		case ASTERISK_SYMBOL:
			s.addToken(ASTERISK)

		case SLASH_SYMBOL:
			s.addToken(SLASH)

		case EQUAL_SYMBOL:
			if s.matchRune(EQUAL_SYMBOL) {
				s.addToken(EQUAL_EQUAL)
			} else {
				s.addToken(EQUAL)
			}

		case LESS_SYMBOL:
			if s.matchRune(EQUAL_SYMBOL) {
				s.addToken(LESS_EQUAL)
			} else {
				s.addToken(LESS)
			}

		case GREATER_SYMBOL:
			if s.matchRune(EQUAL_SYMBOL) {
				s.addToken(GREATER_EQUAL)
			} else {
				s.addToken(GREATER)
			}

		case MARK_SYMBOL:
			if s.matchRune(EQUAL_SYMBOL) {
				s.addToken(MARK_EQUAL)
			} else {
				s.addToken(MARK)
			}

		case START_QUOTE:
			start := s.index

			s.consumeWhile(func(r rune) bool { return r != END_QUOTE })
			literalStr := s.content[start:s.index]

			if !s.isAtEnd() {
				s.advance()
				s.addTokenWithValue(LITERAL, literalStr)
			} else {
				s.addTokenWithValue(STARTED_LITERAL, literalStr)
			}

		case END_QUOTE:
			s.clearTokens()
			s.addTokenWithValue(ENDED_LITERAL, s.content[:s.index-1])
			s.canMergeStart = true

		default:
			start := s.index - 1
			if isAlphanumeric(r) {
				s.consumeWhile(isAlphanumeric)
				lexeme := s.content[start:s.index]

				if tokenType, ok := KEYWORDS[lexeme]; ok {
					s.addTokenWithValue(tokenType, lexeme)
					continue
				}

				if isAllDigits(lexeme) {
					s.addTokenWithValue(NUMBER, lexeme)
					continue
				}

				s.addTokenWithValue(IDENTIFIER, lexeme)
			}
		}
	}

	return nil
}
