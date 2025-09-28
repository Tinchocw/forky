package parser

import (
	"github.com/Tinchocw/Interprete-concurrente/common"
)

type ForkyParser struct {
	numWorkers int
}

func ParallelParse(tokens []common.Token, numWorkers int, debug bool) (segment, error) {

	length := len(tokens)

	if debug {
		println("[DEBUG] ParallelParse: total tokens =", length, "numWorkers =", numWorkers)
	}

	if length == 0 || numWorkers <= 0 {
		return segment{}, nil
	}

	if numWorkers == 1 || length == 1 {
		if debug {
			println("[DEBUG] Sequential parse: total tokens =", length)
		}

		parser := NewParser(tokens, debug)
		sg, err := parser.parse()
		if err != nil {
			return segment{}, err
		}

		if debug {
			println("[DEBUG] Sequential parse result: AST with", len(sg.Program.Statements), "statements")
		}

		return sg, err
	}

	leftWorkers := (numWorkers + 1) / 2 // ceil(workers/2) ensures left >= right when odd
	rightWorkers := numWorkers - leftWorkers

	mid := (length * leftWorkers) / numWorkers

	type res struct {
		sg  segment
		err error
	}

	leftCh := make(chan res, 1)
	go func() {
		if debug {
			println("[DEBUG] Left parse: tokens [0:", mid, ") length =", mid, "workers =", leftWorkers)
		}

		leftTokens := tokens[:mid]
		leftSegment, err := ParallelParse(leftTokens, leftWorkers, debug)
		if err != nil {
			return
		}

		leftCh <- res{sg: leftSegment, err: nil}
	}()

	if debug {
		println("[DEBUG] Right parse: tokens [", mid, ":", length, ") length =", length-mid, "workers =", rightWorkers)
	}

	rightSegment, rightErr := ParallelParse(tokens[mid:], rightWorkers, debug)
	if rightErr != nil {
		return segment{}, rightErr
	}
	leftRes := <-leftCh

	if leftRes.err != nil {
		return segment{}, leftRes.err
	}

	if debug {
		println("[DEBUG] Merging segments: left tokens =", len(leftRes.sg.Tokens), "right tokens =", len(rightSegment.Tokens))
	}

	//leftRes.sg.Merge(&rightSegment)

	if debug {
		println("[DEBUG] Merge complete: total tokens =", len(leftRes.sg.Tokens))
	}

	return leftRes.sg, nil

}

func CreateForkyParser(numWorkers int) *ForkyParser {
	return &ForkyParser{numWorkers: numWorkers}
}

func (fp *ForkyParser) Parse(tokens []common.Token, debug bool) (common.PartialProgram, error) {

	sg, err := ParallelParse(tokens, fp.numWorkers, debug)
	if err != nil {
		return common.PartialProgram{}, err
	}
	// i have to change this to return a Program
	return sg.Program, nil
}
