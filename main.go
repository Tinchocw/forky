package main

import (
	"fmt"
	"os"
	"strconv"
)

func readArguments() (string, int, error) {
	if len(os.Args) < 2 {
		exetuable := os.Args[0]
		return "", 0, fmt.Errorf("uso esperado: %s <file_path> <num_workers>", exetuable)
	}

	args := os.Args[1:]

	file_path := args[0]
	num_workers := args[1]

	num_workers_int, err := strconv.Atoi(num_workers)
	if err != nil {
		return "", 0, fmt.Errorf("el argumento num_workers debe ser un entero, se recibió: %s", num_workers)
	}

	return file_path, num_workers_int, nil
}

type WorkerResult struct {
	id              int
	couldMergeStart bool
	couldMergeEnd   bool
	tokens          []Token
	content         string
	err             error
}

// String implements fmt.Stringer to pretty-print a WorkerResult with indenting.
func (wr WorkerResult) String() string {
	out := fmt.Sprintf("Part %d\n", wr.id)
	out += fmt.Sprintf("  couldMergeStart: %v\n", wr.couldMergeStart)
	out += fmt.Sprintf("  couldMergeEnd:   %v\n", wr.couldMergeEnd)
	out += "  Tokens:\n"
	if len(wr.tokens) == 0 {
		out += "    <no tokens>\n"
		return out
	}
	for _, tok := range wr.tokens {
		out += fmt.Sprintf("    %s\n", tok.String())
	}
	return out
}

func worker(id, ini, fin int, file *os.File, c chan WorkerResult) {
	buf := make([]byte, fin-ini+1)

	_, err := file.ReadAt(buf, int64(ini))
	if err != nil {
		c <- WorkerResult{id: id, content: "", err: err}
		return
	}

	readContent := string(buf)

	scanner := createScanner(readContent)
	err = scanner.scan()

	if err != nil {
		c <- WorkerResult{id: id, err: err}
		return
	}

	result := WorkerResult{
		id:              id,
		couldMergeStart: scanner.canMergeStart,
		couldMergeEnd:   scanner.canMergeEnd,
		tokens:          scanner.tokens,
		content:         readContent,
		err:             nil,
	}

	c <- result
}

func scan(file *os.File, num_workers int) ([]Token, error) {
	file_size, err := getFileSize(file)
	if err != nil {
		return nil, err
	}

	part_size := (file_size + int64(num_workers) - 1) / int64(num_workers)
	result_chan := make(chan WorkerResult, num_workers)

	for i := range num_workers {
		ini := i * int(part_size)
		fin := (i+1)*int(part_size) - 1

		if fin >= int(file_size) {
			fin = int(file_size) - 1
		}

		go worker(i, ini, fin, file, result_chan)
	}

	parts := make([]WorkerResult, num_workers)

	for range num_workers {
		result := <-result_chan
		if result.err != nil {
			return nil, result.err
		}
		parts[result.id] = result
	}

	fmt.Println("Parts tokens (ordered by worker id):")
	for _, part := range parts {
		fmt.Print(part.String())
	}

	mergedTokens := []Token{}

	// TODO: Merge tokens considering couldMergeStart and couldMergeEnd

	return mergedTokens, nil
}

func run(file_path string, num_workers int) error {
	file, err := os.Open(file_path)
	if err != nil {
		return err
	}

	defer file.Close()

	tokens, err := scan(file, num_workers)
	if err != nil {
		return err
	}

	fmt.Println(tokens)

	// 2. Parse tokens
	// 3. Execute AST

	return nil
}

func getFileSize(file *os.File) (int64, error) {
	file_info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return file_info.Size(), nil
}

func main() {
	file_path, num_workers, err := readArguments()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("file_path: %s, num_workers: %d\n", file_path, num_workers)

	err = run(file_path, num_workers)
	if err != nil {
		fmt.Println(err)
		return
	}
}
