package scanner

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type ForkyScanner struct {
	numWorkers int
	debug      bool
}

// parallelScan performs a fork-join recursive scan. It splits the range
// [start,end) roughly in half while there is budget (>1) of workers left
// and the segment length is > 1. Each recursive branch gets a share of the
// remaining worker budget. When budget == 1 or range is minimal, it scans
// sequentially in-place.
func parallelScan(r io.ReaderAt, start, end int64, workers int, debug bool) (segment, error) {
	length := max(end-start, 0)

	if debug {
		fmt.Printf("[DEBUG] parallelScan: range [%d,%d) length=%d workers=%d\n", start, end, length, workers)
	}

	if workers <= 0 || length == 0 {
		return segment{}, nil
	}

	if workers == 1 || length == 1 {
		if debug {
			fmt.Printf("[DEBUG] Sequential scan: range [%d,%d) length=%d\n", start, end, length)
		}
		buf := make([]byte, length)
		if length > 0 {
			if _, err := r.ReadAt(buf, start); err != nil && err != io.EOF {
				return segment{}, err
			}
		}
		sc := createScanner(string(buf))
		sg, err := sc.scan()
		if debug {
			fmt.Printf("[DEBUG] Sequential scan result: %d tokens\n", len(sg.Tokens))
		}
		return sg, err
	}

	leftWorkers := (workers + 1) / 2 // ceil(workers/2) ensures left >= right when odd
	rightWorkers := workers - leftWorkers

	// Proportional midpoint keeps byte distribution aligned with worker share.
	mid := start + (length*int64(leftWorkers))/int64(workers)

	// Adjust mid to rune boundary to avoid splitting multi-byte runes.
	adjustedPos, err := adjustToRuneBoundary(r, mid)
	if err != nil {
		return segment{}, err
	}

	mid = adjustedPos

	if debug {
		fmt.Printf("[DEBUG] Fork-join split: left[%d,%d) workers=%d, right[%d,%d) workers=%d\n",
			start, mid, leftWorkers, mid, end, rightWorkers)
	}

	type res struct {
		sg  segment
		err error
	}

	leftCh := make(chan res, 1)
	func() { // fork left branch
		if debug {
			fmt.Printf("[DEBUG] Forking left branch [%d,%d)\n", start, mid)
		}
		sg, err := parallelScan(r, start, mid, leftWorkers, debug)
		leftCh <- res{sg, err}
	}()

	// Recurse right branch in current goroutine (join pattern)
	if debug {
		fmt.Printf("[DEBUG] Processing right branch [%d,%d)\n", mid, end)
	}
	rightSeg, rightErr := parallelScan(r, mid, end, rightWorkers, debug)
	leftRes := <-leftCh

	if leftRes.err != nil {
		return segment{}, leftRes.err
	}
	if rightErr != nil {
		return segment{}, rightErr
	}

	if debug {
		fmt.Printf("[DEBUG] Merging segments: left=%d tokens, right=%d tokens\n",
			len(leftRes.sg.Tokens), len(rightSeg.Tokens))
	}

	leftRes.sg.Merge(&rightSeg)

	if debug {
		fmt.Printf("[DEBUG] Merge complete: total=%d tokens\n", len(leftRes.sg.Tokens))
	}

	return leftRes.sg, nil
}

func adjustToRuneBoundary(r io.ReaderAt, pos int64) (int64, error) {
	for i := pos; i > 0; i-- {
		readByte := make([]byte, 1)
		r.ReadAt(readByte, i)
		if utf8.RuneStart(readByte[0]) {
			return pos, nil
		}
	}

	return 0, nil
}

func (f *ForkyScanner) Scan(r io.ReaderAt, size int64) ([]common.Token, error) {
	if f.debug {
		fmt.Printf("[DEBUG] Starting scan with %d workers on %d bytes\n", f.numWorkers, size)
	}

	sg, err := parallelScan(r, 0, size, f.numWorkers, f.debug)
	if err != nil {
		return nil, err
	}
	if sg.hasInvalidTokens() {
		return nil, fmt.Errorf("merged segment has invalid tokens")
	}

	if f.debug {
		fmt.Printf("[DEBUG] Scan completed, found %d tokens\n", len(sg.Tokens))
	}

	return sg.Tokens, nil
}

// PUBLIC API

func CreateForkyScanner(numWorkers int, debug bool) ForkyScanner {
	if numWorkers < 1 {
		numWorkers = 1
	}
	return ForkyScanner{numWorkers: numWorkers, debug: debug}
}

func (f *ForkyScanner) ScanBytes(data []byte) ([]common.Token, error) {
	return f.Scan(bytesReader(data), int64(len(data)))
}

// scanString scans an in-memory string using the configured number of workers.
func (f *ForkyScanner) scanString(src string) ([]common.Token, error) {
	return f.ScanBytes([]byte(src))
}

// ScanString is a package-level helper to scan a string without manually creating a ForkJoinScanner.
func ScanString(src string, workers int) ([]common.Token, error) {
	sc := CreateForkyScanner(workers, false)
	return sc.scanString(src)
}
