package bus

import (
	mailhandler "github.com/BagusAK95/go-skeleton/internal/application/mail/delivery/messaging"
	"github.com/BagusAK95/go-skeleton/internal/domain/mail"
	"github.com/BagusAK95/go-skeleton/internal/infrastructure/message/bus"
	"github.com/BagusAK95/go-skeleton/internal/presentation/messaging/middleware"
)

func NewBusListener(mailBus bus.Bus[mail.MailSendMessage], mailUsecase mail.IMailUsecase) {
	mailHandler := mailhandler.NewMailHandler(mailUsecase)

	mailBus.SubscribeAsync(mail.MailSendTopic, middleware.TracerMiddleware(mail.MailSendTopic, mailHandler.Send), false)
}
