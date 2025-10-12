package usecase

import (
	"context"

	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	mailsender "github.com/goodone-dev/go-boilerplate/internal/infrastructure/mail"
	"github.com/goodone-dev/go-boilerplate/internal/utils/tracer"
)

type MailUsecase struct {
	mailSender mailsender.IMailSender
}

func NewMailUsecase(mailSender mailsender.IMailSender) mail.IMailUsecase {
	return &MailUsecase{
		mailSender: mailSender,
	}
}

func (u *MailUsecase) Send(ctx context.Context, req mail.MailSendMessage) (err error) {
	ctx, span := tracer.StartSpan(ctx, req)
	defer func() {
		span.EndSpan(err)
	}()

	err = u.mailSender.SendEmail(ctx, req.To, req.Subject, req.Template, req.Data)
	if err != nil {
		return err
	}

	return nil
}
