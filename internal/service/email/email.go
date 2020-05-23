// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package email

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jordan-wright/email"
	"net"
	"net/smtp"
	"net/textproto"
)

// TODO: 重构邮箱模版
const (
	prefix                      = "[GOTEST]: "
	TemplateActivation          = `<a href="javascript: void 0">点击这里激活</a>或使用激活码: %v`
	TemplateForgotPassword      = `<a href="javascript: void 0">点击连接重置密码</a>或使用重置码: %v`
	TemplateForgotTradePassword = `<a href="javascript: void 0">点击连接重置交易密码</a>或使用重置码: %v`
	TemplateAuth                = `正在验证您的身份，你的验证码是 %s`
	TemplateRegistry            = `<a href="%s" href="target">点击注册您的帐号</a>`
)

type Mailer struct {
	Auth   *smtp.Auth
	Config model.ConfigFieldSMTP
}

type Message struct {
	ReplyTo     []string
	To          []string
	Bcc         []string
	Cc          []string
	Subject     string
	Text        []byte // Plaintext message (optional)
	HTML        []byte // Html message (optional)
	Sender      string // override From as SMTP envelope sender (optional)
	Headers     textproto.MIMEHeader
	Attachments []*email.Attachment
	ReadReceipt []string
}

// 从数据库中获取邮箱设置
func GetMailerConfig() (*model.ConfigFieldSMTP, error) {
	c := model.Config{}
	if err := database.Db.Model(&model.Config{}).Where(&c).First(&c).Error; err != nil {
		return nil, exception.NoConfig
	}

	if err := c.IsValidConfigField(); err != nil {
		return nil, err
	}

	result := model.ConfigFieldSMTP{}

	if err := json.Unmarshal([]byte(c.Fields), &result); err != nil {
		return nil, exception.InvalidParams.New(err.Error())
	}

	if err := validator.ValidateStruct(c); err != nil {
		return nil, err
	}

	return &result, nil
}

func NewMailer() (*Mailer, error) {
	// Set up authentication information.
	c, err := GetMailerConfig()

	if err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth(
		"",
		c.Username,
		c.Password,
		c.Server,
	)

	return &Mailer{
		Auth:   &auth,
		Config: *c,
	}, nil
}

// 发送邮件
func (e *Mailer) Send(message *Message) (err error) {
	if message == nil {
		err = errors.New("message can not be nil")
		return
	}
	msg := &email.Email{
		From:    fmt.Sprintf("%v <%v>", e.Config.FromName, e.Config.FromEmail),
		To:      message.To,
		Subject: message.Subject,
		Text:    message.Text,
		HTML:    message.HTML,
		Headers: textproto.MIMEHeader{},
	}

	var addr = net.JoinHostPort(e.Config.Server, fmt.Sprintf("%d", e.Config.Port))

	if err = msg.SendWithTLS(
		addr,
		*e.Auth, &tls.Config{
			ServerName:         e.Config.Server,
			InsecureSkipVerify: true,
		}); err != nil {
		err = exception.SendEmailFail
		return
	}

	return nil
}

// 发送激活邮件
func (e *Mailer) SendActivationEmail(toEmail string, code string) (err error) {
	if err = e.Send(&Message{
		To:      []string{toEmail},
		Subject: prefix + "账号激活",
		Text:    []byte("请点击连接激活您的账号"),
		HTML:    []byte(fmt.Sprintf(TemplateActivation, code)),
	}); err != nil {
		return
	}

	return nil
}

// 发送认证邮件
func (e *Mailer) SendAuthEmail(toEmail string, code string) (err error) {
	if err = e.Send(&Message{
		To:      []string{toEmail},
		Subject: prefix + "邮箱认证",
		Text:    []byte(fmt.Sprintf("您的验证码是: %s", code)),
		HTML:    []byte(fmt.Sprintf(TemplateAuth, code)),
	}); err != nil {
		return
	}

	return nil
}

// 发送注册邮件
func (e *Mailer) SendRegisterEmail(toEmail string, redirectURL string) (err error) {
	if err = e.Send(&Message{
		To:      []string{toEmail},
		Subject: prefix + "邮箱认证",
		Text:    []byte(fmt.Sprintf("打开链接注册帐号: %s", redirectURL)),
		HTML:    []byte(fmt.Sprintf(TemplateRegistry, redirectURL)),
	}); err != nil {
		return
	}

	return nil
}

// 发送忘记密码邮件
func (e *Mailer) SendForgotPasswordEmail(toEmail string, code string) (err error) {
	if err = e.Send(&Message{
		To:      []string{toEmail},
		Subject: prefix + "忘记登陆密码",
		Text:    []byte("请点击重置您的登陆密码"),
		HTML:    []byte(fmt.Sprintf(TemplateForgotPassword, code)),
	}); err != nil {
		return
	}

	return nil
}

// 发送忘记交易密码邮件
func (e *Mailer) SendForgotTradePasswordEmail(toEmail string, code string) (err error) {
	if err = e.Send(&Message{
		To:      []string{toEmail},
		Subject: prefix + "忘记交易密码",
		Text:    []byte("请点击连接激活您的账号"),
		HTML:    []byte(fmt.Sprintf(TemplateForgotTradePassword, code)),
	}); err != nil {
		return
	}

	return nil
}
