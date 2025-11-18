package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/peterh/liner"
)

const DEFAULT_WORKERS = 4

func main() {
	// Flags
	var (
		debug   bool
		modeStr string
		workers int
	)

	flag.BoolVar(&debug, "debug", false, "Enable debug output")
	flag.StringVar(&modeStr, "mode", "normal", "Run mode: normal, scanning, parsing")
	flag.IntVar(&workers, "workers", DEFAULT_WORKERS, "Number of workers for fork-join scanning")
	flag.Parse()

	// Determine mode based on string flag
	var mode InterpreterMode
	switch modeStr {
	case "scanning":
		mode = ScanningMode
	case "parsing":
		mode = ParsingMode
	case "normal":
		mode = NormalMode
	default:
		fmt.Printf("Invalid mode: %s. Valid modes are: normal, scanning, parsing\n", modeStr)
		os.Exit(1)
	}

	if workers <= 0 {
		workers = DEFAULT_WORKERS
	}

	forky := NewForky(workers, debug, mode)

	// If a file arg remains, run once on that file
	if flag.NArg() > 0 {
		path := flag.Arg(0)
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()
		st, err := f.Stat()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var runErr error
		_, runErr = forky.Run(f, st.Size())
		if runErr != nil {
			fmt.Println(runErr)
			os.Exit(1)
		}

		return
	}

	// REPL mode: read from stdin with arrow key support
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	fmt.Println("Forky - REPL with arrow key support. Ctrl-C or Ctrl-D (on empty line) to exit.")
	fmt.Println("Use ↑↓ arrows for history, ←→ for line editing.")
	fmt.Println()

	for {
		input, err := line.Prompt("> ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				fmt.Println("Aborted")
				break
			} else if err == io.EOF {
				fmt.Println()
				break
			}
			if debug {
				fmt.Printf("read error: %v\n", err)
			}
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Add to history
		line.AppendHistory(input)

		result, err := forky.Run(strings.NewReader(input), int64(len(input)))
		if err != nil {
			fmt.Println(err)
		} else if result != "" {
			fmt.Println(result)
		}
	}

}
