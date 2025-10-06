package primaryExpression

import (
	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type Primary struct {
	Value PrimaryValue
}

type PrimaryValue interface {
	expression.MergableNode
	Print(start string)
}

func (p *Primary) Print(start string) {
	p.Value.Print(start)
}

func (p *Primary) HasLeft() bool {
	return true
}

func (p *Primary) HasRight() bool {
	return false
}

func (p *Primary) HasOperator() bool {
	return false
}

func (p *Primary) GetLeft() expression.MergableNode {
	return p.Value
}

func (p *Primary) GetRight() expression.MergableNode {
	panic("GetRight: Primary does not have a right node")
}

func (p *Primary) GetOperator() *common.Token {
	panic("GetOperator: Primary does not have an operator")
}

func (p *Primary) SetLeft(n expression.MergableNode) {
	primaryValue, ok := n.(PrimaryValue)
	if !ok {
		panic("SetLeft: expected PrimaryValue")
	}

	p.Value = primaryValue
}

func (p *Primary) SetRight(n expression.MergableNode) {
	panic("SetRight: Primary does not have a right node")
}

func (p *Primary) SetOperator(t *common.Token) {
	panic("SetOperator: Primary does not have an operator")
}

func (p *Primary) ExpresionDepth(direction expression.Direction) int {
	return p.Value.ExpresionDepth(direction)
}

func (p *Primary) GetSubExpression(level int, direction expression.Direction) expression.MergableNode {
	return p.Value.GetSubExpression(level, direction)
}
