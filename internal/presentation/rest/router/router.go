package router

import (
	"net/http/pprof"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/domain/health"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/cache"
	"github.com/goodone-dev/go-boilerplate/internal/presentation/rest/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(healthHandler health.HealthHandler, orderHandler order.OrderHandler, cacheClient cache.Cache) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	// ========== Middleware Config ==========
	corsConfig := cors.Config{
		AllowOrigins: config.CorsConfig.AllowOrigins,
		AllowMethods: config.CorsConfig.AllowMethods,
	}

	secureConfig := secure.DefaultConfig()
	secureConfig.SSLRedirect = config.ApplicationConfig.Env == config.EnvProd // Only force HTTPS in production
	if config.ApplicationConfig.Env != config.EnvProd {                       // Disable HSTS in non-production environments
		secureConfig.STSSeconds = 0
		secureConfig.STSIncludeSubdomains = false
	}

	// ========== Initialize Router ==========
	router := gin.New()

	// Library Middleware
	router.Use(otelgin.Middleware(""))
	router.Use(cors.New(corsConfig))
	router.Use(secure.New(secureConfig))
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Internal Middleware
	router.Use(middleware.ErrorMiddleware())
	router.Use(middleware.TimeoutMiddleware())
	router.Use(middleware.ResponseMiddleware())

	router.Use(gin.Recovery())

	// ========== Define Routes ==========
	health := router.Group("/health")
	{
		health.GET("", healthHandler.LiveCheck)
		health.GET("/ready", healthHandler.ReadyCheck)
	}

	// TODO: Add authentication
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
