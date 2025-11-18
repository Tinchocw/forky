package statement

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type Statement interface {
	Print(start string)
	Headline() string
}

func PrintStatements(start string, statements []Statement) {
	for i, stmt := range statements {
		fmt.Printf("%s%4d: %s\n", start, i+1, stmt.Headline())
		stmt.Print(common.AdvanceSuffix(start + string(common.COUNTER_INDENT)))
	}
}
