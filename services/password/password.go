package password

import (
	"github.com/axetroy/redpack/utils"
)

const (
	Prefix = "gotest"
)

func Generate(text string) string {
	return utils.MD5(Prefix + text)
}

// 是否是合法的用户密码
func IsValidUserPassword(text string) bool {
	return true
}

// 是否是合法的交易密码
func isValidPayPassword(text string) bool {
	return true
}
