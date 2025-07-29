package repository

import (
  "database/sql"
  "fmt"

  "github.com/jmoiron/sqlx"
  "vehicle-showroom/internal/entity"
)

type SessionRepository interface {
  Create(session *entity.UserSession) error
  GetByToken(token string) (*entity.UserSession, error)
  UpdateLogout(token string) error
  DeleteByUserID(userID int) error
}

type sessionRepository struct {
  db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) SessionRepository {
  return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *entity.UserSession) error {
  query := `
    INSERT INTO user_sessions (user_id, session_token, ip_address, is_active)
    VALUES ($1, $2, $3, $4)
    RETURNING id, login_at
  `
  
  err := r.db.QueryRow(
    query,
    session.UserID,
    session.SessionToken,
    session.IPAddress,
    session.IsActive,
  ).Scan(&session.ID, &session.LoginAt)
  
  if err != nil {
    return fmt.Errorf("failed to create session: %w", err)
  }
  
  return nil
}

func (r *sessionRepository) GetByToken(token string) (*entity.UserSession, error) {
  session := &entity.UserSession{}
  query := `
    SELECT id, user_id, session_token, login_at, logout_at, ip_address, is_active
    FROM user_sessions
    WHERE session_token = $1 AND is_active = true
  `
  
  err := r.db.Get(session, query, token)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get session by token: %w", err)
  }
  
  return session, nil
}

func (r *sessionRepository) UpdateLogout(token string) error {
  query := `
    UPDATE user_sessions
    SET logout_at = CURRENT_TIMESTAMP, is_active = false
    WHERE session_token = $1
  `
  
  _, err := r.db.Exec(query, token)
  if err != nil {
    return fmt.Errorf("failed to update logout: %w", err)
  }
  
  return nil
}

func (r *sessionRepository) DeleteByUserID(userID int) error {
  query := `
    UPDATE user_sessions
    SET is_active = false
    WHERE user_id = $1
  `
  
  _, err := r.db.Exec(query, userID)
  if err != nil {
    return fmt.Errorf("failed to delete sessions by user id: %w", err)
  }
  
  return nil
}
