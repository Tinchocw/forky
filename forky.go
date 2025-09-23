package main

import (
	"fmt"
	"io"

	"github.com/Tinchocw/Interprete-concurrente/common"
	parserPackage "github.com/Tinchocw/Interprete-concurrente/parser"
	scannerPackage "github.com/Tinchocw/Interprete-concurrente/scanner"
)

type InterpreterMode int

const (
	NormalMode InterpreterMode = iota
	ScanningMode
	ParsingMode
	ColorDemoMode
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

	// Demo de colores si se solicita
	if fj.mode == ColorDemoMode {
		fj.showColorDemo()
		return nil
	}

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
		fmt.Printf("%s== Tokens ==%s\n", common.BoldCyan(""), common.COLOR_RESET)

		for i, t := range tokens {
			fmt.Printf("%s%4d:%s %s\n", common.Dim(""), i, common.COLOR_RESET, t.ColorString())
		}

		fmt.Printf("%s== End of Tokens ==%s\n", common.BoldCyan(""), common.COLOR_RESET)
		fmt.Println()
	}

	ps := parserPackage.NewParser(tokens)
	program, err := ps.Parse()
	if err != nil {
		return err
	}

	if fj.mode == ParsingMode {
		fmt.Println()
		fmt.Printf("%s== PROGRAM (Tree View) ==%s\n", common.BoldCyan(""), common.COLOR_RESET)
		program.Print(0)
		fmt.Printf("%s== End of PROGRAM ==%s\n", common.BoldCyan(""), common.COLOR_RESET)
		fmt.Println()
	}

	return nil
}

// showColorDemo muestra una demostración de las funciones de colores
func (fj *Forky) showColorDemo() {
	fmt.Printf("%s🎨 Demostración de funciones de colores 🎨%s\n\n", common.Rainbow(""), common.COLOR_RESET)

	// Colores básicos
	fmt.Printf("%s=== Colores básicos ===%s\n", common.BoldCyan(""), common.COLOR_RESET)
	fmt.Printf("%s\n", common.Red("Este texto es rojo"))
	fmt.Printf("%s\n", common.Green("Este texto es verde"))
	fmt.Printf("%s\n", common.Blue("Este texto es azul"))
	fmt.Printf("%s\n", common.Yellow("Este texto es amarillo"))
	fmt.Printf("%s\n", common.Magenta("Este texto es magenta"))
	fmt.Printf("%s\n", common.Cyan("Este texto es cian"))
	fmt.Println()

	// Estilos
	fmt.Printf("%s=== Estilos ===%s\n", common.BoldCyan(""), common.COLOR_RESET)
	fmt.Printf("%s\n", common.Bold("Este texto es negrita"))
	fmt.Printf("%s\n", common.Italic("Este texto es cursiva"))
	fmt.Printf("%s\n", common.Underline("Este texto está subrayado"))
	fmt.Printf("%s\n", common.Strike("Este texto está tachado"))
	fmt.Println()

	// Combinaciones
	fmt.Printf("%s=== Combinaciones ===%s\n", common.BoldCyan(""), common.COLOR_RESET)
	fmt.Printf("%s\n", common.BoldRed("Rojo en negrita"))
	fmt.Printf("%s\n", common.BoldGreen("Verde en negrita"))
	fmt.Printf("%s\n", common.UnderlineBlue("Azul subrayado"))
	fmt.Println()

	// Funciones divertidas
	fmt.Printf("%s=== Funciones especiales ===%s\n", common.BoldCyan(""), common.COLOR_RESET)
	fmt.Printf("%s\n", common.Rainbow("¡Texto arco iris!"))
	fmt.Printf("%s\n", common.Gradient("Gradiente", common.COLOR_RED, common.COLOR_BLUE))
	fmt.Println()

	// Box
	fmt.Printf("%s=== Cajas ===%s\n", common.BoldCyan(""), common.COLOR_RESET)
	fmt.Printf("%s\n", common.Box("¡Hola mundo!", common.COLOR_GREEN, common.BG_BLACK))
	fmt.Println()

	// Printf con colores
	fmt.Printf("%s=== Printf con colores ===%s\n", common.BoldCyan(""), common.COLOR_RESET)
	common.PrintfRed("Error: %s\n", "Algo salió mal")
	common.PrintfGreen("Éxito: %s\n", "Todo funcionó bien")
	common.PrintfYellow("Advertencia: %d elementos pendientes\n", 5)
	common.PrintfBlue("Info: Versión %s\n", "1.0.0")
	fmt.Println()

	// Uso manual con Colorize
	fmt.Printf("%s=== Colorize personalizado ===%s\n", common.BoldCyan(""), common.COLOR_RESET)
	fancy := common.Colorize("¡Súper fancy!", common.STYLE_BOLD, common.STYLE_UNDERLINE, common.COLOR_MAGENTA)
	fmt.Printf("%s\n", fancy)
}
