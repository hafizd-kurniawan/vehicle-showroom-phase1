package config

import (
  "os"
  "strconv"
)

type Config struct {
  Database DatabaseConfig
  JWT      JWTConfig
  Server   ServerConfig
}

type DatabaseConfig struct {
  Host     string
  Port     string
  User     string
  Password string
  Name     string
  SSLMode  string
}

type JWTConfig struct {
  Secret      string
  ExpireHours int
}

type ServerConfig struct {
  Port string
  Mode string
}

func New() *Config {
  expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))

  return &Config{
    Database: DatabaseConfig{
      Host:     getEnv("DB_HOST", "localhost"),
      Port:     getEnv("DB_PORT", "5432"),
      User:     getEnv("DB_USER", "postgres"),
      Password: getEnv("DB_PASSWORD", "postgres"),
      Name:     getEnv("DB_NAME", "vehicle_showroom"),
      SSLMode:  getEnv("DB_SSLMODE", "disable"),
    },
    JWT: JWTConfig{
      Secret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
      ExpireHours: expireHours,
    },
    Server: ServerConfig{
      Port: getEnv("PORT", "8080"),
      Mode: getEnv("GIN_MODE", "debug"),
    },
  }
}

func getEnv(key, defaultValue string) string {
  if value := os.Getenv(key); value != "" {
    return value
  }
  return defaultValue
}
