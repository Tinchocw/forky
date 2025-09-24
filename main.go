package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
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
		if err := forky.Run(f, st.Size()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return
	}

	// REPL mode: read from stdin line-by-line
	in := bufio.NewReader(os.Stdin)
	fmt.Println("Forky - REPL. Ctrl-D/Ctrl-C to exit.")
	fmt.Println()
	for {
		fmt.Print("> ")
		line, err := in.ReadString('\n')
		if err == io.EOF {
			fmt.Println()
			break
		}
		if err != nil {
			if debug {
				fmt.Printf("read error: %v\n", err)
			}
			continue
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			continue
		}
		if err := forky.Run(strings.NewReader(line), int64(len(line))); err != nil {
			fmt.Println(err)
		}
	}

}
