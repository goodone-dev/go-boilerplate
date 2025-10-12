package mail

import (
	"bytes"
	"context"
	"crypto/tls"

	"github.com/goodonedev/go-boilerplate/internal/config"
	"github.com/goodonedev/go-boilerplate/internal/utils/html"
	"github.com/goodonedev/go-boilerplate/internal/utils/tracer"
	"gopkg.in/gomail.v2"
)

type IMailSender interface {
	SendEmail(ctx context.Context, to, subject, file string, data any) error
}

type MailSender struct{}

func NewMailSender() IMailSender {
	return &MailSender{}
}

func (s *MailSender) SendEmail(ctx context.Context, to, subject, file string, data any) (err error) {
	_, span := tracer.StartSpan(ctx, to, subject, file, data)
	defer func() {
		span.EndSpan(err)
	}()

	var body bytes.Buffer
	if err := html.ExecuteTemplate(&body, file, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.MailConfig.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(config.MailConfig.Host, config.MailConfig.Port, config.MailConfig.Username, config.MailConfig.Password)

	if config.MailConfig.TLS {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
