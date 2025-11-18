package statement

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type Program struct {
	Statements []Statement
}

func (p *Program) Print() {
	PrintStatements("", p.Statements)
}

func PrintProgram(program Program) {
	fmt.Println()
	fmt.Printf("%s\n", common.Title("PROGRAM (Tree View)"))
	program.Print()
	fmt.Printf("%s\n", common.Title("End of PROGRAM"))
	fmt.Println()
}
