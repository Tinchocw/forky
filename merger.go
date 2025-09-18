package main

import "fmt"

type merger struct {
	results []WorkerResult
	index   int
	tokens  []Token
}

func newMerger(results []WorkerResult) merger {
	return merger{
		results: results,
		index:   0,
		tokens:  []Token{},
	}
}

func (m *merger) isAtEnd() bool {
	return m.index >= len(m.results)
}

func (m *merger) lookahead() WorkerResult {
	if m.isAtEnd() {
		panic("No more parts to lookahead")
	}
	return m.results[m.index]
}

func (m *merger) advance() WorkerResult {
	part := m.lookahead()
	m.index++
	return part
}

func (m *merger) addTokens(tokens []Token) {
	m.tokens = append(m.tokens, tokens...)
}

func (m *merger) addToken(token Token) {
	m.tokens = append(m.tokens, token)
}

func (m *merger) match(typ TokenType) bool {
	if m.isAtEnd() {
		return false
	}
	nextPart := m.lookahead()
	if len(nextPart.tokens) == 0 || nextPart.couldMergeStart || nextPart.tokens[0].typ != typ {
		return false
	}

	nextPart.tokens = nextPart.tokens[1:]
	return true
}

func (m *merger) merge() error {
general:
	for !m.isAtEnd() {
		part := m.advance()

		if len(part.tokens) == 0 {
			continue
		}

		if part.tokens[0].typ == ENDED_LITERAL {

			fmt.Printf("Part %d starts with ENDED_LITERAL token: %v\n", m.index-1, part.tokens[0])
			return fmt.Errorf("part %d starts with ENDED_LITERAL token, cannot merge", m.index-1)
		}

		if !part.couldMergeEnd {
			m.addTokens(part.tokens)
			continue
		}

		m.addTokens(part.tokens[:len(part.tokens)-1])
		token := part.tokens[len(part.tokens)-1]

		switch token.typ {

		case EQUAL:
			if m.match(EQUAL) {
				m.addToken(NewToken(EQUAL_EQUAL))
			} else {
				m.addToken(token)
			}
		case LESS:
			if m.match(EQUAL) {
				m.addToken(NewToken(LESS_EQUAL))
			} else {
				m.addToken(token)
			}
		case GREATER:
			if m.match(EQUAL) {
				m.addToken(NewToken(GREATER_EQUAL))
			} else {
				m.addToken(token)
			}
		case MARK:
			if m.match(EQUAL) {
				m.addToken(NewToken(MARK_EQUAL))
			} else {
				m.addToken(token)
			}
		case NUMBER:
			for !m.isAtEnd() {
				nextPart := m.lookahead()
				if len(nextPart.tokens) == 0 || !nextPart.couldMergeStart || (nextPart.tokens[0].typ != NUMBER && nextPart.tokens[0].typ != IDENTIFIER) {
					break
				}

				nextToken := nextPart.tokens[0]

				if nextToken.typ == IDENTIFIER {
					token.typ = IDENTIFIER
				}

				nextPart.tokens = nextPart.tokens[1:]
				token.value += nextToken.value

				if !nextPart.couldMergeEnd {
					break
				}
			}

			m.addToken(token)
		case IDENTIFIER:
			for !m.isAtEnd() {
				nextPart := m.lookahead()
				if len(nextPart.tokens) == 0 || !nextPart.couldMergeStart || (nextPart.tokens[0].typ != IDENTIFIER && nextPart.tokens[0].typ != NUMBER) {
					break
				}

				nextToken := nextPart.tokens[0]
				nextPart.tokens = nextPart.tokens[1:]
				m.results[m.index].tokens = nextPart.tokens
				token.value += nextToken.value

				if !nextPart.couldMergeEnd || len(nextPart.tokens) > 0 {
					break
				}

				m.index++
			}

			if keywordType, ok := KEYWORDS[token.value]; ok {
				m.addToken(NewTokenWithValue(keywordType, token.value))
			} else {
				m.addToken(token)
			}
		case STARTED_LITERAL:
			start := m.index
			for !m.isAtEnd() {
				nextPart := m.lookahead()
				if len(nextPart.tokens) == 0 || !nextPart.couldMergeStart || nextPart.tokens[0].typ != ENDED_LITERAL {
					m.index++
					continue
				}

				newToken := NewTokenWithValue(LITERAL, token.value)
				endToken := nextPart.tokens[0]
				m.results[m.index].tokens = nextPart.tokens[1:]

				for i := start; i < m.index; i++ {
					newToken.value += m.results[i].content
				}

				newToken.value += endToken.value
				m.addToken(newToken)
				continue general
			}
			return fmt.Errorf("no ENDED_LITERAL token found to close STARTED_LITERAL starting at part %d", m.index-1)
		default:
			m.addToken(token)
		}
	}
	return nil
}
