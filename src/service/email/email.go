// Copyright 2019 Axetroy. All rights reserved. MIT license.
package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/axetroy/go-server/src/config"
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
)

const (
	prefix                 = "[GOTEST]: "
	tmpActivation          = `<a href="javascript: void 0">点击这里激活</a>或使用激活码: %v`
	tmpForgotPassword      = `<a href="javascript: void 0">点击连接重置密码</a>或使用重置码: %v`
	tmpForgotTradePassword = `<a href="javascript: void 0">点击连接重置交易密码</a>或使用重置码: %v`
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

	var addr = Config.Host + ":" + Config.Port

	if err = msg.SendWithTLS(
		addr,
		*e.Auth, &tls.Config{
			ServerName:         Config.Host,
			InsecureSkipVerify: true,
		}); err != nil {
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
		HTML:    []byte(fmt.Sprintf(tmpActivation, code)),
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
		HTML:    []byte(fmt.Sprintf(tmpForgotPassword, code)),
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
		HTML:    []byte(fmt.Sprintf(tmpForgotTradePassword, code)),
	}); err != nil {
		return
	}

	return nil
}
