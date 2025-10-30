package mail

import "context"

type IMailHandler interface {
	Send(ctx context.Context, msg MailSendMessage) error
}
