package mail

const MailSendTopic string = "mail.send"

type MailSendMessage struct {
	To       string
	Subject  string
	Template string
	Data     map[string]any
}
