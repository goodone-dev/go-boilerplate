package mail

import "context"

type MailHandler interface {
	Send(ctx context.Context, payload any, headers map[string]any) error
}
