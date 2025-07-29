package database

import (
  "fmt"

  "github.com/jmoiron/sqlx"
  _ "github.com/lib/pq"
  "vehicle-showroom/internal/config"
)

func NewPostgreSQL(cfg config.DatabaseConfig) (*sqlx.DB, error) {
  dsn := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
  )

  db, err := sqlx.Connect("postgres", dsn)
  if err != nil {
    return nil, fmt.Errorf("failed to connect to database: %w", err)
  }

  if err := db.Ping(); err != nil {
    return nil, fmt.Errorf("failed to ping database: %w", err)
  }

  return db, nil
}
