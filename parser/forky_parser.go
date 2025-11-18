package parser

import (
	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/statement"
)

type ForkyParser struct {
	numWorkers int
	debug      bool
}

func CreateForkyParser(numWorkers int, debug bool) *ForkyParser {
	return &ForkyParser{numWorkers: numWorkers, debug: debug}
}

func (fp *ForkyParser) Parse(tokens []common.Token) (statement.Program, error) {
	p := NewParser(tokens, fp.debug)
	return p.parse()
}
