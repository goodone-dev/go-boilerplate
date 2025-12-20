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
	ctx, span := tracer.Start(ctx, payload, headers)
	defer func() {
		span.Stop(err)
	}()

	req := payload.(mail.MailSendMessage)

	if errs := validator.Validate(req); errs != nil {
		logger.With().Error(ctx, errors.New(strings.Join(errs, ", ")), "❌ Failed to validate email send request")
		return fmt.Errorf("request contains invalid or missing fields: %v", errs)
	}

	err = h.mailUsecase.Send(ctx, req)
	if err != nil {
		logger.With().Error(ctx, err, "❌ Failed to send email")
		return
	}

	return
}
