package main

import (
	"context"
	"fmt"
	"log"
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
		log.Fatalf("❌ Could not load config: %v", err)
	}

	// ========== Infrastructure Setup ==========
	tracerProvider := tracer.NewProvider(ctx)
	postgresConn := postgres.Open()
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
		log.Printf("🚀 Starting server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Could not to start server: %v", err)
		}
	}()

	// ========== Graceful Shutdown ==========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println()
	log.Println("💤 Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server shutdown gracefully.")

	utils.GracefulShutdown(ctx, postgresConn, redisClient, tracerProvider)
}
