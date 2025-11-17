package expression

import (
	"github.com/Tinchocw/Interprete-concurrente/common"
)

type TokenLiteralNode struct {
	Token *common.Token
}

func (tl TokenLiteralNode) Print(start string) {
	if tl.Token != nil {
		println(start, tl.Token.String())
	} else {
		println(start, "<nil>")
	}
}
