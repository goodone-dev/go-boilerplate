package bus

import (
	mailhandler "github.com/goodone-dev/go-boilerplate/internal/application/mail/delivery/messaging"
	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/bus"
	"github.com/goodone-dev/go-boilerplate/internal/presentation/messaging/middleware"
)

func NewBusListener(mailBus bus.Bus[mail.MailSendMessage], mailUsecase mail.IMailUsecase) {
	mailHandler := mailhandler.NewMailHandler(mailUsecase)

	mailBus.SubscribeAsync(mail.MailSendTopic, middleware.TracerMiddleware(mail.MailSendTopic, mailHandler.Send), false)
}
