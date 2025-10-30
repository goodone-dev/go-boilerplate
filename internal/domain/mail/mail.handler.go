package mail

import "context"

type MailHandler interface {
	Send(ctx context.Context, msg MailSendMessage) error
}
