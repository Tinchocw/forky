package main

import (
	"io"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/statement"
	"github.com/Tinchocw/Interprete-concurrente/interpreter"
	parserPackage "github.com/Tinchocw/Interprete-concurrente/parser"
	scannerPackage "github.com/Tinchocw/Interprete-concurrente/scanner"
)

type InterpreterMode int

const (
	NormalMode InterpreterMode = iota
	ScanningMode
	ParsingMode
	// ResolveMode
)

// Forky is the top-level runner that coordinates the scanning (and future phases).
type Forky struct {
	workers     int
	debug       bool
	mode        InterpreterMode
	interpreter *interpreter.Interpreter
}

func NewForky(workers int, debug bool, mode InterpreterMode) *Forky {
	if workers < 1 {
		workers = 1
	}
	i := interpreter.NewInterpreter()
	return &Forky{workers: workers, debug: debug, mode: mode, interpreter: &i}
}

// Run executes the configured mode against the provided ReaderAt of given size.
func (forky *Forky) Run(r io.ReaderAt, size int64) (string, error) {
	sc := scannerPackage.CreateForkyScanner(forky.workers, forky.debug)

	// Read entire input into memory and scan from bytes
	buf := make([]byte, size)
	if size > 0 {
		if _, err := r.ReadAt(buf, 0); err != nil && err != io.EOF {
			return "", err
		}
	}

	tokens, err := sc.Scan(r, size)
	if err != nil {
		return "", err
	}

	if forky.mode == ScanningMode {
		common.PrintTokens(tokens)
		return "", nil
	}

	ps := parserPackage.CreateForkyParser(forky.workers, forky.debug)
	program, err := ps.Parse(tokens)
	if err != nil {
		return "", err
	}

	if forky.mode == ParsingMode {
		statement.PrintProgram(program)
		return "", nil
	}

	return forky.interpreter.Execute(program)
}
