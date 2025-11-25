package expression

import (
	"github.com/Tinchocw/forky/common"
)

type TokenLiteralNode struct {
	Token common.Token
}

func (tl TokenLiteralNode) Print(start string) {
	println(start, tl.Token.String())
}
