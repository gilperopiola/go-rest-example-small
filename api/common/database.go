package common

import (
	"fmt"
	"log"
	"time"

	"github.com/gilperopiola/go-rest-example-small/api/common/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type database struct {
	*gorm.DB
}

func NewDatabase(config *config.Config, logger *logrus.Logger) *database {
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

func (database *database) connectToDB(config *config.Config, logger *logrus.Logger) error {
	dbConfig := config.Database
	retries := 0
	var err error

	// Retry connection if it fails due to Docker's orchestration
	for retries < dbConfig.MaxRetries {
		if database.DB, err = gorm.Open(mysql.Open(dbConfig.GetConnectionString())); err == nil {
			break
		}

		retries++
		if retries >= dbConfig.MaxRetries {
			logger.Error(fmt.Sprintf("error connecting to database after %d retries: %v", dbConfig.MaxRetries, err), nil)
			return err
		}

		logger.Info("error connecting to database, retrying... ", map[string]interface{}{})
		time.Sleep(time.Duration(dbConfig.RetryDelay) * time.Second)
	}
	return nil
}

func (database *database) configure(config *config.Config) {
	mySQLDB, _ := database.DB.DB()
	dbConfig := config.Database

	// Set connection pool limits
	mySQLDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	mySQLDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	mySQLDB.SetConnMaxLifetime(time.Hour)

	// Destroy or clean tables
	if dbConfig.Destroy {
		for _, model := range AllModels {
			database.DB.Migrator().DropTable(model)
		}
	} else if dbConfig.Clean {
		for _, model := range AllModels {
			database.DB.Delete(model)
		}
	}

	// AutoMigrate fields
	database.DB.AutoMigrate(AllModels...)

	// Insert admin user
	if dbConfig.AdminInsert {
		admin := makeAdminModel("ferra.main@gmail.com", Hash(dbConfig.AdminPassword, config.Auth.HashSalt))
		if err := database.DB.Create(admin).Error; err != nil {
			fmt.Println(err.Error())
		}
	}

	// Just for formatting the logs :)
	if config.General.LogInfo {
		fmt.Println("")
	}
}

func makeAdminModel(email, password string) *User {
	return &User{
		Username: "admin",
		Email:    email,
		Password: password,
		IsAdmin:  true,
	}
}
