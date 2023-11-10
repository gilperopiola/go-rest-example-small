package common

import (
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type database struct {
	*gorm.DB
}

func NewDatabase(config *Config, logger *logrus.Logger) *database {
	var database database

	// Create connection. It's deferred closed in main.go.
	// Retry connection if it fails due to Docker's orchestration.
	if err := database.connectToDB(config, logger); err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// Set connection pool limits
	// Log queries if debug = true
	// Destroy or clean tables
	// AutoMigrate fields
	// Create admin
	database.configure(config)

	return &database
}

func (database *database) connectToDB(config *Config, logger *logrus.Logger) error {

	var (
		dbConfig   = config.Database
		retries    = 0
		maxRetries = 5
		err        error
	)

	// Retry connection if it fails due to Docker's orchestration
	for retries < maxRetries {
		if database.DB, err = gorm.Open(mysql.Open(dbConfig.GetConnectionString())); err == nil {
			break
		}

		retries++
		if retries >= maxRetries {
			logger.Error(fmt.Sprintf("error connecting to database after %d retries: %v", maxRetries, err), nil)
			return err
		}

		logger.Info("error connecting to database, retrying... ", map[string]interface{}{})
		time.Sleep(time.Duration(5) * time.Second)
	}
	return nil
}

func (database *database) configure(config *Config) {
	mySQLDB, _ := database.DB.DB()

	// Set connection pool limits
	mySQLDB.SetMaxIdleConns(100)
	mySQLDB.SetMaxOpenConns(100)
	mySQLDB.SetConnMaxLifetime(time.Hour)

	// AutoMigrate fields
	database.DB.AutoMigrate(AllModels...)
}
