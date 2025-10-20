package router

import (
	"net/http/pprof"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	healthhandler "github.com/goodone-dev/go-boilerplate/internal/application/health/http"
	orderhandler "github.com/goodone-dev/go-boilerplate/internal/application/order/delivery/http"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/cache"
	"github.com/goodone-dev/go-boilerplate/internal/presentation/rest/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(healthHandler *healthhandler.HealthHandler, orderHandler *orderhandler.OrderHandler, cacheClient cache.ICache) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	// ========== Initialize Router ==========
	router := gin.New()
	router.Use(otelgin.Middleware(""))
	router.Use(cors.New(cors.Config{
		AllowOrigins: config.CorsConfig.AllowOrigins,
		AllowMethods: config.CorsConfig.AllowMethods,
	}))
	// router.Use(secure.New(secure.DefaultConfig()))
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.Use(middleware.ErrorMiddleware())
	router.Use(middleware.TimeoutMiddleware())
	router.Use(middleware.ResponseMiddleware())
	router.Use(gin.Recovery())

	// ========== Define Routes ==========
	health := router.Group("/health")
	{
		health.GET("", healthHandler.HealthCheck)
		health.GET("/ready", healthHandler.HealthReadyCheck)
	}

	debug := router.Group("/debug/pprof")
	{
		debug.GET("/goroutine", gin.WrapF(pprof.Index))
		debug.GET("/profile", gin.WrapF(pprof.Profile))
		debug.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		debug.GET("/symbol", gin.WrapF(pprof.Symbol))
		debug.GET("/trace", gin.WrapF(pprof.Trace))
	}

	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST(
				"",
				middleware.SingleLimiterMiddleware(cacheClient, 60, 1*time.Minute),
				middleware.IdempotencyMiddleware(cacheClient, 5*time.Minute),
				orderHandler.Create,
			)
		}
	}

	return router
}
