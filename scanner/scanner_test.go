package scanner

import (
	"testing"

	"github.com/Tinchocw/forky/common"
)

type expectedToken struct {
	typ   common.TokenType
	value string
}

func checkTokens(t *testing.T, got []common.Token, exp []expectedToken) {
	t.Helper()
	if len(got) != len(exp) {
		t.Fatalf("token length mismatch: got %d expected %d\nGot: %v", len(got), len(exp), got)
	}
	for i, e := range exp {
		if got[i].Typ != e.typ {
			t.Fatalf("token %d type mismatch: got %s expected %s", i, got[i].Typ, e.typ)
		}
		if e.value != "" && got[i].Value != e.value {
			t.Fatalf("token %d value mismatch: got %q expected %q (type %s)", i, got[i].Value, e.value, got[i].Typ)
		}
	}
}

func workerVariants(input string) []int {
	l := len(input)
	out := []int{1}
	if l >= 2 {
		out = append(out, 2)
	}
	if l >= 3 {
		out = append(out, 3)
	}
	if l > 3 {
		out = append(out, l)
	}
	return out
}

func TestSingleCharTokens(t *testing.T) {
	input := "+-*/(),:;"
	expected := []expectedToken{{common.PLUS, ""}, {common.MINUS, ""}, {common.ASTERISK, ""}, {common.SLASH, ""}, {common.OPEN_PARENTHESIS, ""}, {common.CLOSE_PARENTHESIS, ""}, {common.COMMA, ""}, {common.COLON, ""}, {common.SEMICOLON, ""}}
	for _, w := range workerVariants(input) {
		toks, err := ScanString(input, w)
		if err != nil {
			t.Fatalf("scan error workers=%d: %v", w, err)
		}
		checkTokens(t, toks, expected)
	}
}

func TestMultiCharOperators(t *testing.T) {
	input := "== != <= >= = ! < >"
	expected := []expectedToken{{common.EQUAL_EQUAL, ""}, {common.BANG_EQUAL, ""}, {common.LESS_EQUAL, ""}, {common.GREATER_EQUAL, ""}, {common.EQUAL, ""}, {common.BANG, ""}, {common.LESS, ""}, {common.GREATER, ""}}
	for _, w := range workerVariants(input) {
		toks, err := ScanString(input, w)
		if err != nil {
			t.Fatalf("scan error workers=%d: %v", w, err)
		}
		checkTokens(t, toks, expected)
	}
}

func TestKeywordsIdentifiersNumbers(t *testing.T) {
	input := "if else while return true false none func var foo bar123 123abc 123 007"
	expected := []expectedToken{{common.IF, ""}, {common.ELSE, ""}, {common.WHILE, ""}, {common.RETURN, ""}, {common.TRUE, ""}, {common.FALSE, ""}, {common.NONE, ""}, {common.FUNC, ""}, {common.VAR, ""}, {common.IDENTIFIER, "foo"}, {common.IDENTIFIER, "bar123"}, {common.IDENTIFIER, "123abc"}, {common.NUMBER, "123"}, {common.NUMBER, "007"}}
	for _, w := range workerVariants(input) {
		toks, err := ScanString(input, w)
		if err != nil {
			t.Fatalf("scan error workers=%d: %v", w, err)
		}
		checkTokens(t, toks, expected)
	}
}

// TestLiteralsHappy validates properly terminated literals are tokenized.
func TestLiteralsHappy(t *testing.T) {
	input := "\"hi' \"multi word literal'" // two complete literals
	expected := []expectedToken{{common.LITERAL, "hi"}, {common.LITERAL, "multi word literal"}}
	for _, w := range workerVariants(input) {
		toks, err := ScanString(input, w)
		if err != nil {
			t.Fatalf("unexpected error workers=%d: %v", w, err)
		}
		checkTokens(t, toks, expected)
	}
}

// TestLiteralErrors ensures unterminated literals produce an error instead of STARTED_LITERAL tokens.
func TestLiteralErrors(t *testing.T) {
	// Inputs ending with an open double quote without matching single quote terminator (per current design).
	cases := []string{
		"\"unterminated",
		"prefix \"stillopen",
		"\"", // empty literal not closed
	}
	for _, input := range cases {
		for _, w := range workerVariants(input) {
			_, err := ScanString(input, w)
			if err == nil {
				t.Fatalf("expected error for unterminated literal %q workers=%d, got none", input, w)
			}
		}
	}
}

func TestCrossBoundaryIdentifier(t *testing.T) {
	ident := "AlphaBetaGammaDelta" // reasonably long
	expected := []expectedToken{{common.IDENTIFIER, ident}}
	toks, err := ScanString(ident, len(ident))
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}
	checkTokens(t, toks, expected)
}

func TestCrossBoundaryOperators(t *testing.T) {
	input := "== != <= >="
	expected := []expectedToken{{common.EQUAL_EQUAL, ""}, {common.BANG_EQUAL, ""}, {common.LESS_EQUAL, ""}, {common.GREATER_EQUAL, ""}}
	toks, err := ScanString(input, 8)
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}
	checkTokens(t, toks, expected)
}
