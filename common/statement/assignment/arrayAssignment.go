package assignment

import (
	"fmt"

	"github.com/Tinchocw/forky/common"
	"github.com/Tinchocw/forky/common/expression"
)

type ArrayAssignment struct {
	Name    string
	Indexes []*expression.ExpressionNode
	Value   *expression.ExpressionNode
}

func (aa ArrayAssignment) Print(start string) {
	totalChildren := 1 // Name
	if len(aa.Indexes) > 0 {
		totalChildren++ // Indexes group
	}
	if aa.Value != nil {
		totalChildren++
	}

	childIndex := 0
	// advance suffix so vertical connectors propagate to child groups
	start = common.AdvanceSuffix(start)
	basePrefix := start + string(common.SIMPLE_INDENT)

	// Name
	nameConn := string(common.BRANCH_CONNECTOR)
	if totalChildren == 1 { // only name
		nameConn = string(common.LAST_CONNECTOR)
	}

	fmt.Printf("%s%s %s\n",
		basePrefix+nameConn,
		common.Colorize("Name:", common.COLOR_YELLOW),
		common.Colorize(aa.Name, common.COLOR_WHITE),
	)
	childIndex++

	// Indexes (as a group)
	if len(aa.Indexes) > 0 {
		isLast := (childIndex == totalChildren-1)
		indexesConn := string(common.BRANCH_CONNECTOR)
		if isLast {
			indexesConn = string(common.LAST_CONNECTOR)
		}

		fmt.Printf("%s%s\n",
			basePrefix+indexesConn,
			common.Colorize("Indexes:", common.COLOR_YELLOW),
		)

		// Individual indexes
		for i, indexExpr := range aa.Indexes {
			indexIsLast := (i == len(aa.Indexes)-1)
			indexConn := string(common.BRANCH_CONNECTOR)
			if indexIsLast {
				indexConn = string(common.LAST_CONNECTOR)
			}

			// print index label
			fmt.Printf("%s%s\n",
				basePrefix+string(common.SIMPLE_INDENT)+indexConn,
				common.Colorize(fmt.Sprintf("%d:", i), common.COLOR_CYAN),
			)

			// render the index expression
			exprPrefix := basePrefix + string(common.SIMPLE_INDENT) + string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
			if indexExpr != nil {
				(*indexExpr).Print(exprPrefix)
			}
		}
		childIndex++
	}

	// Value (always last if exists)
	if aa.Value != nil {
		fmt.Printf("%s%s\n",
			basePrefix+string(common.LAST_CONNECTOR),
			common.Colorize("Value:", common.COLOR_YELLOW),
		)
		valuePrefix := basePrefix + string(common.SIMPLE_INDENT) + string(common.LAST_CONNECTOR)
		(*aa.Value).Print(valuePrefix)
	}
}

func (aa ArrayAssignment) Headline() string {
	return common.Colorize("Array Assignment", common.COLOR_GREEN)
}
