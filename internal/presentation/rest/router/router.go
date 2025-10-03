package router

import (
	"log"
	"time"

	"net/http/pprof"

	orderHandler "github.com/BagusAK95/go-skeleton/internal/application/order/delivery/http"
	"github.com/BagusAK95/go-skeleton/internal/config"
	"github.com/BagusAK95/go-skeleton/internal/domain/order"
	"github.com/BagusAK95/go-skeleton/internal/infrastructure/cache"
	"github.com/BagusAK95/go-skeleton/internal/presentation/rest/middleware"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(orderUsecase order.IOrderUsecase, cacheClient cache.ICache) *gin.Engine {
	router := gin.New()
	router.Use(otelgin.Middleware(""))
	router.Use(middleware.ErrorMiddleware())
	router.Use(middleware.ContextMiddleware())
	router.Use(gin.Recovery())

	// Initialize handlers
	orderHandler := orderHandler.NewOrderHandler(orderUsecase)

	// Define routes
	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/order")
		{
			orders.POST(
				"",
				middleware.RateLimitMiddleware(cacheClient, 1, 1*time.Second),
				middleware.IdempotencyMiddleware(cacheClient, 5*time.Minute),
				orderHandler.Create,
			)
		}
	}

	if config.ApplicationConfig.Env == "production" {
		return router
	}

	// Enabling pprof for profiling
	log.Printf("🔎 Enabling pprof for profiling")
	debug := router.Group("/debug/pprof")
	{
		debug.GET("/goroutine", gin.WrapF(pprof.Index))
		debug.GET("/profile", gin.WrapF(pprof.Profile))
		debug.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		debug.GET("/symbol", gin.WrapF(pprof.Symbol))
		debug.GET("/trace", gin.WrapF(pprof.Trace))
	}

	return router
}
