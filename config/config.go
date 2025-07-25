package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var DBName = os.Getenv("DB.NAME")
var BaseURL = os.Getenv("BASE_URL")
var SecretKey = os.Getenv("JWT.SECRET")

type Config struct {
	ServerPort string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	BaseURL    string
	SecretKey  string
}

func LoadConfig() (*Config, error) {
	// Load config here
	env := os.Getenv("ENV")
	if env == "production" {
		// For Docker, where .env is mounted at the root
		viper.AddConfigPath("/")
		viper.SetConfigFile("/.env")
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigFile(".env")
	}

	viper.SetConfigType("env")
	viper.SetDefault("sever.port", "8089")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return nil, err
	}

	viper.AutomaticEnv()
	config := &Config{
		ServerPort: viper.GetString("WEB.PORT"),
		DBUser:     viper.GetString("DB.USERNAME"),
		DBPassword: viper.GetString("DB.PASSWORD"),
		DBName:     viper.GetString("DB.NAME"),
		DBHost:     viper.GetString("DB.HOST"),
		DBPort:     viper.GetString("DB.PORT"),
		BaseURL:    viper.GetString("BASE_URL"),
		SecretKey:  viper.GetString("JWT.SECRET"),
	}
	DBName = config.DBName
	BaseURL = config.BaseURL
	SecretKey = config.SecretKey
	return config, nil
}
