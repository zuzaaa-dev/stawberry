package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	ServerPort    string
	AccessKey     string
	SecretKEY     string
	BucketName    string
	URL           string
	SigningRegion string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "price_negotiation"),
		DBPort:     getEnv("DB_PORT", "5432"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		AccessKey:     getEnv("ACCESS_KEY", "YCAJEA42mQqWIYN-INQRF0DTI"),
		SecretKEY:     getEnv("SECRET_KEY", "YCM1PQy6UGDHD7sKi2Q3nBDDyur8fZFEfSEC33Fn"),
		BucketName:    getEnv("BUCKET_NAME", "stawberry"),
		URL:           getEnv("URL", "https://storage.yandexcloud.net"),
		SigningRegion: getEnv("SIGNING_REGION", "ru-central1"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDBConnString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}
