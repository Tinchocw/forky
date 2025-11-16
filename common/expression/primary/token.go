package primaryExpression

import (
	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type TokenValue struct {
	Token *common.Token
}

func (t TokenValue) Print(start string) {
	if t.Token != nil {
		println(start, t.Token.String())
	} else {
		println(start, "<nil>")
	}
}

func (t TokenValue) HasLeft() bool {
	return t.Token != nil
}

func (t TokenValue) HasRight() bool {
	return false
}

func (t TokenValue) HasOperator() bool {
	return false
}

func (t TokenValue) GetLeft() expression.MergableNode {
	if t.Token == nil {
		panic("GetLeft: expected *Token")
	}
	return t
}

func (t TokenValue) GetRight() expression.MergableNode {
	panic("GetRight: TokenValue does not have a right node")
}

func (t TokenValue) GetOperator() *common.Token {
	panic("GetOperator: TokenValue does not have an operator")
}

func (t TokenValue) SetLeft(n expression.MergableNode) {
	panic("SetLeft: TokenValue does not allow setting a left node")
}

func (t TokenValue) SetRight(n expression.MergableNode) {
	panic("SetRight: TokenValue does not have a right node")
}

func (t TokenValue) SetOperator(op *common.Token) {
	panic("SetOperator: TokenValue does not have an operator")
}

func (t TokenValue) ExpresionDepth(direction expression.Direction) int {
	return 0
}

func (t TokenValue) GetSubExpression(levell int, direction expression.Direction) expression.MergableNode {
	panic("GetSubExpression: TokenValue does not have a sub-expression")
}
