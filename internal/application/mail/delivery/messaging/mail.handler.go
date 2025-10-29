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
	usecase mail.IMailUsecase
}

func NewMailHandler(usecase mail.IMailUsecase) *mailHandler {
	return &mailHandler{
		usecase: usecase,
	}
}

func (h *mailHandler) Send(ctx context.Context, msg mail.MailSendMessage) (err error) {
	ctx, span := tracer.Start(ctx, msg)
	defer func() {
		span.Stop(err)
	}()

	logger.Infof(ctx, "processing email send request to: %s", msg.To)

	if errs := validator.Validate(msg); errs != nil {
		logger.Errorf(ctx, errors.New(strings.Join(errs, ", ")), "failed to validate email send request to: %s", msg.To)
		return fmt.Errorf("request contains invalid or missing fields: %v", errs)
	}

	err = h.usecase.Send(ctx, msg)
	if err != nil {
		logger.Errorf(ctx, err, "failed to send email to: %s", msg.To)
		return
	}

	logger.Infof(ctx, "successfully sent email to: %s", msg.To)

	return
}
