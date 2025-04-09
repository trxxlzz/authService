package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DBType string

const (
	MongoDB    DBType = "mongo"
	PostgreSQL DBType = "postgres"
)

type Config struct {
	AppEnv string

	// MongoDB
	MongoHost string
	MongoPort string
	MongoUser string
	MongoPass string
	MongoDB   string

	// PostgreSQL
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

// LoadConfig загружает конфиг из .env файла и переменных окружения
func LoadConfig(envFile string, dbType DBType) (*Config, error) {
	if err := godotenv.Load(envFile); err != nil {
		return nil, fmt.Errorf("error loading env file %s: %v", envFile, err)
	}

	config := &Config{
		AppEnv:        os.Getenv("APP_ENV"),
		MongoHost:     os.Getenv("MONGO_HOST"),
		MongoPort:     os.Getenv("MONGO_PORT"),
		MongoUser:     os.Getenv("MONGO_USER"),
		MongoPass:     os.Getenv("MONGO_PASSWORD"),
		MongoDB:       os.Getenv("MONGO_DB"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}

	if err := config.validate(dbType); err != nil {
		return nil, err
	}

	return config, nil
}

// validate проверяет, что все обязательные поля заполнены
func (c *Config) validate(dbType DBType) error {
	missingFields := []string{}

	switch dbType {
	case MongoDB:
		if c.MongoHost == "" {
			missingFields = append(missingFields, "MONGO_HOST")
		}
		if c.MongoPort == "" {
			missingFields = append(missingFields, "MONGO_PORT")
		}
		if c.MongoUser == "" {
			missingFields = append(missingFields, "MONGO_USER")
		}
		if c.MongoPass == "" {
			missingFields = append(missingFields, "MONGO_PASSWORD")
		}
		if c.MongoDB == "" {
			missingFields = append(missingFields, "MONGO_DB")
		}
	case PostgreSQL:
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
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required config fields: %v", missingFields)
	}

	return nil
}

// DSN возвращает строку подключения к нужной базе данных
func (c *Config) DSN(dbType DBType) string {
	switch dbType {
	case MongoDB:
		return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
			c.MongoUser, c.MongoPass, c.MongoHost, c.MongoPort, c.MongoDB)
	case PostgreSQL:
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	default:
		return ""
	}
}

// RedisAddr возвращает строку подключения к Redis
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
