package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

type InfrastructureConfig struct {
}

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
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("No .env file found, loading environment variables from system")
	}

	Settings = Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "app"),
		},
		Infrastructure: InfrastructureConfig{},
		Security: SecurityConfig{
			JWTSecret: getEnv("JWT_SECRET", "default-secret"),
		},
	}
	log.Println("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
