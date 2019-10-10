// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import "github.com/axetroy/go-server/src/service/dotenv"

type aliyun struct {
	AccessKeyId               string `json:"access_key_id"`                // access key
	AccessSecret              string `json:"access_secret"`                // access secret
	SignName                  string `json:"sign_name"`                    // 短信签名名字
	TemplateCodeAuth          string `json:"template_code_auth"`           // 短信模版代码 - 身份验证
	TemplateCodeResetPassword string `json:"template_code_reset_password"` // 短信模版代码 - 重置密码
	TemplateCodeRegister      string `json:"template_code_register"`       // 短信模版代码 - 注册帐号
}

type telephone struct {
	Provider string `json:"provider"` // 选用哪家短信提供商
	Aliyun   aliyun `json:"aliyun"`   // 阿里云服务商相关配置
}

var Telephone telephone

func init() {
	Telephone = telephone{
		Provider: dotenv.GetByDefault("TELEPHONE_PROVIDER", "aliyun"),
		Aliyun: aliyun{
			AccessKeyId:               dotenv.Get("TELEPHONE_ALIYUN_ACCESS_KEY"),
			AccessSecret:              dotenv.Get("TELEPHONE_ALIYUN_ACCESS_SECRET"),
			SignName:                  dotenv.Get("TELEPHONE_ALIYUN_SIGN_NAME"),
			TemplateCodeAuth:          dotenv.Get("TELEPHONE_ALIYUN_TEMPLATE_CODE_AUTH"),
			TemplateCodeResetPassword: dotenv.Get("TELEPHONE_ALIYUN_TEMPLATE_CODE_RESET_PASSWORD"),
			TemplateCodeRegister:      dotenv.Get("TELEPHONE_ALIYUN_TEMPLATE_CODE_REGISTER"),
		},
	}
}
