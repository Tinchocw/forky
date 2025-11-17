package primaryExpression

import (
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type Primary struct {
	Value PrimaryValue
}

type PrimaryValue interface {
	expression.ExpressionNode
}

func (p *Primary) Print(start string) {
	p.Value.Print(start)
}
