package messaging

import (
	"context"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
)

type mailHandler struct {
	usecase mail.IMailUsecase
}

func NewMailHandler(usecase mail.IMailUsecase) *mailHandler {
	return &mailHandler{
		usecase: usecase,
	}
}

func (h *mailHandler) Send(ctx context.Context, msg mail.MailSendMessage) (err error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	logger.Infof(ctx, "processing email send request to: %s", msg.To)

	err = h.usecase.Send(ctx, msg)
	if err != nil {
		logger.Errorf(ctx, err, "failed to send email to: %s", msg.To)
		return
	}

	logger.Infof(ctx, "successfully sent email to: %s", msg.To)

	return
}
