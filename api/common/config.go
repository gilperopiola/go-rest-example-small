package common

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	General
	Database   Database
	Monitoring Monitoring
}

func NewConfig() *Config {
	config := Config{}
	config.setup()
	return &config
}

// If you add something here, remember to add it to the .env_example file

type General struct {
	AppName   string `envconfig:"GO_REST_EXAMPLE_APP_NAME"`
	Debug     bool   `envconfig:"GO_REST_EXAMPLE_DEBUG"`
	Port      string `envconfig:"GO_REST_EXAMPLE_PORT"`
	JWTSecret string `envconfig:"GO_REST_EXAMPLE_JWT_SECRET"`
	HashSalt  string `envconfig:"GO_REST_EXAMPLE_HASH_SALT"`
}

type Database struct {
	Type     string `envconfig:"GO_REST_EXAMPLE_DATABASE_TYPE"`
	Username string `envconfig:"GO_REST_EXAMPLE_DATABASE_USERNAME"`
	Password string `envconfig:"GO_REST_EXAMPLE_DATABASE_PASSWORD"`
	Hostname string `envconfig:"GO_REST_EXAMPLE_DATABASE_HOSTNAME"`
	Port     string `envconfig:"GO_REST_EXAMPLE_DATABASE_PORT"`
	Schema   string `envconfig:"GO_REST_EXAMPLE_DATABASE_SCHEMA"`
}

type Monitoring struct {
	NewRelicEnabled    bool   `envconfig:"GO_REST_EXAMPLE_MONITORING_NEW_RELIC_ENABLED"`
	NewRelicAppName    string `envconfig:"GO_REST_EXAMPLE_MONITORING_NEW_RELIC_APP_NAME"`
	NewRelicLicenseKey string `envconfig:"GO_REST_EXAMPLE_MONITORING_NEW_RELIC_LICENSE_KEY"`

	PrometheusEnabled bool   `envconfig:"GO_REST_EXAMPLE_MONITORING_PROMETHEUS_ENABLED"`
	PrometheusAppName string `envconfig:"GO_REST_EXAMPLE_MONITORING_PROMETHEUS_APP_NAME"`
}

func (config *Config) setup() {

	// We may be on the cmd folder or not. Hacky, I know.
	envFilePath := ".env"
	if currentFolderIsCMD() {
		envFilePath = "../.env"
	}

	// Load .env file into environment variables
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	// Parse the environment variables into the Config struct
	err = envconfig.Process("", config)
	if err != nil {
		log.Fatalf("error parsing environment variables: %v", err)
	}
}

func (dbConfig *Database) GetConnectionString() string {
	var (
		username = dbConfig.Username
		password = dbConfig.Password
		hostname = dbConfig.Hostname
		port     = dbConfig.Port
		schema   = dbConfig.Schema
		params   = "?charset=utf8&parseTime=True&loc=Local"
	)
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s%s",
		username, password, hostname, port, schema, params,
	)
}

func currentFolderIsCMD() bool {
	dir, _ := os.Getwd()

	// Extract the current directory name from the path
	dirName := strings.Split(dir, string(os.PathSeparator))
	currentDir := dirName[len(dirName)-1]

	// Check if the last 3 letters are "cmd"
	return strings.HasSuffix(currentDir, "cmd")
}
