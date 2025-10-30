package messaging

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	"github.com/goodone-dev/go-boilerplate/internal/utils/validator"
)

type mailHandler struct {
	usecase mail.MailUsecase
}

func NewMailHandler(usecase mail.MailUsecase) mail.MailHandler {
	return &mailHandler{
		usecase: usecase,
	}
}

func (h *mailHandler) Send(ctx context.Context, msg mail.MailSendMessage) (err error) {
	ctx, span := tracer.Start(ctx, msg)
	defer func() {
		span.Stop(err)
	}()

	logger.Infof(ctx, "📧 Processing email send request to: %s", msg.To)

	if errs := validator.Validate(msg); errs != nil {
		logger.Errorf(ctx, errors.New(strings.Join(errs, ", ")), "❌ Failed to validate email send request to: %s", msg.To)
		return fmt.Errorf("request contains invalid or missing fields: %v", errs)
	}

	err = h.usecase.Send(ctx, msg)
	if err != nil {
		logger.Errorf(ctx, err, "❌ Failed to send email to: %s", msg.To)
		return
	}

	logger.Infof(ctx, "✅ Successfully sent email to: %s", msg.To)

	return
}
