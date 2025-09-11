package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Type     string // "sqlite" or "mongodb"
	Path     string // SQLite database file path
	MongoURI string // MongoDB connection URI
	MongoDB  string // MongoDB database name
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn int // hours
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DATABASE_TYPE", "sqlite"),
			Path:     getEnv("DB_PATH", "data/app.db"),
			MongoURI: getEnv("MONGODB_URI", ""),
			MongoDB:  getEnv("MONGODB_DATABASE", "chatbot"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET", "your-secret-key"),
			ExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 24),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
