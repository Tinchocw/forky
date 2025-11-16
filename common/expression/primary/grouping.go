package primaryExpression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type GroupingExpression struct {
	Expression          *expression.Expression
	StartingParenthesis bool
	ClosingParenthesis  bool
}

func (ge GroupingExpression) Print(start string) {
	fmt.Printf("%s%s\n", start, common.Colorize("GroupingExpression", common.COLOR_GREEN))
	start = common.AdvanceSuffix(start)
	ge.Expression.Print(start + string(common.LAST_CONNECTOR))
}

func (ge *GroupingExpression) HasLeft() bool {
	return ge.Expression != nil
}

func (ge *GroupingExpression) HasRight() bool {
	return false
}

func (ge *GroupingExpression) HasOperator() bool {
	if ge.ClosingParenthesis {
		return !ge.StartingParenthesis
	}

	return false
}

func (ge *GroupingExpression) GetLeft() expression.MergableNode {
	if ge.Expression == nil {
		panic("GetLeft: expected *Expression")
	}
	return ge.Expression

}

func (ge *GroupingExpression) GetRight() expression.MergableNode {
	panic("GetRight: GroupingExpression does not have a right node")
}

// (3 + 3 | )

func (ge *GroupingExpression) GetOperator() *common.Token {
	if ge.ClosingParenthesis {
		if ge.StartingParenthesis {
			panic("GetOperator: expected StartingParenthesis")
		}

		return &common.Token{Typ: common.CLOSE_PARENTHESIS}
	}

	panic("GetOperator: GroupingExpression does not have an operator")
}

func (ge *GroupingExpression) SetLeft(n expression.MergableNode) {
	expression, ok := n.(*expression.Expression)
	if !ok {
		panic("SetLeft: expected GroupingExpression")
	}

	ge.Expression = expression
}

func (ge *GroupingExpression) SetRight(n expression.MergableNode) {
	panic("SetRight: GroupingExpression does not have a right node")
}

func (ge *GroupingExpression) SetOperator(t *common.Token) {
	if t.Typ != common.CLOSE_PARENTHESIS {
		panic("SetOperator: expected CLOSE_PARENTHESIS")
	}

	if ge.StartingParenthesis {
		ge.ClosingParenthesis = true
	} else {
		panic("SetOperator: GroupingExpression does not have an operator")
	}
}

//  2 * ( 1 + 2 +  |  ( 2  | * 3 ) + 4 )

func (ge *GroupingExpression) ExpresionDepth(direction expression.Direction) int {
	if direction == expression.Left {
		if ge.ClosingParenthesis {
			return 1 + ge.Expression.ExpresionDepth(direction)
		}
		return ge.Expression.ExpresionDepth(direction)
	} else {
		if ge.StartingParenthesis {
			return 1 + ge.Expression.ExpresionDepth(direction)
		}
		return ge.Expression.ExpresionDepth(direction)
	}
}

func (ge *GroupingExpression) GetSubExpression(level int, direction expression.Direction) expression.MergableNode {
	if direction == expression.Left {
		if !ge.ClosingParenthesis {
			panic("Invalid state: expected ClosingParenthesis")
		}
	} else {
		if !ge.StartingParenthesis {
			panic("Invalid state: expected StartingParenthesis")
		}
	}

	return ge.Expression.GetSubExpression(level-1, direction)
}
