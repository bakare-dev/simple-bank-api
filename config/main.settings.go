package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type InfrastructureConfig struct{}

type SecurityConfig struct {
	JWTSecret string
}

type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	Infrastructure InfrastructureConfig
	Security       SecurityConfig
}

var Settings Config

func LoadConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	Settings = Config{
		Server: ServerConfig{
			Host: viper.GetString("SERVER_HOST"),
			Port: viper.GetInt("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		Infrastructure: InfrastructureConfig{},
		Security: SecurityConfig{
			JWTSecret: viper.GetString("JWT_SECRET"),
		},
	}

	log.Println("Configuration loaded successfully")
}
