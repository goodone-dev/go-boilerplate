package consumer

import (
	"context"

	"github.com/goodone-dev/go-boilerplate/internal/application/mail/handler/worker"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/messaging/rabbitmq"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/messaging/rabbitmq/direct"
)

type consumer struct {
	client      rabbitmq.Client
	mailHandler mail.MailHandler
}

func NewConsumer(rmqClient rabbitmq.Client, mailUsecase mail.MailUsecase) *consumer {
	return &consumer{
		client:      rmqClient,
		mailHandler: worker.NewMailHandler(mailUsecase),
	}
}

func (c *consumer) Consume(ctx context.Context) {
	mailConsumer := direct.NewConsumer(ctx, c.client, direct.ConsumerConfig{
		ExchangeName: config.RabbitMQConfig.DirectExchangeName,
		QueueName:    "mail.send.queue",
		RoutingKey:   "mail.send",
		DLXEnabled:   true,
	})

	err := mailConsumer.ConsumeJSON(ctx, c.mailHandler.Send, mail.MailSendMessage{})
	if err != nil {
		logger.Fatal(ctx, err, "❌ Failed to start email consumer").Write()
	}
}
