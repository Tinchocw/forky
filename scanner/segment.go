package scanner

import "github.com/Tinchocw/Interprete-concurrente/common"

type Segment struct {
	ID              int
	CouldMergeStart bool
	CouldMergeEnd   bool
	Tokens          []common.Token
	Content         string
	Err             error
}

func NewSegment(id int, content string) Segment {
	return Segment{
		ID:              id,
		CouldMergeStart: true,
		CouldMergeEnd:   true,
		Tokens:          []common.Token{},
		Content:         content,
		Err:             nil,
	}
}

func (s *Segment) AddToken(token common.Token) {
	s.Tokens = append(s.Tokens, token)
}

func (s *Segment) AddTokens(tokens []common.Token) {
	s.Tokens = append(s.Tokens, tokens...)
}

func (s *Segment) hasTokens() bool {
	return len(s.Tokens) > 0
}

func (s *Segment) lastToken() *common.Token {
	if !s.hasTokens() {
		panic("No tokens available")
	}
	return &s.Tokens[len(s.Tokens)-1]
}

func (s *Segment) firstToken() *common.Token {
	if !s.hasTokens() {
		panic("No tokens available")
	}
	return &s.Tokens[0]
}

func (s *Segment) consume(n int) []common.Token {
	if n > len(s.Tokens) {
		panic("Cannot consume more tokens than available")
	}

	consumed := s.Tokens[:n]
	s.Tokens = s.Tokens[n:]
	return consumed
}

func (s *Segment) consumeOne() common.Token {
	return s.consume(1)[0]
}

func (s *Segment) clearTokens() {
	s.Tokens = []common.Token{}
}

func (current *Segment) Merge(other *Segment) {
	defer func() {
		current.Content += other.Content
	}()

	if !other.hasTokens() {
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
		return
	}

	if !current.CouldMergeEnd || !current.hasTokens() {
		current.Tokens = append(current.Tokens, other.Tokens...)
		return
	}

	if other.CouldMergeStart {
		switch other.firstToken().Typ {
		case common.EQUAL:
			if current.lastToken().Typ == common.EQUAL {
				current.lastToken().Typ = common.EQUAL_EQUAL
				other.consumeOne()
			}
		case common.BANG:
			if current.lastToken().Typ == common.EQUAL {
				current.lastToken().Typ = common.BANG_EQUAL
				other.consumeOne()
			}
		case common.LESS:
			if current.lastToken().Typ == common.EQUAL {
				current.lastToken().Typ = common.LESS_EQUAL
				other.consumeOne()
			}
		case common.GREATER:
			if current.lastToken().Typ == common.EQUAL {
				current.lastToken().Typ = common.GREATER_EQUAL
				other.consumeOne()
			}
		case common.NUMBER, common.IDENTIFIER:
			if other.firstToken().Typ == common.NUMBER || other.firstToken().Typ == common.IDENTIFIER {
				if other.firstToken().Typ != current.lastToken().Typ {
					current.firstToken().Typ = common.IDENTIFIER
				}

				current.lastToken().Value += other.consumeOne().Value
			} else if keyword, ok := common.KEYWORDS_VALUES[other.firstToken().Typ]; ok {
				other.consumeOne()
				current.lastToken().Value += keyword
			}
		}

		if currentKeyword, ok := common.KEYWORDS_VALUES[current.lastToken().Typ]; ok {
			if other.firstToken().Typ == common.NUMBER || other.firstToken().Typ == common.IDENTIFIER {
				current.firstToken().Typ = common.IDENTIFIER
				current.lastToken().Value = currentKeyword + other.consumeOne().Value
			} else if otherKeyword, ok := common.KEYWORDS_VALUES[other.firstToken().Typ]; ok {
				other.consumeOne()
				current.firstToken().Typ = common.IDENTIFIER
				current.lastToken().Value = currentKeyword + otherKeyword
			}
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
