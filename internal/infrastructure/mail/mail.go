package mail

import (
	"bytes"
	"context"
	"crypto/tls"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	"github.com/goodone-dev/go-boilerplate/internal/utils/html"
	"gopkg.in/gomail.v2"
)

type MailSender interface {
	SendEmail(ctx context.Context, to, subject, file string, data any) error
}

type mailSender struct {
	dialer *gomail.Dialer
}

func NewMailSender() MailSender {
	d := gomail.NewDialer(config.MailConfig.Host, config.MailConfig.Port, config.MailConfig.Username, config.MailConfig.Password)

	if config.MailConfig.TLS {
		d.TLSConfig = &tls.Config{
			ServerName: config.MailConfig.Host,
			MinVersion: tls.VersionTLS13,
		}
	}

	return &mailSender{
		dialer: d,
	}
}

func (s *mailSender) SendEmail(ctx context.Context, to, subject, file string, data any) (err error) {
	_, span := tracer.Start(ctx, to, subject, file, data)
	defer func() {
		span.Stop(err)
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

	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
