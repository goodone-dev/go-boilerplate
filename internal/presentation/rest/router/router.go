package router

import (
	"log"
	"time"

	"net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	orderhandler "github.com/goodone-dev/go-boilerplate/internal/application/order/delivery/http"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/cache"
	"github.com/goodone-dev/go-boilerplate/internal/presentation/rest/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(orderUsecase order.IOrderUsecase, cacheClient cache.ICache) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	// Initialize router
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins: config.CorsAllowOrigins,
	}))
	router.Use(otelgin.Middleware(""))
	router.Use(middleware.ErrorMiddleware())
	router.Use(middleware.ContextMiddleware())
	router.Use(gin.Recovery())

	// Initialize handlers
	orderHandler := orderhandler.NewOrderHandler(orderUsecase)

	// Define routes
	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST(
				"",
				middleware.RateLimitMiddleware(cacheClient, 1, 1*time.Second),
				middleware.IdempotencyMiddleware(cacheClient, 5*time.Minute),
				orderHandler.Create,
			)
		}
	}

	if config.ApplicationConfig.Env == config.ProdEnv {
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
