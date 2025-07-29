package repository

import (
  "database/sql"
  "fmt"

  "github.com/jmoiron/sqlx"
  "vehicle-showroom/internal/entity"
)

type UserRepository interface {
  Create(user *entity.User) error
  GetByUsername(username string) (*entity.User, error)
  GetByEmail(email string) (*entity.User, error)
  GetByID(id int) (*entity.User, error)
  Update(user *entity.User) error
  Delete(id int) error
}

type userRepository struct {
  db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
  return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
  query := `
    INSERT INTO users (username, email, password_hash, full_name, phone, role, is_active)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, created_at, updated_at
  `
  
  err := r.db.QueryRow(
    query,
    user.Username,
    user.Email,
    user.PasswordHash,
    user.FullName,
    user.Phone,
    user.Role,
    user.IsActive,
  ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
  
  if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
  }
  
  return nil
}

func (r *userRepository) GetByUsername(username string) (*entity.User, error) {
  user := &entity.User{}
  query := `
    SELECT id, username, email, password_hash, full_name, phone, role, is_active, created_at, updated_at
    FROM users
    WHERE username = $1 AND is_active = true
  `
  
  err := r.db.Get(user, query, username)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get user by username: %w", err)
  }
  
  return user, nil
}

func (r *userRepository) GetByEmail(email string) (*entity.User, error) {
  user := &entity.User{}
  query := `
    SELECT id, username, email, password_hash, full_name, phone, role, is_active, created_at, updated_at
    FROM users
    WHERE email = $1 AND is_active = true
  `
  
  err := r.db.Get(user, query, email)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get user by email: %w", err)
  }
  
  return user, nil
}

func (r *userRepository) GetByID(id int) (*entity.User, error) {
  user := &entity.User{}
  query := `
    SELECT id, username, email, password_hash, full_name, phone, role, is_active, created_at, updated_at
    FROM users
    WHERE id = $1 AND is_active = true
  `
  
  err := r.db.Get(user, query, id)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get user by id: %w", err)
  }
  
  return user, nil
}

func (r *userRepository) Update(user *entity.User) error {
  query := `
    UPDATE users
    SET username = $1, email = $2, full_name = $3, phone = $4, role = $5, updated_at = CURRENT_TIMESTAMP
    WHERE id = $6
  `
  
  _, err := r.db.Exec(query, user.Username, user.Email, user.FullName, user.Phone, user.Role, user.ID)
  if err != nil {
    return fmt.Errorf("failed to update user: %w", err)
  }
  
  return nil
}

func (r *userRepository) Delete(id int) error {
  query := `UPDATE users SET is_active = false WHERE id = $1`
  
  _, err := r.db.Exec(query, id)
  if err != nil {
    return fmt.Errorf("failed to delete user: %w", err)
  }
  
  return nil
}
