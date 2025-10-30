package usecase

import (
	"context"

	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	mailsender "github.com/goodone-dev/go-boilerplate/internal/infrastructure/mail"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
)

type mailUsecase struct {
	mailSender mailsender.IMailSender
}

func NewMailUsecase(mailSender mailsender.IMailSender) mail.IMailUsecase {
	return &mailUsecase{
		mailSender: mailSender,
	}
}

func (u *mailUsecase) Send(ctx context.Context, req mail.MailSendMessage) (err error) {
	ctx, span := tracer.Start(ctx, req)
	defer func() {
		span.Stop(err)
	}()

	err = u.mailSender.SendEmail(ctx, req.To, req.Subject, req.Template, req.Data)
	if err != nil {
		return err
	}

	return nil
}
