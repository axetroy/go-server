// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/axetroy/go-server/internal/config"
	"github.com/axetroy/go-server/internal/exception"
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

var Config = config.SMTP

type Mailer struct {
	Auth *smtp.Auth
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

func NewMailer() *Mailer {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		Config.Username,
		Config.Password,
		Config.Host,
	)
	return &Mailer{
		Auth: &auth,
	}
}

// 发送邮件
func (e *Mailer) Send(message *Message) (err error) {
	if message == nil {
		err = errors.New("message can not be nil")
		return
	}
	msg := &email.Email{
		From:    fmt.Sprintf("%v <%v>", Config.Sender.Name, Config.Sender.Email),
		To:      message.To,
		Subject: message.Subject,
		Text:    message.Text,
		HTML:    message.HTML,
		Headers: textproto.MIMEHeader{},
	}

	var addr = net.JoinHostPort(Config.Host, Config.Port)

	if err = msg.SendWithTLS(
		addr,
		*e.Auth, &tls.Config{
			ServerName:         Config.Host,
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
