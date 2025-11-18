package expression

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
)

type TermNode struct {
	Left     Expression // TermNode or FactorNode
	Operator *common.Token
	Right    *FactorNode
}

func (t *TermNode) Print(start string) {
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

func (t *TermNode) skipPrinting() bool {
	return t.Operator == nil && t.Right == nil
}
