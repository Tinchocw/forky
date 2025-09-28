package parser

import "github.com/Tinchocw/Interprete-concurrente/common"

type segment struct {
	CouldMergeStart bool
	CouldMergeEnd   bool
	Program         common.PartialProgram
	Tokens          []common.Token
}

func NewSegment() segment {
	return segment{
		CouldMergeStart: true,
		CouldMergeEnd:   true,
		Program:         common.PartialProgram{Statements: []common.IncompleteStatement{}},
		Tokens:          []common.Token{},
	}
}

func (s *segment) AddStatement(content common.IncompleteStatement) {
	s.Program.Statements = append(s.Program.Statements, content)
}
func (s *segment) AddStatements(contents []common.IncompleteStatement) {
	s.Program.Statements = append(s.Program.Statements, contents...)
}

func (s *segment) hasStatements() bool {
	return len(s.Program.Statements) > 0
}

func (s *segment) hasInvalidStatements() bool {

	if !s.hasStatements() {
		return false
	}

	return false
}
