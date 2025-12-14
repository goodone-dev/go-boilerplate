package worker

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
	mailUsecase mail.MailUsecase
}

func NewMailHandler(mailUsecase mail.MailUsecase) mail.MailHandler {
	return &mailHandler{
		mailUsecase: mailUsecase,
	}
}

func (h *mailHandler) Send(ctx context.Context, payload any, headers map[string]any) (err error) {
	ctx, span := tracer.Start(ctx, payload)
	defer func() {
		span.Stop(err)
	}()

	body := payload.(mail.MailSendMessage)

	if errs := validator.Validate(body); errs != nil {
		logger.Errorf(ctx, errors.New(strings.Join(errs, ", ")), "❌ Failed to validate email send request to: %s", body.To)
		return fmt.Errorf("request contains invalid or missing fields: %v", errs)
	}

	err = h.mailUsecase.Send(ctx, body)
	if err != nil {
		logger.Errorf(ctx, err, "❌ Failed to send email to: %s", body.To)
		return
	}

	return
}
