package main

import (
	"fmt"
	"os"
	"strconv"

	scanpkg "github.com/Tinchocw/Interprete-concurrente/scanner"
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

func run(filePath string, workers int) error {
	scanner := scanpkg.CreateForkJoinScanner(workers)

	tokens, err := scanner.ScanFile(filePath)
	if err != nil {
		return err
	}

	for _, t := range tokens {
		fmt.Println(t.String())
	}
	// TODO: Parse & execute
	return nil
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
