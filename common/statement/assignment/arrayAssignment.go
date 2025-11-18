package assignment

import (
	"fmt"

	"github.com/Tinchocw/Interprete-concurrente/common"
	"github.com/Tinchocw/Interprete-concurrente/common/expression"
)

type ArrayAssignment struct {
	Name    string
	Indexes []*expression.ExpressionNode
	Value   *expression.ExpressionNode
}

func (aa ArrayAssignment) Print(start string) {

	// Contar hijos totales (Name + Indexes (como grupo) + Value)
	totalChildren := 1 // Name
	if len(aa.Indexes) > 0 {
		totalChildren++ // Indexes como grupo
	}
	if aa.Value != nil {
		totalChildren++
	}

	childIndex := 0
	basePrefix := start + string(common.SIMPLE_INDENT)

	// Name
	nameConn := string(common.BRANCH_CONNECTOR)
	if totalChildren == 1 { // solo name
		nameConn = string(common.LAST_CONNECTOR)
	}

	fmt.Printf("%s%s %s\n",
		basePrefix+nameConn,
		common.Colorize("Name:", common.COLOR_YELLOW),
		common.Colorize(aa.Name, common.COLOR_CYAN),
	)
	childIndex++

	// Indexes (como grupo padre)
	if len(aa.Indexes) > 0 {
		isLast := (childIndex == totalChildren-1)
		indexesConn := string(common.BRANCH_CONNECTOR)
		indexesBasePrefix := basePrefix + "│   "
		if isLast {
			indexesConn = string(common.LAST_CONNECTOR)
			indexesBasePrefix = basePrefix + "    "
		}

		fmt.Printf("%s%s\n",
			basePrefix+indexesConn,
			common.Colorize("Indexes:", common.COLOR_YELLOW),
		)

		// Índices individuales
		for i, indexExpr := range aa.Indexes {
			indexIsLast := (i == len(aa.Indexes)-1)
			indexConn := string(common.BRANCH_CONNECTOR)
			nextPrefix := indexesBasePrefix + "│   "
			if indexIsLast {
				indexConn = string(common.LAST_CONNECTOR)
				nextPrefix = indexesBasePrefix + "    "
			}

			fmt.Printf("%s%s\n",
				indexesBasePrefix+indexConn,
				common.Colorize(fmt.Sprintf("%d:", i), common.COLOR_CYAN),
			)

			// Renderizar el árbol de la expresión del índice
			if indexExpr != nil {
				(*indexExpr).Print(nextPrefix)
			}
		}
		childIndex++
	}

	// Value (siempre último si existe)
	if aa.Value != nil {
		fmt.Printf("%s%s\n",
			basePrefix+string(common.LAST_CONNECTOR),
			common.Colorize("Value:", common.COLOR_YELLOW),
		)
		// Renderizar el árbol de la expresión del valor
		valuePrefix := basePrefix + "    "
		(*aa.Value).Print(valuePrefix)
	}
}

func (aa ArrayAssignment) Headline() string {
	return common.Colorize("Array Assignment", common.COLOR_GREEN)
}
