package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string
	Port int
	Mode string
}

type DatabaseConfig struct {
	Development Development
	Production  Production
	Test        Test
}

type Development struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type Test struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type Production struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type InfrastructureConfig struct{}

type SecurityConfig struct {
	PasetoSecret string
	BcryptCost   int
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
		log.Fatalf("Error reading config file: %v", err)
	}

	viper.WatchConfig()

	Settings = Config{
		Server: ServerConfig{
			Host: viper.GetString("SERVER_HOST"),
			Port: viper.GetInt("SERVER_PORT"),
			Mode: viper.GetString("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			Production: Production{
				Host:     viper.GetString("PROD_DB_HOST"),
				Port:     viper.GetInt("PROD_DB_PORT"),
				User:     viper.GetString("PROD_DB_USER"),
				Password: viper.GetString("PROD_DB_PASSWORD"),
				Name:     viper.GetString("PROD_DB_NAME"),
			},
			Development: Development{
				Host:     viper.GetString("DEV_DB_HOST"),
				Port:     viper.GetInt("DEV_DB_PORT"),
				User:     viper.GetString("DEV_DB_USER"),
				Password: viper.GetString("DEV_DB_PASSWORD"),
				Name:     viper.GetString("DEV_DB_NAME"),
			},
			Test: Test{
				Host:     viper.GetString("TEST_DB_HOST"),
				Port:     viper.GetInt("TEST_DB_PORT"),
				User:     viper.GetString("TEST_DB_USER"),
				Password: viper.GetString("TEST_DB_PASSWORD"),
				Name:     viper.GetString("TEST_DB_NAME"),
			},
		},
		Infrastructure: InfrastructureConfig{},
		Security: SecurityConfig{
			PasetoSecret: viper.GetString("PASETO_SECRET"),
			BcryptCost:   viper.GetInt("BCRYPT_COST"),
		},
	}

	log.Println("Configuration loaded successfully")
}
