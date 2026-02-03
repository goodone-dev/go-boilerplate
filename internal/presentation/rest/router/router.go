package router

import (
	"net/http/pprof"

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
		AllowOrigins: config.Cors.AllowOrigins,
		AllowMethods: config.Cors.AllowMethods,
	}

	secureConfig := secure.DefaultConfig()
	secureConfig.SSLRedirect = config.Application.Env == config.EnvProd // Only force HTTPS in production
	if config.Application.Env != config.EnvProd {                       // Disable HSTS in non-production environments
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
	router.Use(middleware.ContextTimeoutHandler())
	router.Use(middleware.RequestIdHandler())
	router.Use(middleware.IdempotencyHandler(cacheClient, config.IdempotencyDuration))
	router.Use(middleware.ErrorHandler())

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
				middleware.RateLimiterHandler(cacheClient, middleware.RateLimitConfig{
					Limit: config.RateLimiter.SingleLimit,
					TTL:   config.RateLimiter.SingleDuration,
					Mode:  middleware.SingleLimiter,
				}),
				orderHandler.Create,
			)
		}
	}

	return router
}
