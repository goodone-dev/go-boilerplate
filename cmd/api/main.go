package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/goodone-dev/go-boilerplate/cmd/utils"
	customerrepo "github.com/goodone-dev/go-boilerplate/internal/application/customer/repository"
	healthhandler "github.com/goodone-dev/go-boilerplate/internal/application/health/delivery/http"
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
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	// ========== Environment Setup ==========
	err := config.Load()
	if err != nil {
		panic(err)
	}

	// ========== Observability Setup ==========
	loggerProvider := logger.NewProvider(ctx)
	tracerProvider := tracer.NewProvider(ctx)

	// ========== Infrastructure Setup ==========
	postgresConn := postgres.Open(ctx)
	redisClient := redis.NewClient(ctx)
	mailSender := mailsender.NewMailSender()

	// ========== Repositories Setup ==========
	customerBaseRepo := postgres.NewBaseRepository[gorm.DB, uuid.UUID, customer.Customer](postgresConn)
	customerRepo := customerrepo.NewCustomerRepository(customerBaseRepo)
	productBaseRepo := postgres.NewBaseRepository[gorm.DB, uuid.UUID, product.Product](postgresConn)
	productRepo := productrepo.NewProductRepository(productBaseRepo)
	orderBaseRepo := postgres.NewBaseRepository[gorm.DB, uuid.UUID, order.Order](postgresConn)
	orderRepo := orderrepo.NewOrderRepository(orderBaseRepo)
	orderItemBaseRepo := postgres.NewBaseRepository[gorm.DB, uuid.UUID, order.OrderItem](postgresConn)
	orderItemRepo := orderrepo.NewOrderItemRepository(orderItemBaseRepo)

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
		ReadTimeout:       config.HttpServerConfig.ReadTimeout,
		ReadHeaderTimeout: config.HttpServerConfig.ReadHeaderTimeout,
		WriteTimeout:      config.HttpServerConfig.WriteTimeout,
		IdleTimeout:       config.HttpServerConfig.IdleTimeout,
	}

	go func() {
		logger.Infof(ctx, "🚀 Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, err, "❌ Failed to start server")
		}
	}()

	// ========== Graceful Shutdown ==========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println()
	logger.Info(ctx, "🛑 Initiating server shutdown...")
	logger.Info(ctx, "⏳ Waiting for in-flight requests to complete...")

	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(ctx, err, "❌ Server forced to shutdown due to error")
	}

	logger.Info(ctx, "✅ Server shutdown gracefully")

	utils.GracefulShutdown(ctx, loggerProvider, tracerProvider, postgresConn, redisClient, mailBus)
}
