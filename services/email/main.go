package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/axetroy/go-server/env"
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"os"
)

type from struct {
	Username string
	Email    string
}

type Config struct {
	From     from
	Username string
	Password string
	Server   string
	Port     string
}

var config Config

const (
	Prefix        = "[GOTEST]: "
	TmpActivation = `
<a href="javascript: void 0">点击这里激活</a>或使用激活码: %v
`
	TmpForgotPassword = `
<a href="javascript: void 0">点击连接重置密码</a>或使用重置码: %v
`
)

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

func init() {
	if err := env.Load(); err != nil {
		panic(err)
	}

	config = Config{
		From: from{
			Username: os.Getenv("SMTP_FROM_NAME"),
			Email:    os.Getenv("SMTP_FROM_EMAIL"),
		},
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Server:   os.Getenv("SMTP_SERVER"),
		Port:     os.Getenv("SMTP_SERVER_PORT"),
	}
}

func New() *Mailer {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		config.Username,
		config.Password,
		config.Server,
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
		From:    fmt.Sprintf("%v <%v>", config.From.Username, config.From.Email),
		To:      message.To,
		Subject: message.Subject,
		Text:    message.Text,
		HTML:    message.HTML,
		Headers: textproto.MIMEHeader{},
	}

	var addr = config.Server + ":" + config.Port

	if err = msg.SendWithTLS(
		addr,
		*e.Auth, &tls.Config{
			ServerName:         config.Server,
			InsecureSkipVerify: true,
		}); err != nil {
		return
	}

	return nil
}

// 发送激活邮件
func (e *Mailer) SendActivationEmail(to string, code string) (err error) {
	if err = e.Send(&Message{
		To:      []string{to},
		Subject: Prefix + "账号激活",
		Text:    []byte("请点击连接激活您的账号"),
		HTML:    []byte(fmt.Sprintf(TmpActivation, code)),
	}); err != nil {
		return
	}

	return nil
}

// 发送忘记密码邮件
func (e *Mailer) SendForgotPasswordEmail(to string, code string) (err error) {
	if err = e.Send(&Message{
		To:      []string{to},
		Subject: Prefix + "忘记密码",
		Text:    []byte("请点击连接激活您的账号"),
		HTML:    []byte(fmt.Sprintf(TmpForgotPassword, code)),
	}); err != nil {
		return
	}

	return nil
}
