package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/axetroy/go-server/src/util"
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
	"os"
)

type emailFrom struct {
	Username string
	Email    string
}

type emailConfig struct {
	From     emailFrom
	Username string
	Password string
	Server   string
	Port     string
}

var config emailConfig

const (
	prefix                 = "[GOTEST]: "
	tmpActivation          = `<a href="javascript: void 0">点击这里激活</a>或使用激活码: %v`
	tmpForgotPassword      = `<a href="javascript: void 0">点击连接重置密码</a>或使用重置码: %v`
	tmpForgotTradePassword = `<a href="javascript: void 0">点击连接重置交易密码</a>或使用重置码: %v`
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
	if err := util.LoadEnv(); err != nil {
		panic(err)
	}

	config = emailConfig{
		From: emailFrom{
			Username: os.Getenv("SMTP_FROM_NAME"),
			Email:    os.Getenv("SMTP_FROM_EMAIL"),
		},
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Server:   os.Getenv("SMTP_SERVER"),
		Port:     os.Getenv("SMTP_SERVER_PORT"),
	}
}

func NewMailer() *Mailer {
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
