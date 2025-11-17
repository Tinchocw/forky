package primaryExpression

import (
	"github.com/Tinchocw/Interprete-concurrente/common"
)

type TokenValue struct {
	Token *common.Token
}

func (t TokenValue) Print(start string) {
	if t.Token != nil {
		println(start, t.Token.String())
	} else {
		println(start, "<nil>")
	}
}
