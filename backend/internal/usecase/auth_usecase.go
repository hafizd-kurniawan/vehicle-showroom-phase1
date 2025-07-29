package usecase

import (
  "errors"
  "fmt"
  "time"

  "github.com/golang-jwt/jwt/v5"
  "golang.org/x/crypto/bcrypt"
  "vehicle-showroom/internal/config"
  "vehicle-showroom/internal/entity"
  "vehicle-showroom/internal/repository"
)

type AuthUsecase interface {
  Login(req *entity.LoginRequest, ipAddress string) (*entity.LoginResponse, error)
  Register(req *entity.RegisterRequest) (*entity.User, error)
  Logout(token string) error
  GetProfile(token string) (*entity.User, error)
  ValidateToken(token string) (*entity.User, error)
}

type authUsecase struct {
  userRepo    repository.UserRepository
  sessionRepo repository.SessionRepository
  jwtConfig   config.JWTConfig
}

func NewAuthUsecase(
  userRepo repository.UserRepository,
  sessionRepo repository.SessionRepository,
  jwtConfig config.JWTConfig,
) AuthUsecase {
  return &authUsecase{
    userRepo:    userRepo,
    sessionRepo: sessionRepo,
    jwtConfig:   jwtConfig,
  }
}

func (u *authUsecase) Login(req *entity.LoginRequest, ipAddress string) (*entity.LoginResponse, error) {
  // Get user by username
  user, err := u.userRepo.GetByUsername(req.Username)
  if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
  }
  
  if user == nil {
    return nil, errors.New("invalid username or password")
  }
  
  // Verify password
  if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
    return nil, errors.New("invalid username or password")
  }
  
  // Generate JWT token
  token, err := u.generateJWT(user)
  if err != nil {
    return nil, fmt.Errorf("failed to generate token: %w", err)
  }
  
  // Create session
  session := &entity.UserSession{
    UserID:       user.ID,
    SessionToken: token,
    IPAddress:    &ipAddress,
    IsActive:     true,
  }
  
  if err := u.sessionRepo.Create(session); err != nil {
    return nil, fmt.Errorf("failed to create session: %w", err)
  }
  
  return &entity.LoginResponse{
    Token: token,
    User:  *user,
  }, nil
}

func (u *authUsecase) Register(req *entity.RegisterRequest) (*entity.User, error) {
  // Check if username already exists
  existingUser, err := u.userRepo.GetByUsername(req.Username)
  if err != nil {
    return nil, fmt.Errorf("failed to check username: %w", err)
  }
  if existingUser != nil {
    return nil, errors.New("username already exists")
  }
  
  // Check if email already exists
  existingUser, err = u.userRepo.GetByEmail(req.Email)
  if err != nil {
    return nil, fmt.Errorf("failed to check email: %w", err)
  }
  if existingUser != nil {
    return nil, errors.New("email already exists")
  }
  
  // Hash password
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
  if err != nil {
    return nil, fmt.Errorf("failed to hash password: %w", err)
  }
  
  // Create user
  user := &entity.User{
    Username:     req.Username,
    Email:        req.Email,
    PasswordHash: string(hashedPassword),
    FullName:     req.FullName,
    Phone:        &req.Phone,
    Role:         req.Role,
    IsActive:     true,
  }
  
  if err := u.userRepo.Create(user); err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
  }
  
  return user, nil
}

func (u *authUsecase) Logout(token string) error {
  return u.sessionRepo.UpdateLogout(token)
}

func (u *authUsecase) GetProfile(token string) (*entity.User, error) {
  return u.ValidateToken(token)
}

func (u *authUsecase) ValidateToken(tokenString string) (*entity.User, error) {
  // Parse JWT token
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(u.jwtConfig.Secret), nil
  })
  
  if err != nil {
    return nil, fmt.Errorf("failed to parse token: %w", err)
  }
  
  if !token.Valid {
    return nil, errors.New("invalid token")
  }
  
  // Extract claims
  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    return nil, errors.New("invalid token claims")
  }
  
  userID, ok := claims["user_id"].(float64)
  if !ok {
    return nil, errors.New("invalid user id in token")
  }
  
  // Check session
  session, err := u.sessionRepo.GetByToken(tokenString)
  if err != nil {
    return nil, fmt.Errorf("failed to get session: %w", err)
  }
  
  if session == nil || !session.IsActive {
    return nil, errors.New("session not found or inactive")
  }
  
  // Get user
  user, err := u.userRepo.GetByID(int(userID))
  if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
  }
  
  if user == nil {
    return nil, errors.New("user not found")
  }
  
  return user, nil
}

func (u *authUsecase) generateJWT(user *entity.User) (string, error) {
  claims := jwt.MapClaims{
    "user_id":  user.ID,
    "username": user.Username,
    "role":     user.Role,
    "exp":      time.Now().Add(time.Hour * time.Duration(u.jwtConfig.ExpireHours)).Unix(),
    "iat":      time.Now().Unix(),
  }
  
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString([]byte(u.jwtConfig.Secret))
}
