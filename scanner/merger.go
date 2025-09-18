package scanner

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

// merger combines tokens from independently scanned segments handling
// multi-character operators, identifiers/numbers split across boundaries,
// and multi-part literals (STARTED_LITERAL ... ENDED_LITERAL).
// It is intentionally unexported; the public API is ScanFile / ScanBytes.
type merger struct {
	results []segment
	index   int
	tokens  []common.Token
}

func newMerger(results []segment) *merger {
	return &merger{results: results, index: 0, tokens: []common.Token{}}
}

func (m *merger) isAtEnd() bool { return m.index >= len(m.results) }

func (m *merger) lookahead() *segment {
	if m.isAtEnd() {
		panic("No more parts to lookahead")
	}
	return &m.results[m.index]
}

func (m *merger) advance() *segment {
	part := m.lookahead()
	m.index++
	return part
}

func (m *merger) addTokens(tokens []common.Token) { m.tokens = append(m.tokens, tokens...) }
func (m *merger) addToken(token common.Token)     { m.tokens = append(m.tokens, token) }

func (m *merger) match(typ common.TokenType) bool {
	if m.isAtEnd() {
		return false
	}
	next := m.lookahead()
	if len(next.tokens) == 0 || next.couldMergeStart || next.tokens[0].Typ != typ {
		return false
	}
	next.tokens = next.tokens[1:]
	return true
}

func (m *merger) merge() error {
general:
	for !m.isAtEnd() {
		part := m.advance()
		if len(part.tokens) == 0 {
			continue
		}

		if part.tokens[0].Typ == common.ENDED_LITERAL {
			fmt.Printf("Part %d starts with ENDED_LITERAL token: %v\n", m.index-1, part.tokens[0])
			return fmt.Errorf("part %d starts with ENDED_LITERAL token, cannot merge", m.index-1)
		}

		if !part.couldMergeEnd {
			m.addTokens(part.tokens)
			continue
		}

		// Keep all but last token; process last for potential merging.
		m.addTokens(part.tokens[:len(part.tokens)-1])
		token := part.tokens[len(part.tokens)-1]

		switch token.Typ {
		case common.EQUAL:
			if m.match(common.EQUAL) {
				m.addToken(common.NewToken(common.EQUAL_EQUAL))
			} else {
				m.addToken(token)
			}
		case common.LESS:
			if m.match(common.EQUAL) {
				m.addToken(common.NewToken(common.LESS_EQUAL))
			} else {
				m.addToken(token)
			}
		case common.GREATER:
			if m.match(common.EQUAL) {
				m.addToken(common.NewToken(common.GREATER_EQUAL))
			} else {
				m.addToken(token)
			}
		case common.BANG:
			if m.match(common.EQUAL) {
				m.addToken(common.NewToken(common.BANG_EQUAL))
			} else {
				m.addToken(token)
			}
		case common.NUMBER:
			for !m.isAtEnd() {
				next := m.lookahead()
				if len(next.tokens) == 0 || !next.couldMergeStart || (next.tokens[0].Typ != common.NUMBER && next.tokens[0].Typ != common.IDENTIFIER) {
					break
				}
				nextToken := next.tokens[0]
				if nextToken.Typ == common.IDENTIFIER {
					token.Typ = common.IDENTIFIER
				}
				next.tokens = next.tokens[1:]
				token.Value += nextToken.Value
				if !next.couldMergeEnd {
					break
				}
			}
			m.addToken(token)
		case common.IDENTIFIER:
			for !m.isAtEnd() {
				next := m.lookahead()
				if len(next.tokens) == 0 || !next.couldMergeStart || (next.tokens[0].Typ != common.IDENTIFIER && next.tokens[0].Typ != common.NUMBER) {
					break
				}
				nextToken := next.tokens[0]
				next.tokens = next.tokens[1:]
				token.Value += nextToken.Value
				if !next.couldMergeEnd || len(next.tokens) > 0 {
					break
				}
				m.index++
			}
			if keywordType, ok := common.KEYWORDS[token.Value]; ok {
				m.addToken(common.NewTokenWithValue(keywordType, token.Value))
			} else {
				m.addToken(token)
			}
		case common.STARTED_LITERAL:
			start := m.index
			for !m.isAtEnd() {
				next := m.lookahead()
				if len(next.tokens) == 0 || !next.couldMergeStart || next.tokens[0].Typ != common.ENDED_LITERAL {
					m.index++
					continue
				}
				newToken := common.NewTokenWithValue(common.LITERAL, token.Value)
				endToken := next.tokens[0]
				next.tokens = next.tokens[1:]
				for i := start; i < m.index; i++ {
					newToken.Value += m.results[i].content
				}
				newToken.Value += endToken.Value
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

func (m *merger) tokensOut() []common.Token { return m.tokens }
