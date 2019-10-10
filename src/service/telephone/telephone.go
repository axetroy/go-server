package telephone

import (
	"github.com/axetroy/go-server/src/validator"
	"log"
)

type TemplateID string

const (
	TemplateAuth          TemplateID = "1" // 身份验证的模版
	TemplateResetPassword TemplateID = "2" // 重置密码的模版
)

// 发送短信验证码
func Send(phone string, templateID TemplateID, values ...interface{}) error {
	if err := validator.ValidatePhone(phone); err != nil {
		return err
	}
	// TODO: 接入发送短信验证码
	log.Printf("发送短信验证码")
	return nil
}
