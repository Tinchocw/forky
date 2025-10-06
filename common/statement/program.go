package statement

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Program struct {
	Statements []Statement
}

func (p *Program) Print() {
	printStatements("", p.Statements)
}

func PrintProgram(program Program) {
	fmt.Println()
	fmt.Printf("%s\n", common.Title("PROGRAM (Tree View)"))
	program.Print()
	fmt.Printf("%s\n", common.Title("End of PROGRAM"))
	fmt.Println()
}
