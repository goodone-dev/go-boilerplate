package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goodone-dev/go-boilerplate/cmd/utils"
	customerrepo "github.com/goodone-dev/go-boilerplate/internal/application/customer/repository"
	healthhandler "github.com/goodone-dev/go-boilerplate/internal/application/health/http"
	mailuc "github.com/goodone-dev/go-boilerplate/internal/application/mail/usecase"
	orderhandler "github.com/goodone-dev/go-boilerplate/internal/application/order/delivery/http"
	orderrepo "github.com/goodone-dev/go-boilerplate/internal/application/order/repository"
	orderuc "github.com/goodone-dev/go-boilerplate/internal/application/order/usecase"
	productrepo "github.com/goodone-dev/go-boilerplate/internal/application/product/repository"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/domain/customer"
	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/domain/product"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/cache/redis"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database/postgres"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	mailsender "github.com/goodone-dev/go-boilerplate/internal/infrastructure/mail"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/bus"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	buslistener "github.com/goodone-dev/go-boilerplate/internal/presentation/messaging/bus"
	"github.com/goodone-dev/go-boilerplate/internal/presentation/rest/router"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	// ========== Environment Setup ==========
	err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	// ========== Observability Setup ==========
	loggerProvider := logger.NewProvider(ctx)
	tracerProvider := tracer.NewProvider(ctx)

	// ========== Infrastructure Setup ==========
	postgresConn := postgres.Open(ctx)
	redisClient := redis.NewClient(ctx)
	mailSender := mailsender.NewMailSender()

	// ========== Repositories Setup ==========
	customerBaseRepo := postgres.NewBaseRepo[gorm.DB, uuid.UUID, customer.Customer](postgresConn)
	customerRepo := customerrepo.NewCustomerRepo(customerBaseRepo)
	productBaseRepo := postgres.NewBaseRepo[gorm.DB, uuid.UUID, product.Product](postgresConn)
	productRepo := productrepo.NewProductRepo(productBaseRepo)
	orderBaseRepo := postgres.NewBaseRepo[gorm.DB, uuid.UUID, order.Order](postgresConn)
	orderRepo := orderrepo.NewOrderRepo(orderBaseRepo)
	orderItemBaseRepo := postgres.NewBaseRepo[gorm.DB, uuid.UUID, order.OrderItem](postgresConn)
	orderItemRepo := orderrepo.NewOrderItemRepo(orderItemBaseRepo)

	// ========== Bus Setup ==========
	mailBus := bus.NewBus[mail.MailSendMessage]()

	// ========== Usecase Setup ==========
	mailUsecase := mailuc.NewMailUsecase(mailSender)
	orderUsecase := orderuc.NewOrderUsecase(
		customerRepo,
		productRepo,
		orderRepo,
		orderItemRepo,
		mailBus,
	)

	// ========== HTTP Handler Setup ==========
	healthHandler := healthhandler.NewHealthHandler(postgresConn, redisClient)
	orderHandler := orderhandler.NewOrderHandler(orderUsecase)

	// ========== Bus Listener Setup ==========
	buslistener.NewBusListener(mailBus, mailUsecase)

	// ========== HTTP Server Setup ==========
	r := router.NewRouter(healthHandler, orderHandler, redisClient)
	addr := fmt.Sprintf(":%d", config.ApplicationConfig.Port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		logger.Infof(ctx, "starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, err, "failed to start server")
		}
	}()

	// ========== Graceful Shutdown ==========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println()
	logger.Info(ctx, "initiating server shutdown...")

	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(ctx, err, "server forced to shutdown due to error")
	}

	logger.Info(ctx, "server shutdown gracefully")

	utils.GracefulShutdown(ctx, loggerProvider, tracerProvider, postgresConn, redisClient)
}
