package mail

const MailSendTopic string = "mail.send"

type MailSendMessage struct {
	To       string `validate:"required,email"`
	Subject  string `validate:"required"`
	Template string `validate:"required"`
	Data     map[string]any
}
