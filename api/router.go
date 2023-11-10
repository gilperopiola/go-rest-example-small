package api

import (
	"github.com/gilperopiola/go-rest-example-small/api/common"
	"github.com/gilperopiola/go-rest-example-small/api/endpoints"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type router struct {
	*gin.Engine
}

func NewRouter(h endpoints.Handler, cfg *common.Config, auth common.AuthI, middlewares ...gin.HandlerFunc) router {
	var router router
	router.setup(h, cfg, auth, middlewares...)
	return router
}

func (router *router) setup(h endpoints.Handler, cfg *common.Config, auth common.AuthI, middlewares ...gin.HandlerFunc) {

	// Create router. Set debug/release mode
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Engine = gin.New()
	router.Engine.SetTrustedProxies(nil)

	// Add middlewares
	for _, middleware := range middlewares {
		router.Use(middleware)
	}

	// Set endpoints
	router.setEndpoints(h, cfg, auth)
}

/*-----------------------------
//     ROUTES / ENDPOINTS
//---------------------------*/

func (router *router) setEndpoints(h endpoints.Handler, cfg *common.Config, authI common.AuthI) {

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
