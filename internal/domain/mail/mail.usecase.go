package mail

import "context"

type MailUsecase interface {
	Send(ctx context.Context, req MailSendMessage) error
}
