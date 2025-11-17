package expression

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
)

type Term struct {
	Left     ExpressionNode
	Operator *common.Token
	Right    *Factor
}

func (t *Term) Print(start string) {
	if t.skipPrinting() {
		t.Left.Print(start)
		return
	}

	nodeName := fmt.Sprintf("Term (%s)", t.Operator.FriendlyOperatorName())

	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	t.Left.Print(start + string(common.BRANCH_CONNECTOR))
	t.Right.Print(start + string(common.LAST_CONNECTOR))
}

func (t *Term) skipPrinting() bool {
	return t.Operator == nil && t.Right == nil
}
