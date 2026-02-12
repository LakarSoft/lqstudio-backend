package handlers

import (
	"strings"

	"lqstudio-backend/internal/config"
	"lqstudio-backend/internal/repositories"
	"lqstudio-backend/pkg/auth"
	"lqstudio-backend/pkg/errors"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userRepo *repositories.UserRepository
	cfg      *config.Config
	logger   *zap.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo *repositories.UserRepository, cfg *config.Config, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		cfg:      cfg,
		logger:   logger,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo represents user information in the response
type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// Login godoc
// @Summary Admin login
// @Description Authenticate admin user and receive JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body handlers.LoginRequest true "Login credentials"
// @Success 200 {object} dto.ApiResponse{data=handlers.LoginResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Failed to bind login request", zap.Error(err))
		return errors.NewValidationError("Invalid request body")
	}

	// Validate input
	if req.Email == "" {
		return errors.NewValidationError("Email is required")
	}
	if req.Password == "" {
		return errors.NewValidationError("Password is required")
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	h.logger.Info("Login attempt", zap.String("email", req.Email))

	// Get admin user by email
	user, err := h.userRepo.GetAdminByEmail(c.Request().Context(), req.Email)
	if err != nil {
		h.logger.Warn("Admin user not found", zap.String("email", req.Email), zap.Error(err))
		// Don't reveal if user exists or not for security
		return errors.NewUnauthorizedError("Invalid email or password")
	}

	// Mask hash for logging (show first 20 chars only for debugging)
	maskedHash := user.PasswordHash
	if len(maskedHash) > 20 {
		maskedHash = maskedHash[:20] + "..."
	}

	h.logger.Debug("User found",
		zap.String("user_id", user.ID),
		zap.String("user_email", user.Email),
		zap.String("user_role", user.Role),
		zap.Bool("has_password_hash", user.PasswordHash != ""),
		zap.String("hash_prefix", maskedHash),
	)

	// Verify password
	if user.PasswordHash == "" {
		h.logger.Warn("User has no password hash", zap.String("email", req.Email))
		return errors.NewUnauthorizedError("Invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		h.logger.Warn("Password verification failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return errors.NewUnauthorizedError("Invalid email or password")
	}

	h.logger.Info("Login successful", zap.String("email", req.Email))

	// Generate JWT token
	token, err := auth.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
		h.cfg.JWT.Secret,
		h.cfg.JWT.ExpiryHours,
	)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.String("email", req.Email), zap.Error(err))
		return errors.NewInternalError("Failed to generate token", err)
	}

	// Return response
	response := LoginResponse{
		Token: token,
		User: UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}

	return SendOK(c, response, "Login successful")
}
