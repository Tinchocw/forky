package main

import (
	"io"
	"strings"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/statement"
	"github.com/Tinchocw/forky/interpreter"
	"github.com/Tinchocw/forky/parser"
	"github.com/Tinchocw/forky/scanner"
	"github.com/peterh/liner"
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
	sc := scanner.CreateForkyScanner(forky.workers, forky.debug)

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

	ps := parser.CreateForkyParser(forky.workers, forky.debug)
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

func (f *Forky) WordCompleter() liner.WordCompleter {
	keywords := make([]string, 0, len(common.KEYWORDS))
	for k := range common.KEYWORDS {
		keywords = append(keywords, k)
	}

	return func(line string, pos int) (string, []string, string) {
		globalVars := f.interpreter.GetGlobalVariables()
		options := make([]string, 0, len(keywords)+len(globalVars))
		options = append(options, keywords...)
		options = append(options, globalVars...)

		// Find the current word being typed
		start := pos
		for start > 0 && line[start-1] != ' ' && line[start-1] != '\t' {
			start--
		}
		end := pos
		for end < len(line) && line[end] != ' ' && line[end] != '\t' {
			end++
		}
		word := line[start:end]

		// Filter completions that start with the current word
		completions := []string{}

		for _, opt := range options {
			if strings.HasPrefix(opt, word) {
				completions = append(completions, opt)
			}
		}

		return line[:start], completions, line[end:]
	}
}
