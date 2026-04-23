package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system env")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	env := os.Getenv("APP_ENV")

	var logLevel logger.LogLevel

	switch env {
	case "dev":
		logLevel = logger.Info
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		fmt.Println("Running in DEV mode")
	case "prod":
		logLevel = logger.Error
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		logLevel = logger.Warn
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		fmt.Println("Running in DEFAULT mode")
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println(" Database connected successfully!")
}
