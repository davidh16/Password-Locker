package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tkanos/gonfig"
	"golang.org/x/oauth2/google"
)

const (
	LOCAL_ENVIRONEMNT      = "local"
	DEBUG_ENVIRONMENT      = "debug"
	PRODUCTION_ENVIRONMENT = "production"
)

// Config struct can be expanded if more env variables are introduced
type Config struct {
	PgUrl                           string `env:"PG_URL"`
	DbHost                          string `env:"DB_HOST"`
	DbUser                          string `env:"DB_USER"`
	DbPassword                      string `env:"DB_PASSWORD"`
	DbName                          string `env:"DB_NAME"`
	RedisHost                       string `env:"REDIS_HOST"`
	RedisPort                       string `env:"REDIS_PORT"`
	RedisPassword                   string `env:"REDIS_PASSWORD"`
	Port                            string `env:"PORT"`
	UserSecretKey                   string `env:"USER_SECRET_KEY"`
	EntitySecretKey                 string `env:"ENTITY_SECRET_KEY"`
	EntitySecretVector              string `env:"ENTITY_SECRET_VECTOR"`
	ResponseSecretKey               string `env:"RESPONSE_SECRET_KEY"`
	ResponseSecretVector            string `env:"RESPONSE_SECRET_VECTOR"`
	StorageBucket                   string `env:"STORAGE_BUCKET"`
	SmtpPort                        string `env:"SMTP_PORT"`
	SmtpHost                        string `env:"SMTP_HOST"`
	SmtpFrom                        string `env:"SMTP_FROM"`
	SmtpPassword                    string `env:"SMTP_PASSWORD"`
	BaseUrl                         string `env:"BASE_URL"`
	DefaultEntityIconPath           string `env:"DEFAULT_ENTITY_ICON_PATH"`
	GinMode                         string `env:"GIN_MODE"`
	Environment                     string `env:"ENVIRONMENT"`
	GoogleServiceAccountCredentials *google.Credentials
}

func GetConfig() *Config {

	configuration := &Config{}

	// fileName could be changed dynamically if there are more env files (like production env and development env ), but for the purpose of this app, it will be hardcoded
	//fileName := ".env"
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env file")
	}

	err := gonfig.GetConf("", configuration)
	if err != nil {
		fmt.Println(err.Error())
	}

	return configuration
}
