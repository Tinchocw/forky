package block

import (
	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/statement"
)

type BlockStatement struct {
	Statements []statement.Statement
}

func (bs BlockStatement) Print(start string) {
	statement.PrintStatements(start, bs.Statements)
}

func (bs BlockStatement) Headline() string {
	return common.Colorize("Block Statement", common.COLOR_MAGENTA)
}
