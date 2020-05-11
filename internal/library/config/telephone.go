// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import "github.com/axetroy/go-server/internal/service/dotenv"

type aliyunCloud struct {
	AccessKeyId               string `json:"access_key_id"`                // access key
	AccessSecret              string `json:"access_secret"`                // access secret
	SignName                  string `json:"sign_name"`                    // 短信签名名字
	TemplateCodeAuth          string `json:"template_code_auth"`           // 短信模版代码 - 身份验证
	TemplateCodeResetPassword string `json:"template_code_reset_password"` // 短信模版代码 - 重置密码
	TemplateCodeRegister      string `json:"template_code_register"`       // 短信模版代码 - 注册帐号
}

type tencentCloud struct {
	AppId                     string `json:"app_id"`                       // sdkappid请填写您在 短信控制台 添加应用后生成的实际 SDK AppID
	AppKey                    string `json:"app_key"`                      // sdkappid 对应的 appkey，需要业务方高度保密
	Sign                      string `json:"sign"`                         // 短信签名内容，使用 UTF-8 编码，必须填写已审核通过的签名。签名信息可登录 短信控制台 查看
	TemplateCodeAuth          string `json:"template_code_auth"`           // 短信模版代码 - 身份验证
	TemplateCodeResetPassword string `json:"template_code_reset_password"` // 短信模版代码 - 重置密码
	TemplateCodeRegister      string `json:"template_code_register"`       // 短信模版代码 - 注册帐号
}

type telephone struct {
	Provider string       `json:"provider"` // 选用哪家短信提供商
	Aliyun   aliyunCloud  `json:"aliyun"`   // 阿里云服务商相关配置
	Tencent  tencentCloud `json:"tencent"`  // 腾讯云服务商相关配置
}

var Telephone telephone

func init() {
	Telephone = telephone{
		Provider: dotenv.GetByDefault("TELEPHONE_PROVIDER", "aliyun"),
		Aliyun: aliyunCloud{
			AccessKeyId:               dotenv.Get("TELEPHONE_ALIYUN_ACCESS_KEY"),
			AccessSecret:              dotenv.Get("TELEPHONE_ALIYUN_ACCESS_SECRET"),
			SignName:                  dotenv.Get("TELEPHONE_ALIYUN_SIGN_NAME"),
			TemplateCodeAuth:          dotenv.Get("TELEPHONE_ALIYUN_TEMPLATE_CODE_AUTH"),
			TemplateCodeResetPassword: dotenv.Get("TELEPHONE_ALIYUN_TEMPLATE_CODE_RESET_PASSWORD"),
			TemplateCodeRegister:      dotenv.Get("TELEPHONE_ALIYUN_TEMPLATE_CODE_REGISTER"),
		},
		Tencent: tencentCloud{
			AppId:                     dotenv.Get("TELEPHONE_TENCENT_APP_ID"),
			AppKey:                    dotenv.Get("TELEPHONE_TENCENT_APP_KEY"),
			Sign:                      dotenv.Get("TELEPHONE_TENCENT_SIGN"),
			TemplateCodeAuth:          dotenv.Get("TELEPHONE_TENCENT_TEMPLATE_CODE_AUTH"),
			TemplateCodeResetPassword: dotenv.Get("TELEPHONE_TENCENT_TEMPLATE_CODE_RESET_PASSWORD"),
			TemplateCodeRegister:      dotenv.Get("TELEPHONE_TENCENT_TEMPLATE_CODE_REGISTER"),
		},
	}
}
