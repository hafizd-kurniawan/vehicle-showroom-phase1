package http

import (
  "net/http"
  "strings"

  "github.com/gin-gonic/gin"
  "vehicle-showroom/internal/entity"
  "vehicle-showroom/internal/usecase"
)

type AuthHandler struct {
  authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
  return &AuthHandler{
    authUsecase: authUsecase,
  }
}

func (h *AuthHandler) Login(c *gin.Context) {
  var req entity.LoginRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid request",
      "message": err.Error(),
    })
    return
  }

  ipAddress := c.ClientIP()
  
  response, err := h.authUsecase.Login(&req, ipAddress)
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{
      "error":   "Login failed",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    response,
  })
}

func (h *AuthHandler) Register(c *gin.Context) {
  var req entity.RegisterRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid request",
      "message": err.Error(),
    })
    return
  }

  user, err := h.authUsecase.Register(&req)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Registration failed",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusCreated, gin.H{
    "success": true,
    "data":    user,
  })
}

func (h *AuthHandler) Logout(c *gin.Context) {
  token := h.extractToken(c)
  if token == "" {
    c.JSON(http.StatusUnauthorized, gin.H{
      "error":   "Unauthorized",
      "message": "No token provided",
    })
    return
  }

  if err := h.authUsecase.Logout(token); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error":   "Logout failed",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "message": "Logged out successfully",
  })
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
  token := h.extractToken(c)
  if token == "" {
    c.JSON(http.StatusUnauthorized, gin.H{
      "error":   "Unauthorized",
      "message": "No token provided",
    })
    return
  }

  user, err := h.authUsecase.GetProfile(token)
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{
      "error":   "Unauthorized",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    user,
  })
}

func (h *AuthHandler) extractToken(c *gin.Context) string {
  authHeader := c.GetHeader("Authorization")
  if authHeader == "" {
    return ""
  }

  parts := strings.Split(authHeader, " ")
  if len(parts) != 2 || parts[0] != "Bearer" {
    return ""
  }

  return parts[1]
}
