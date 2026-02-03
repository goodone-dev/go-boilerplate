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
	d := gomail.NewDialer(config.Mail.Host, config.Mail.Port, config.Mail.Username, config.Mail.Password)

	if config.Mail.TLS {
		d.TLSConfig = &tls.Config{
			ServerName: config.Mail.Host,
			MinVersion: tls.VersionTLS13,
		}
	}

	return &mailSender{
		dialer: d,
	}
}

func (s *mailSender) SendEmail(ctx context.Context, to, subject, file string, data any) (err error) {
	_, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"to":      to,
		"subject": subject,
		"file":    file,
		"data":    data,
	})

	defer func() {
		span.End(err)
	}()

	var body bytes.Buffer
	if err := html.ExecuteTemplate(&body, file, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Mail.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
