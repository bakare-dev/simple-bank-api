package util

import (
	"fmt"
	"log"

	"github.com/bakare-dev/simple-bank-api/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func ConnectDB() *gorm.DB {
	var dbConfig DbConfig

	if config.Settings.Server.Env == "production" {
		dbConfig = DbConfig{
			Host:     config.Settings.Database.Production.Host,
			Port:     config.Settings.Database.Production.Port,
			User:     config.Settings.Database.Production.User,
			Password: config.Settings.Database.Production.Password,
			Name:     config.Settings.Database.Production.Name,
		}
	} else {
		dbConfig = DbConfig{
			Host:     config.Settings.Database.Development.Host,
			Port:     config.Settings.Database.Development.Port,
			User:     config.Settings.Database.Development.User,
			Password: config.Settings.Database.Development.Password,
			Name:     config.Settings.Database.Development.Name,
		}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(200)

	log.Println("Database connected successfully")

	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		log.Fatalf("failed to auto-migrate database schemas: %v", err)
	}

	log.Println("Database migrated successfully")
}
