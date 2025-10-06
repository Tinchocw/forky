package expression

import "github.com/Tinchocw/Interprete-concurrente/common"

type Direction int

const (
	Left Direction = iota + 1
	Right
)

type MergableNode interface {
	HasLeft() bool
	HasRight() bool
	HasOperator() bool

	GetLeft() MergableNode
	GetRight() MergableNode
	GetOperator() *common.Token

	SetLeft(n MergableNode)
	SetRight(n MergableNode)
	SetOperator(t *common.Token)

	ExpresionDepth(direction Direction) int
	GetSubExpression(level int, direction Direction) MergableNode
}
