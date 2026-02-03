package usecase

import (
	"context"

	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	mailsender "github.com/goodone-dev/go-boilerplate/internal/infrastructure/mail"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
)

type mailUsecase struct {
	mailSender mailsender.MailSender
}

func NewMailUsecase(mailSender mailsender.MailSender) mail.MailUsecase {
	return &mailUsecase{
		mailSender: mailSender,
	}
}

func (u *mailUsecase) Send(ctx context.Context, req mail.MailSendMessage) (err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"request": req,
	})

	defer func() {
		span.End(err)
	}()

	err = u.mailSender.SendEmail(ctx, req.To, req.Subject, req.Template, req.Data)
	if err != nil {
		return err
	}

	return nil
}
