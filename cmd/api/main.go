package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	customerrepo "github.com/BagusAK95/go-skeleton/internal/application/customer/repository"
	mailuc "github.com/BagusAK95/go-skeleton/internal/application/mail/usecase"
	orderrepo "github.com/BagusAK95/go-skeleton/internal/application/order/repository"
	orderuc "github.com/BagusAK95/go-skeleton/internal/application/order/usecase"
	productrepo "github.com/BagusAK95/go-skeleton/internal/application/product/repository"
	"github.com/BagusAK95/go-skeleton/internal/config"
	"github.com/BagusAK95/go-skeleton/internal/domain/customer"
	"github.com/BagusAK95/go-skeleton/internal/domain/mail"
	"github.com/BagusAK95/go-skeleton/internal/domain/order"
	"github.com/BagusAK95/go-skeleton/internal/domain/product"
	"github.com/BagusAK95/go-skeleton/internal/infrastructure/cache/redis"
	"github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql/postgres"
	mailsender "github.com/BagusAK95/go-skeleton/internal/infrastructure/mail"
	"github.com/BagusAK95/go-skeleton/internal/infrastructure/message/bus"
	"github.com/BagusAK95/go-skeleton/internal/infrastructure/tracer/jaeger"
	buslistener "github.com/BagusAK95/go-skeleton/internal/presentation/messaging/bus"
	"github.com/BagusAK95/go-skeleton/internal/presentation/rest/router"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	// Load .env config
	err := config.Load()
	if err != nil {
		log.Fatalf("❌ Could not load config: %v", err)
	}

	// Open connection
	postgresConn := postgres.Open()
	redisClient := redis.NewClient(ctx)
	tracerProvider := jaeger.Start()

	// Setup mail server
	mailSender := mailsender.NewMailSender()
	mailBus := bus.NewBus[mail.MailSendMessage]()

	// Initialize repository
	customerBaseRepo := postgres.NewBaseRepo[customer.Customer](postgresConn)
	customerRepo := customerrepo.NewCustomerRepo(customerBaseRepo)
	productBaseRepo := postgres.NewBaseRepo[product.Product](postgresConn)
	productRepo := productrepo.NewProductRepo(productBaseRepo)
	orderBaseRepo := postgres.NewBaseRepo[order.Order](postgresConn)
	orderRepo := orderrepo.NewOrderRepo(orderBaseRepo)
	orderItemBaseRepo := postgres.NewBaseRepo[order.OrderItem](postgresConn)
	orderItemRepo := orderrepo.NewOrderItemRepo(orderItemBaseRepo)

	// Initialize usecase
	mailUsecase := mailuc.NewMailUsecase(mailSender)
	orderUsecase := orderuc.NewOrderUsecase(
		customerRepo,
		productRepo,
		orderRepo,
		orderItemRepo,
		mailBus,
	)

	// Initialize bus listener
	buslistener.NewBusListener(mailBus, mailUsecase)

	// Start server
	gin.SetMode(gin.ReleaseMode)

	r := router.NewRouter(orderUsecase, redisClient)
	addr := fmt.Sprintf(":%d", config.ApplicationConfig.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Starting server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Could not to start server: %v", err)
		}
	}()

	// Gracefull Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println()
	log.Println("💤 Shutting down server...")

	// Close connection
	postgresConn.Close()
	redisClient.Close()
	tracerProvider.Shutdown(ctx)

	// Gracefull shutdown timeout
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exiting")
}
