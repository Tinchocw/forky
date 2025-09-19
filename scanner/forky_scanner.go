package scanner

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

const windowSize = 4 // Number of runes to look ahead for merging segments

type ForkyScanner struct{ numWorkers int }

// parallelScan performs a fork-join recursive scan. It splits the range
// [start,end) roughly in half while there is budget (>1) of workers left
// and the segment length is > 1. Each recursive branch gets a share of the
// remaining worker budget. When budget == 1 or range is minimal, it scans
// sequentially in-place.
func parallelScan(r io.ReaderAt, start, end int64, workers int) (segment, error) {
	// 1. Invalid worker count
	if workers <= 0 {
		return segment{}, fmt.Errorf("workers must be >= 1 (got %d)", workers)
	}

	length := end - start
	if length < 0 {
		length = 0
	}

	// 2. Single worker (or empty length) -> direct scan
	if workers == 1 || length == 0 {
		buf := make([]byte, length)
		if length > 0 {
			if _, err := r.ReadAt(buf, start); err != nil && err != io.EOF {
				return segment{}, err
			}
		}
		sc := createScanner(string(buf))
		sg, err := sc.scan()
		// fmt.Println("Scanned segment: \n", sg.String())
		return sg, err
	}

	// 3. workers >= 2: proportional split based on worker allocation.
	leftWorkers := (workers + 1) / 2 // ceil(workers/2) ensures left >= right when odd
	rightWorkers := workers - leftWorkers

	// Proportional midpoint keeps byte distribution aligned with worker share.
	mid := start + (length*int64(leftWorkers))/int64(workers)

	fmt.Printf("Splitting range [%d,%d) at %d (L=%d R=%d)\n", start, end, mid, leftWorkers, rightWorkers)

	// Adjust mid to rune boundary to avoid splitting multi-byte runes.
	mid = int64(adjustToRuneBoundary(r, int(mid)))

	fmt.Printf("Adjusted mid to rune boundary: %d\n", mid)

	// Ensure forward progress if integer division collapses interval.
	if mid < start && end-start > 1 {
		mid = start + 1
	}
	if mid >= end && end-start > 1 {
		mid = end - 1
	}

	type res struct {
		sg  segment
		err error
	}
	leftCh := make(chan res, 1)
	go func() { // fork left branch
		sg, err := parallelScan(r, start, mid, leftWorkers)
		leftCh <- res{sg, err}
	}()

	// Recurse right branch in current goroutine (join pattern)
	rightSeg, rightErr := parallelScan(r, mid, end, rightWorkers)
	leftRes := <-leftCh

	if leftRes.err != nil {
		return segment{}, leftRes.err
	}
	if rightErr != nil {
		return segment{}, rightErr
	}

	leftRes.sg.Merge(&rightSeg)
	// fmt.Println("Merged segment: \n", leftRes.sg.String())
	return leftRes.sg, nil
}

func adjustToRuneBoundary(r io.ReaderAt, pos int) int {

	bufStart := pos - windowSize
	if bufStart < 0 {
		bufStart = 0
	}
	bufEnd := pos + windowSize
	if bufEnd < bufStart {
		bufEnd = bufStart
	}

	bufLen := bufEnd - bufStart
	data := make([]byte, bufLen)

	if _, err := r.ReadAt(data, int64(bufStart)); err != nil && err != io.EOF {
		// On read error, just return original pos (may split rune)
		return pos
	}

	fixPos := int(pos - bufStart)
	if fixPos >= len(data) {
		return bufLen
	}

	for !utf8.RuneStart(data[fixPos]) {
		fixPos--
		if fixPos <= 0 {
			return 0
		}
	}
	return bufStart + fixPos
}

func (f *ForkyScanner) Scan(r io.ReaderAt, size int64) ([]common.Token, error) {
	sg, err := parallelScan(r, 0, size, f.numWorkers)
	if err != nil {
		return nil, err
	}
	if sg.hasInvalidTokens() {
		return nil, fmt.Errorf("merged segment has invalid tokens")
	}
	return sg.Tokens, nil
}

// PUBLIC API

func CreateForkyScanner(numWorkers int) ForkyScanner {
	if numWorkers < 1 {
		numWorkers = 1
	}
	return ForkyScanner{numWorkers: numWorkers}
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
	sc := CreateForkyScanner(workers)
	return sc.scanString(src)
}
