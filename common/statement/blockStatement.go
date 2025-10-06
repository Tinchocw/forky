package statement

type BlockStatement struct {
	Statements []Statement
}

func (bs BlockStatement) Print(start string) {
	printStatements(start, bs.Statements)
}
