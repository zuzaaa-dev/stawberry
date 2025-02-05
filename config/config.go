package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	ServerPort    string
	AccessKey     string
	SecretKey     string
	BucketName    string
	URL           string
	SigningRegion string
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config reading failed: %v", err)
	}

	config := &Config{
		DBHost:        viper.GetString("DB_HOST"),
		DBUser:        viper.GetString("DB_USER"),
		DBPassword:    viper.GetString("DB_PASSWORD"),
		DBName:        viper.GetString("DB_NAME"),
		DBPort:        viper.GetString("DB_PORT"),
		ServerPort:    viper.GetString("SERVER_PORT"),
		AccessKey:     viper.GetString("ACCESS_KEY"),
		SecretKey:     viper.GetString("SECRET_KEY"),
		BucketName:    viper.GetString("BUCKET_NAME"),
		URL:           viper.GetString("URL"),
		SigningRegion: viper.GetString("SIGNING_REGION"),
	}

	return config
}

func (c *Config) GetDBConnString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}
