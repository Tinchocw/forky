package scanner

import (
	"fmt"
	"io"
	"os"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

// (No external imports required for scanner logic.)

type scanner struct {
	content       string
	tokens        []common.Token
	index         int
	canMergeStart bool
	canMergeEnd   bool
}

func CreateScanner(content string) scanner {
	return scanner{
		content:       content,
		tokens:        []common.Token{},
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

func (s *scanner) Scan() error {
	for !s.isAtEnd() {
		r := s.advance()

		if common.IsWhitespace(r) {
			if s.isAtStart() {
				s.canMergeStart = false
			}

			if s.isAtEnd() {
				s.canMergeEnd = false
			}
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
			start := s.index - 1
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (r >= '0' && r <= '9') {
				s.consumeWhile(func(r2 rune) bool {
					return (r2 >= 'a' && r2 <= 'z') || (r2 >= 'A' && r2 <= 'Z') || r2 == '_' || (r2 >= '0' && r2 <= '9')
				})
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
				return fmt.Errorf("unexpected character: %c", r)
			}
		}
	}

	return nil
}

// Public accessors
func (s *scanner) Tokens() []common.Token { return s.tokens }
func (s *scanner) CanMergeStart() bool    { return s.canMergeStart }
func (s *scanner) CanMergeEnd() bool      { return s.canMergeEnd }

// ----- Concurrent high-level API -----

// segment is an internal chunk scanned independently.
type segment struct {
	id              int
	couldMergeStart bool
	couldMergeEnd   bool
	tokens          []common.Token
	content         string
	err             error
}

// splitOffsets divides a total size into up to workers nearly-even ranges.
func splitOffsets(size int64, workers int) [][2]int64 {
	if workers < 1 {
		workers = 1
	}
	part := (size + int64(workers) - 1) / int64(workers)
	out := make([][2]int64, 0, workers)
	var start int64
	for start < size {
		end := min(start+part, size)
		out = append(out, [2]int64{start, end})
		start = end
	}
	return out
}

// ScanFile scans a file with the given number of workers.
func ScanFile(path string, workers int) ([]common.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return scanReaderConcurrently(f, st.Size(), workers)
}

// ScanBytes scans an in-memory byte slice.
func ScanBytes(data []byte, workers int) ([]common.Token, error) {
	return scanReaderConcurrently(bytesReader(data), int64(len(data)), workers)
}

// bytesReader implements io.ReaderAt for a byte slice.
type bytesReader []byte

func (b bytesReader) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(b)) {
		return 0, io.EOF
	}
	n := copy(p, b[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func scanReaderConcurrently(r io.ReaderAt, size int64, workers int) ([]common.Token, error) {
	offsets := splitOffsets(size, workers)
	segCh := make(chan segment, len(offsets))
	for i, off := range offsets {
		iLocal := i
		start := off[0]
		end := off[1]
		go func() {
			length := end - start
			buf := make([]byte, length)
			_, err := r.ReadAt(buf, start)
			if err != nil && err != io.EOF {
				segCh <- segment{id: iLocal, err: err}
				return
			}
			sc := CreateScanner(string(buf))
			if err := sc.Scan(); err != nil {
				segCh <- segment{id: iLocal, err: err}
				return
			}
			segCh <- segment{
				id:              iLocal,
				couldMergeStart: sc.CanMergeStart(),
				couldMergeEnd:   sc.CanMergeEnd(),
				tokens:          sc.Tokens(),
				content:         string(buf),
			}
		}()
	}

	segs := make([]segment, len(offsets))
	for range offsets {
		s := <-segCh
		if s.err != nil {
			return nil, s.err
		}
		segs[s.id] = s
	}

	m := newMerger(segs)
	if err := m.merge(); err != nil {
		return nil, err
	}
	return m.tokensOut(), nil
}
