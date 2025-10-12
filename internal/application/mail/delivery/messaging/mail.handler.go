package messaging

import (
	"context"
	"log"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
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
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	log.Printf("✉️ Receiving mail.send message: %s", msg.To)

	err = h.usecase.Send(ctx, msg)
	if err != nil {
		log.Printf("❌ Could not to send email: %v", err)
		return
	}

	log.Printf("✅ Email sent to %s", msg.To)

	return
}
