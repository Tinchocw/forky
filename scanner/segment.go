package scanner

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type segment struct {
	CouldMergeStart bool
	CouldMergeEnd   bool
	Tokens          []common.Token
	Content         string
}

func NewSegment(content string) segment {
	return segment{
		CouldMergeStart: true,
		CouldMergeEnd:   true,
		Tokens:          []common.Token{},
		Content:         content,
	}
}

func (s *segment) String() string {
	b := fmt.Sprintf("Segment: tokens=%d canMergeStart=%v canMergeEnd=%v\n", len(s.Tokens), s.CouldMergeStart, s.CouldMergeEnd)
	b += fmt.Sprintf("Content: '%s'\n", s.Content)
	b += "Tokens:\n"
	for _, t := range s.Tokens {
		b += fmt.Sprintf("  %s\n", t.String())
	}
	return b
}

func (s *segment) AddToken(token common.Token) {
	s.Tokens = append(s.Tokens, token)
}

func (s *segment) AddTokens(tokens []common.Token) {
	s.Tokens = append(s.Tokens, tokens...)
}

func (s *segment) hasInvalidTokens() bool {
	if !s.hasTokens() {
		return false
	}

	if s.firstToken().Typ == common.ENDED_LITERAL {
		return true
	}

	if s.lastToken().Typ == common.STARTED_LITERAL {
		return true
	}

	return false
}

func (s *segment) hasTokens() bool {
	return len(s.Tokens) > 0
}

func (s *segment) lastToken() *common.Token {
	if !s.hasTokens() {
		panic("No tokens available")
	}
	return &s.Tokens[len(s.Tokens)-1]
}

func (s *segment) firstToken() *common.Token {
	if !s.hasTokens() {
		panic("No tokens available")
	}
	return &s.Tokens[0]
}

func (s *segment) consume(n int) []common.Token {
	if n > len(s.Tokens) {
		panic("Cannot consume more tokens than available")
	}

	consumed := s.Tokens[:n]
	s.Tokens = s.Tokens[n:]
	return consumed
}

func (s *segment) consumeOne() common.Token {
	return s.consume(1)[0]
}

func (s *segment) clearTokens() {
	s.Tokens = []common.Token{}
}

func (current *segment) Merge(other *segment) {
	defer func() {
		current.Content += other.Content
	}()

	if !other.hasTokens() {
		current.CouldMergeEnd = other.CouldMergeEnd
		return
	}

	if other.firstToken().Typ == common.ENDED_LITERAL {
		if current.CouldMergeEnd && current.hasTokens() && current.lastToken().Typ == common.STARTED_LITERAL {
			current.lastToken().Typ = common.LITERAL
			current.lastToken().Value += other.consumeOne().Value
		} else {
			other.lastToken().Value = current.Content + other.lastToken().Value
			current.clearTokens()
		}

		current.AddTokens(other.Tokens)
		current.CouldMergeStart = true
		current.CouldMergeEnd = other.CouldMergeEnd
		return
	}

	if !current.CouldMergeEnd || !current.hasTokens() {
		current.Tokens = append(current.Tokens, other.Tokens...)
		current.CouldMergeEnd = other.CouldMergeEnd
		return
	}

	if other.CouldMergeStart {
		switch current.lastToken().Typ {
		case common.EQUAL:
			if other.firstToken().Typ == common.EQUAL {
				current.lastToken().Typ = common.EQUAL_EQUAL
				other.consumeOne()
			}
		case common.BANG:
			if other.firstToken().Typ == common.EQUAL {
				current.lastToken().Typ = common.BANG_EQUAL
				other.consumeOne()
			}
		case common.LESS:
			if other.firstToken().Typ == common.EQUAL {
				current.lastToken().Typ = common.LESS_EQUAL
				other.consumeOne()
			}
		case common.GREATER:
			if other.firstToken().Typ == common.EQUAL {
				current.lastToken().Typ = common.GREATER_EQUAL
				other.consumeOne()
			}
		case common.NUMBER, common.IDENTIFIER:
			if other.firstToken().Typ == common.NUMBER || other.firstToken().Typ == common.IDENTIFIER {
				if other.firstToken().Typ != current.lastToken().Typ {
					current.lastToken().Typ = common.IDENTIFIER
				}

				current.lastToken().Value += other.consumeOne().Value
			} else if keyword, ok := common.KEYWORDS_VALUES[other.firstToken().Typ]; ok {
				other.consumeOne()
				current.lastToken().Value += keyword
			}
		}

		if currentKeyword, ok := common.KEYWORDS_VALUES[current.lastToken().Typ]; ok {
			if other.firstToken().Typ == common.NUMBER || other.firstToken().Typ == common.IDENTIFIER {
				current.lastToken().Typ = common.IDENTIFIER
				current.lastToken().Value = currentKeyword + other.consumeOne().Value
			} else if otherKeyword, ok := common.KEYWORDS_VALUES[other.firstToken().Typ]; ok {
				other.consumeOne()
				current.lastToken().Typ = common.IDENTIFIER
				current.lastToken().Value = currentKeyword + otherKeyword
			}

			// maybe handle the case of continuous keywords with operators or other stuff
		}

		// If after merging the last token is a keyword, update its type
		if current.lastToken().Typ == common.IDENTIFIER {
			if keywordType, ok := common.KEYWORDS[current.lastToken().Value]; ok {
				current.lastToken().Typ = keywordType
			}
		}
	}

	if current.lastToken().Typ == common.STARTED_LITERAL {
		if other.firstToken().Typ == common.ENDED_LITERAL {
			current.firstToken().Value += other.consumeOne().Value
		} else {
			current.Tokens[len(current.Tokens)-1].Value += other.Content
			other.clearTokens()
		}
	}

	current.Tokens = append(current.Tokens, other.Tokens...)
	current.CouldMergeEnd = other.CouldMergeEnd
}
