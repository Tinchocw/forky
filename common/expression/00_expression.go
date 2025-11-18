package expression

type Expression interface {
	Print(start string)
}

type ExpressionNode struct {
	Root *LogicalOrNode
}

func (e *ExpressionNode) Print(start string) {
	e.Root.Print(start)
}
