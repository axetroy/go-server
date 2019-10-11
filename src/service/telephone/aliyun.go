package telephone

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/axetroy/go-server/src/config"
	"github.com/axetroy/go-server/src/exception"
)

func NewAliyun() *Aliyun {
	c := &Aliyun{}

	return c
}

type Aliyun struct {
}

func (c *Aliyun) getAuthTemplateID() string {
	return config.Telephone.Aliyun.TemplateCodeAuth
}

func (c *Aliyun) getResetPasswordTemplateID() string {
	return config.Telephone.Aliyun.TemplateCodeResetPassword
}

func (c *Aliyun) getRegisterTemplateID() string {
	return config.Telephone.Aliyun.TemplateCodeRegister
}

func (c *Aliyun) send(phone string, templateID string, templateMap map[string]string) error {
	aliClient, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", config.Telephone.Aliyun.AccessKeyId, config.Telephone.Aliyun.AccessSecret)

	if err != nil {
		return err
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phone
	request.SignName = config.Telephone.Aliyun.SignName
	request.TemplateCode = templateID

	b, err := json.Marshal(templateMap)

	if err != nil {
		return err
	}

	request.TemplateParam = string(b)

	res, err := aliClient.SendSms(request)

	if err != nil || !res.IsSuccess() {
		return exception.SendMsgFail
	}

	return nil
}

func (c *Aliyun) SendAuthCode(phone string, code string) error {
	return c.send(phone, c.getAuthTemplateID(), map[string]string{
		"code": code,
	})
}

func (c *Aliyun) SendResetPasswordCode(phone string, code string) error {
	return c.send(phone, c.getResetPasswordTemplateID(), map[string]string{
		"code": code,
	})
}

func (c *Aliyun) SendRegisterCode(phone string, code string) error {
	return c.send(phone, c.getRegisterTemplateID(), map[string]string{
		"code": code,
	})
}
