package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadConfig загружает конфиг из .env файла и переменных окружения
func LoadConfig(envFile string) (*Config, error) {
	if err := godotenv.Load(envFile); err != nil {
		return nil, fmt.Errorf("error loading env file %s: %v", envFile, err)
	}

	config := &Config{
		AppEnv:     os.Getenv("APP_ENV"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate проверяет, что все обязательные поля заполнены
func (c *Config) validate() error {
	missingFields := []string{}

	if c.DBHost == "" {
		missingFields = append(missingFields, "DB_HOST")
	}
	if c.DBPort == "" {
		missingFields = append(missingFields, "DB_PORT")
	}
	if c.DBUser == "" {
		missingFields = append(missingFields, "DB_USER")
	}
	if c.DBPassword == "" {
		missingFields = append(missingFields, "DB_PASSWORD")
	}
	if c.DBName == "" {
		missingFields = append(missingFields, "DB_NAME")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required config fields: %v", missingFields)
	}

	return nil
}

// DSN возвращает строку подключения к базе
func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
