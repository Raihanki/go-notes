package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type MainConfig struct {
	APP_NAME         string
	APP_PORT         string
	APP_BACKEND_URL  string
	APP_FRONTEND_URL string

	DB_DIALECT  string
	DB_HOST     string
	DB_PORT     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string

	JWT_SECRET string
}

var ENV *MainConfig

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	errReadConfig := viper.ReadInConfig()
	if errReadConfig != nil {
		log.Fatal("Failed to load configuration file")
	}

	errUnmarshal := viper.Unmarshal(&ENV)
	if errUnmarshal != nil {
		log.Fatal("Failed when mapping configuration file")
	}

	log.Info("Configuration file successfully loaded")
}
