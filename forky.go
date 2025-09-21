package main

import (
	"fmt"
	"io"

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
	workers int
	debug   bool
	mode    InterpreterMode
}

func NewForky(workers int, debug bool, mode InterpreterMode) *Forky {
	if workers < 1 {
		workers = 1
	}
	return &Forky{workers: workers, debug: debug, mode: mode}
}

// Run executes the configured mode against the provided ReaderAt of given size.
func (fj *Forky) Run(r io.ReaderAt, size int64) error {
	sc := scannerPackage.CreateForkyScanner(fj.workers)

	// Read entire input into memory and scan from bytes
	buf := make([]byte, size)
	if size > 0 {
		if _, err := r.ReadAt(buf, 0); err != nil && err != io.EOF {
			return err
		}
	}

	tokens, err := sc.Scan(r, size)
	if err != nil {
		return err
	}

	if fj.mode == ScanningMode {
		fmt.Println()
		fmt.Println("== Tokens ==")

		for i, t := range tokens {
			fmt.Printf("%4d: %s\n", i, t.ColorString())
		}

		fmt.Println("== End of Tokens ==")
		fmt.Println()
	}

	ps := parserPackage.NewParser(tokens)
	program, err := ps.Parse()
	if err != nil {
		return err
	}

	if fj.mode == ParsingMode {
		fmt.Println()
		fmt.Println("== PROGRAM ==")
		fmt.Println(program)
		fmt.Println("== End of PROGRAM ==")
		fmt.Println()
	}

	return nil
}
