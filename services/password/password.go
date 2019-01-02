package password

import (
	"github.com/axetroy/go-server/utils"
)

const (
	prefix = "gotest"
	suffix = "test"
)

func Generate(text string) string {
	password := utils.MD5(prefix + text + suffix)
	return password
}

// 是否是合法的用户密码
func IsValidUserPassword(text string) bool {
	return true
}

// 是否是合法的交易密码
func isValidPayPassword(text string) bool {
	return true
}
