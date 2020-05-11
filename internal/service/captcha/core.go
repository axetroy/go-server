package captcha

import "github.com/axetroy/go-server/internal/util"

// 该包生成各种码
// 1. 验证码
// 2. 重置码

// 生成邮箱验证码
func GenerateEmailCaptcha() string {
	return util.RandomNumeric(6)
}

// 生成短信验证码
func GeneratePhoneCaptcha() string {
	return util.RandomNumeric(6)
}

// 生成密码重置码
func GenerateResetCode(uid string) string {
	codeId := util.GenerateId() + uid
	return util.MD5(codeId)
}
