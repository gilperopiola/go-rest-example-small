package main

import (
	"log"

	"github.com/gilperopiola/go-rest-example-small/api"
	"github.com/gilperopiola/go-rest-example-small/api/common"
	"github.com/gilperopiola/go-rest-example-small/api/endpoints"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Note: This is the entrypoint of the application.
// The HTTP Requests entrypoint is the Prometheus HandlerFunc in pkg/common/middleware/prometheus.go

func main() {

	log.Println("Server starting ;)")

	/*-------------------------
	//      DEPENDENCIES
	//------------------------*/

	config := common.NewConfig()
	log.Println("Config OK")

	logger := logrus.New()
	logger.Info("Logger OK")

	middlewares := []gin.HandlerFunc{
		gin.Recovery(), // Panic recovery
		common.NewRateLimiterMiddleware(common.NewRateLimiter(200)),                     // Rate Limiter
		common.NewCORSConfigMiddleware(),                                                // CORS
		common.NewNewRelicMiddleware(common.NewNewRelic(config.Monitoring, logger)),     // New Relic (monitoring)
		common.NewPrometheusMiddleware(common.NewPrometheus(config.Monitoring, logger)), // Prometheus (metrics)
		common.NewTimeoutMiddleware(45),                                                 // Timeout
		common.NewErrorHandlerMiddleware(logger),                                        // Error Handler
	}
	logger.Info("Middlewares OK")

	auth := common.NewAuth(config.JWTSecret, 7)
	logger.Info("Auth OK")

	database := common.NewDatabase(config, logger)
	logger.Info("Database OK")

	handler := endpoints.NewHandler(config, database.DB, auth)
	logger.Info("Handler OK")

	router := api.NewRouter(handler, config, auth, middlewares...)
	logger.Info("Router & Endpoints OK")

	/*---------------------------
	//       START SERVER
	//--------------------------*/

	logger.Info("Running server on port " + config.Port)

	err := router.Run(":" + config.Port)
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
