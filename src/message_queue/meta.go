package message_queue

type Topic string
type Chanel string

var (
	TopicSendEmail  Topic  = "send_email"
	ChanelSendEmail Chanel = "send_email"
)

type SendEmailBody struct {
	Email string `json:"email"` // 要发送的邮箱
	Code  string `json:"code"`  // 发送的激活码
}
