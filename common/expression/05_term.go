package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type TermNode struct {
	Left     Expression
	Operator common.Token
	Right    Expression
}

func (t *TermNode) Print(start string) {
	nodeName := fmt.Sprintf("Term (%s)", t.Operator.FriendlyOperatorName())
	fmt.Printf("%s%s\n", start, common.Colorize(nodeName, common.COLOR_MAGENTA))
	start = common.AdvanceSuffix(start)
	t.Left.Print(start + string(common.BRANCH_CONNECTOR))
	t.Right.Print(start + string(common.LAST_CONNECTOR))
}
