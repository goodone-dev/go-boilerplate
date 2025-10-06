package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"syscall"

	customerrepo "github.com/BagusAK95/go-boilerplate/internal/application/customer/repository"
	mailuc "github.com/BagusAK95/go-boilerplate/internal/application/mail/usecase"
	orderrepo "github.com/BagusAK95/go-boilerplate/internal/application/order/repository"
	orderuc "github.com/BagusAK95/go-boilerplate/internal/application/order/usecase"
	productrepo "github.com/BagusAK95/go-boilerplate/internal/application/product/repository"
	"github.com/BagusAK95/go-boilerplate/internal/config"
	"github.com/BagusAK95/go-boilerplate/internal/domain/customer"
	"github.com/BagusAK95/go-boilerplate/internal/domain/mail"
	"github.com/BagusAK95/go-boilerplate/internal/domain/order"
	"github.com/BagusAK95/go-boilerplate/internal/domain/product"
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/cache/redis"
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql/postgres"
	mailsender "github.com/BagusAK95/go-boilerplate/internal/infrastructure/mail"
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/message/bus"
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/tracer/jaeger"
	buslistener "github.com/BagusAK95/go-boilerplate/internal/presentation/messaging/bus"
	"github.com/BagusAK95/go-boilerplate/internal/presentation/rest/router"
)

func main() {
	ctx := context.Background()

	// ========== Environment Setup ==========
	err := config.Load()
	if err != nil {
		log.Fatalf("❌ Could not load config: %v", err)
	}

	// ========== Infrastructure ==========
	postgresConn := postgres.Open()
	redisClient := redis.NewClient(ctx)
	tracerProvider := jaeger.NewProvider(ctx)
	mailSender := mailsender.NewMailSender()

	// ========== Repositories ==========
	customerBaseRepo := postgres.NewBaseRepo[customer.Customer](postgresConn)
	customerRepo := customerrepo.NewCustomerRepo(customerBaseRepo)
	productBaseRepo := postgres.NewBaseRepo[product.Product](postgresConn)
	productRepo := productrepo.NewProductRepo(productBaseRepo)
	orderBaseRepo := postgres.NewBaseRepo[order.Order](postgresConn)
	orderRepo := orderrepo.NewOrderRepo(orderBaseRepo)
	orderItemBaseRepo := postgres.NewBaseRepo[order.OrderItem](postgresConn)
	orderItemRepo := orderrepo.NewOrderItemRepo(orderItemBaseRepo)

	// ========== Bus ==========
	mailBus := bus.NewBus[mail.MailSendMessage]()

	// ========== Usecase ==========
	mailUsecase := mailuc.NewMailUsecase(mailSender)
	orderUsecase := orderuc.NewOrderUsecase(
		customerRepo,
		productRepo,
		orderRepo,
		orderItemRepo,
		mailBus,
	)

	// ========== Bus Listener ==========
	buslistener.NewBusListener(mailBus, mailUsecase)

	// ========== HTTP Server Setup ==========
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

	infraShutdown(ctx, postgresConn, redisClient, tracerProvider)
}

type Infrastructure interface {
	Shutdown(ctx context.Context) error
}

func infraShutdown(ctx context.Context, infras ...Infrastructure) {
	var wg sync.WaitGroup

	for _, infra := range infras {
		wg.Add(1)

		go func(c Infrastructure) {
			defer wg.Done()

			packageName := parsePackageName(c)

			if err := c.Shutdown(ctx); err != nil {
				log.Printf("❌ %s forced to shutdown: %v", packageName, err)
				return
			}

			log.Printf("✅ %s shutdown gracefully.\n", packageName)
		}(infra)
	}

	wg.Wait()
}

func parsePackageName(infra Infrastructure) string {
	n := reflect.TypeOf(infra).String()
	r := regexp.MustCompile(`\*?([^.]+)`)

	matches := r.FindStringSubmatch(n)
	if len(matches) > 1 {
		name := matches[1]
		return strings.ToUpper(string(name[0])) + name[1:]
	}

	return ""
}
