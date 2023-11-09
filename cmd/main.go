package main

import (
	"log"

	"github.com/gilperopiola/go-rest-example-small/api"
	"github.com/gilperopiola/go-rest-example-small/api/common"
	"github.com/gilperopiola/go-rest-example-small/api/common/config"
	"github.com/gilperopiola/go-rest-example-small/api/endpoints"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// Note: This is the entrypoint of the application.
// The HTTP Requests entrypoint is the Prometheus HandlerFunc in pkg/common/middleware/prometheus.go

func main() {

	log.Println("Server starting ;)")

	/*-------------------------
	//      DEPENDENCIES
	//------------------------*/

	config := config.New()
	log.Println("Config OK")

	logger := logrus.New()
	logger.Info("Logger OK", nil)

	middlewares := []gin.HandlerFunc{
		gin.Recovery(), // Panic recovery
		common.NewRateLimiterMiddleware(common.NewRateLimiter(200)),                     // Rate Limiter
		common.NewCORSConfigMiddleware(),                                                // CORS
		common.NewNewRelicMiddleware(common.NewNewRelic(config.Monitoring, logger)),     // New Relic (monitoring)
		common.NewPrometheusMiddleware(common.NewPrometheus(config.Monitoring, logger)), // Prometheus (metrics)
		common.NewTimeoutMiddleware(config.General.Timeout),                             // Timeout
		common.NewErrorHandlerMiddleware(logger),                                        // Error Handler
	}
	logger.Info("Middlewares OK", nil)

	auth := common.NewAuth(config.Auth.JWTSecret, config.Auth.SessionDurationDays)
	logger.Info("Auth OK", nil)

	database := common.NewDatabase(config, logger)
	logger.Info("Database OK", nil)

	handler := endpoints.NewHandler(config, database.DB, auth)
	logger.Info("Handler OK", nil)

	router := api.NewRouter(handler, config, auth, logger, middlewares...)
	logger.Info("Router & Endpoints OK", nil)

	/*---------------------------
	//       START SERVER
	//--------------------------*/

	port := config.General.Port
	logger.Info("Running server on port "+port, nil)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	/* Have a great day! :) */
}

// TODO
// - Redis
// - More tests
// - Batch insert
// - Reset password
// - Roles to DB
// - Request IDs
// - Logic from DeleteUser to service layer
// - Search & Fix TODOs
// - Replace user.Exists when you can
// - OpenAPI (Swagger)
