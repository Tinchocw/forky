package parser

import (
	"github.com/Tinchocw/Interprete-concurrente/common"
)

type ForkyParser struct {
	numWorkers int
}

func NewForkyParser(numWorkers int) *ForkyParser {
	return &ForkyParser{numWorkers: numWorkers}
}

func (fp *ForkyParser) Parse(tokens []common.Token) common.Program {
	parser := NewParser(tokens, false)
	program, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	return program
}
