package main

import (
	"fmt"
	"io"

	scannerPackage "github.com/Tinchocw/Interprete-concurrente/scanner"
)

type InterpreterMode int

const (
	NormalMode InterpreterMode = iota
	ScanningMode
	// ParsingMode
	// ResolveMode
)

// ForkJoiner is the top-level runner that coordinates the scanning (and future phases).
type ForkJoiner struct {
	workers int
	debug   bool
	mode    InterpreterMode
}

func NewForkJoiner(workers int, debug bool, mode InterpreterMode) *ForkJoiner {
	if workers < 1 {
		workers = 1
	}
	return &ForkJoiner{workers: workers, debug: debug, mode: mode}
}

// Run executes the configured mode against the provided ReaderAt of given size.
func (fj *ForkJoiner) Run(r io.ReaderAt, size int64) error {
	sc := scannerPackage.CreateForkJoinScanner(fj.workers)

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
		fmt.Println("== Tokens ==")

		for i, t := range tokens {
			fmt.Printf("%4d: %s\n", i, t.ColorString())
		}

		fmt.Println("== End of Tokens ==")
	}

	return nil
}
