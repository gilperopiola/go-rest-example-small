package api

import (
	"fmt"
	"net/http/pprof"

	"github.com/gilperopiola/go-rest-example/api/common"
	"github.com/gilperopiola/go-rest-example/api/common/config"
	"github.com/gilperopiola/go-rest-example/api/endpoints"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type router struct {
	*gin.Engine
}

func NewRouter(h endpoints.Handler, cfg *config.Config, auth common.AuthI, logger *logrus.Logger, middlewares ...gin.HandlerFunc) router {
	var router router
	router.setup(h, cfg, auth, logger, middlewares...)
	return router
}

func (router *router) setup(h endpoints.Handler, cfg *config.Config, auth common.AuthI, logger *logrus.Logger, middlewares ...gin.HandlerFunc) {

	// Create router. Set debug/release mode
	router.prepare(!cfg.General.Debug, logger)

	// Add middlewares
	for _, middleware := range middlewares {
		router.Use(middleware)
	}

	// Set endpoints
	router.setEndpoints(h, cfg, auth)
}

func (router *router) prepare(isProd bool, logger *logrus.Logger) {
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Engine = gin.New()
	router.Engine.SetTrustedProxies(nil)
}

/*-----------------------------
//     ROUTES / ENDPOINTS
//---------------------------*/

func (router *router) setEndpoints(h endpoints.Handler, cfg *config.Config, authI common.AuthI) {

	// Standard endpoints
	router.GET("/health", h.HealthCheck)

	// V1
	v1 := router.Group("/v1")
	{
		router.setV1Endpoints(v1, h, authI)
	}

	// Monitoring
	if cfg.Monitoring.PrometheusEnabled {
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// Profiling
	if cfg.General.Profiling {
		router.profiling()
	}

	fmt.Println("")
}

func (router *router) setV1Endpoints(v1 *gin.RouterGroup, h endpoints.Handler, authI common.AuthI) {

	// Auth
	v1.POST("/signup", h.Signup)
	v1.POST("/login", h.Login)

	// Users
	users := v1.Group("/users", authI.ValidateToken(common.AnyRole, true))
	{
		users.GET("/:user_id", h.GetUser)
		users.PATCH("/:user_id", h.UpdateUser)
		users.DELETE("/:user_id", h.DeleteUser)
		users.PATCH("/:user_id/password", h.ChangePassword)

		// User posts
		posts := users.Group("/:user_id/posts")
		{
			posts.POST("", h.CreateUserPost)
		}
	}

	// Admins
	admin := v1.Group("/admin", authI.ValidateToken(common.AdminRole, false))
	{
		admin.POST("/user", h.CreateUser)
		admin.GET("/users", h.SearchUsers)
	}
}

// Profiling, only called if the config is set to true
func (r *router) profiling() {
	pprofGroup := r.Group("/debug/pprof")
	pprofGroup.GET("/", gin.WrapF(pprof.Index))
	pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
	pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
	pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
	pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
	pprofGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
	pprofGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
	pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
	pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
	pprofGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
	pprofGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
}
