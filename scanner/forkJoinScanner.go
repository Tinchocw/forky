package scanner

import (
	"fmt"
	"io"
	"os"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type ForkJoinScanner struct{ numWorkers int }

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
		return sc.scan()
	}

	// 3. workers >= 2: proportional split based on worker allocation.
	leftWorkers := (workers + 1) / 2 // ceil(workers/2) ensures left >= right when odd
	rightWorkers := workers - leftWorkers

	// Proportional midpoint keeps byte distribution aligned with worker share.
	mid := start + (length*int64(leftWorkers))/int64(workers)
	// Ensure forward progress if integer division collapses interval.
	if mid <= start && end-start > 1 {
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
	return leftRes.sg, nil
}

func (f *ForkJoinScanner) scan(r io.ReaderAt, size int64) ([]common.Token, error) {
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

func CreateForkJoinScanner(numWorkers int) ForkJoinScanner {
	if numWorkers < 1 {
		numWorkers = 1
	}
	return ForkJoinScanner{numWorkers: numWorkers}
}

func (f *ForkJoinScanner) ScanFile(filePath string) ([]common.Token, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := info.Size()
	return f.scan(file, size)
}

func (f *ForkJoinScanner) ScanBytes(data []byte) ([]common.Token, error) {
	return f.scan(bytesReader(data), int64(len(data)))
}
