package entity

import (
  "time"
)

type User struct {
  ID           int       `json:"id" db:"id"`
  Username     string    `json:"username" db:"username"`
  Email        string    `json:"email" db:"email"`
  PasswordHash string    `json:"-" db:"password_hash"`
  FullName     string    `json:"full_name" db:"full_name"`
  Phone        *string   `json:"phone" db:"phone"`
  Role         string    `json:"role" db:"role"`
  IsActive     bool      `json:"is_active" db:"is_active"`
  CreatedAt    time.Time `json:"created_at" db:"created_at"`
  UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserSession struct {
  ID           int        `json:"id" db:"id"`
  UserID       int        `json:"user_id" db:"user_id"`
  SessionToken string     `json:"session_token" db:"session_token"`
  LoginAt      time.Time  `json:"login_at" db:"login_at"`
  LogoutAt     *time.Time `json:"logout_at" db:"logout_at"`
  IPAddress    *string    `json:"ip_address" db:"ip_address"`
  IsActive     bool       `json:"is_active" db:"is_active"`
}

type LoginRequest struct {
  Username string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
  Username string `json:"username" binding:"required,min=3,max=50"`
  Email    string `json:"email" binding:"required,email"`
  Password string `json:"password" binding:"required,min=6"`
  FullName string `json:"full_name" binding:"required"`
  Phone    string `json:"phone"`
  Role     string `json:"role" binding:"required,oneof=admin mechanic cashier"`
}

type LoginResponse struct {
  Token string `json:"token"`
  User  User   `json:"user"`
}
