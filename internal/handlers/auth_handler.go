package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"regexp"
	"sterling-hms-backend/internal/config"
	"sterling-hms-backend/internal/middleware"
	"sterling-hms-backend/internal/models"
	"sterling-hms-backend/internal/repositories"
	"sterling-hms-backend/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo          *repositories.UserRepository
	passwordResetRepo *repositories.PasswordResetRepository
	cfg               *config.Config
	emailService      *utils.EmailService
	auditLog          *utils.AuditLog
	passwordValidator *utils.PasswordValidator
	resetFlow         *utils.PasswordResetFlow
}

func NewAuthHandler(userRepo *repositories.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo:          userRepo,
		passwordResetRepo: repositories.NewPasswordResetRepository(config.DB),
		cfg:               cfg,
		emailService:      utils.NewEmailService(),
		auditLog:          utils.NewAuditLog(),
		passwordValidator: utils.NewPasswordValidator(),
		resetFlow:         utils.NewPasswordResetFlow(),
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req models.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email format"})
		return
	}

	// Check if user already exists
	exists, err := h.userRepo.EmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	// Create user
	user, err := h.userRepo.Create(req.FirstName, req.LastName, req.Email, string(hashedPassword))
	if err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(h.cfg, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Message: "User registered successfully",
		User:    user,
		Token:   token,
	})
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req models.SignInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email and password are required"})
		return
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// Update last login
	h.userRepo.UpdateLastLogin(user.ID)

	// Generate token
	token, err := middleware.GenerateToken(h.cfg, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	// Don't send password in response
	user.Password = ""

	c.JSON(http.StatusOK, models.AuthResponse{
		Message: "Logged in successfully",
		User:    user,
		Token:   token,
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	user, err := h.userRepo.FindByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		User: user,
	})
}

func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.GetAll()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// ForgotPassword initiates password reset process - sends reset link via email
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": utils.ErrorMessages.InvalidEmail})
		return
	}

	ipAddress := h.auditLog.GetIPAddress(c)
	userAgent := h.auditLog.GetUserAgent(c)

	// Check rate limiting: 5 requests per hour per email
	attemptCount, err := h.passwordResetRepo.CountResetAttempts(req.Email, 1)
	if err != nil {
		log.Printf("Error checking rate limit: %v", err)
		h.passwordResetRepo.CreatePasswordResetLog(nil, req.Email, "forgot_password_request", ipAddress, userAgent, false, stringPtr(utils.ErrorMessages.ServerError))
		c.JSON(http.StatusInternalServerError, gin.H{"message": utils.ErrorMessages.ServerError})
		return
	}

	if attemptCount >= h.resetFlow.MaxResetRequests {
		log.Printf("Rate limit exceeded for email: %s", req.Email)
		h.passwordResetRepo.CreatePasswordResetLog(nil, req.Email, "forgot_password_request", ipAddress, userAgent, false, stringPtr("rate_limited"))
		c.JSON(http.StatusTooManyRequests, gin.H{"message": utils.ErrorMessages.RateLimited})
		return
	}

	// Check if user exists (without revealing if email exists)
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists - send generic success response
		log.Printf("User not found for email: %s", req.Email)
		h.passwordResetRepo.CreatePasswordResetLog(nil, req.Email, "forgot_password_request", ipAddress, userAgent, true, nil)
		c.JSON(http.StatusOK, models.ForgotPasswordResponse{
			Message: "If an account exists with this email, you will receive a password reset link.",
			Success: true,
		})
		return
	}

	// Generate reset token
	tokenGen := utils.NewTokenGenerator(32)
	resetToken, err := tokenGen.GenerateToken()
	if err != nil {
		log.Printf("Error generating reset token: %v", err)
		h.passwordResetRepo.CreatePasswordResetLog(&user.ID, req.Email, "forgot_password_request", ipAddress, userAgent, false, stringPtr(utils.ErrorMessages.ServerError))
		c.JSON(http.StatusInternalServerError, gin.H{"message": utils.ErrorMessages.ServerError})
		return
	}

	// Hash token
	tokenHash := hashValue(resetToken)

	// Store in database with 1 hour expiry
	expiresAt := time.Now().Add(1 * time.Hour)
	_, err = h.passwordResetRepo.CreatePasswordResetToken(user.ID, tokenHash, expiresAt)
	if err != nil {
		log.Printf("Error creating reset token: %v", err)
		h.passwordResetRepo.CreatePasswordResetLog(&user.ID, req.Email, "forgot_password_request", ipAddress, userAgent, false, stringPtr(utils.ErrorMessages.ServerError))
		c.JSON(http.StatusInternalServerError, gin.H{"message": utils.ErrorMessages.ServerError})
		return
	}

	// Build reset link - frontend will handle this
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	resetLink := frontendURL + "/reset-password?token=" + resetToken

	// Send password reset email with link
	err = h.emailService.SendPasswordResetEmail(req.Email, user.FirstName, resetLink)
	if err != nil {
		log.Printf("Error sending reset email: %v", err)
		// Don't fail the request if email fails (for development purposes)
	}

	// Log successful request
	h.passwordResetRepo.CreatePasswordResetLog(&user.ID, req.Email, "forgot_password_request", ipAddress, userAgent, true, nil)

	c.JSON(http.StatusOK, models.ForgotPasswordResponse{
		Message: "If an account exists with this email, you will receive a password reset link.",
		Success: true,
	})
}

// ResetPassword resets user password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	ipAddress := h.auditLog.GetIPAddress(c)
	userAgent := h.auditLog.GetUserAgent(c)

	// Validate password
	if err := h.passwordValidator.ValidatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Verify reset token
	token, err := h.passwordResetRepo.GetPasswordResetTokenByHash(req.ResetToken)
	if err != nil {
		log.Printf("Invalid reset token: %v", err)
		h.passwordResetRepo.CreatePasswordResetLog(nil, "", "reset_failed", ipAddress, userAgent, false, stringPtr("invalid_token"))
		c.JSON(http.StatusUnauthorized, gin.H{"message": utils.ErrorMessages.InvalidToken})
		return
	}

	// Get user
	user, err := h.userRepo.FindByID(token.UserID)
	if err != nil {
		log.Printf("User not found: %v", err)
		h.passwordResetRepo.CreatePasswordResetLog(nil, "", "reset_failed", ipAddress, userAgent, false, stringPtr("user_not_found"))
		c.JSON(http.StatusUnauthorized, gin.H{"message": utils.ErrorMessages.InvalidToken})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		h.passwordResetRepo.CreatePasswordResetLog(&user.ID, user.Email, "reset_failed", ipAddress, userAgent, false, stringPtr(utils.ErrorMessages.ServerError))
		c.JSON(http.StatusInternalServerError, gin.H{"message": utils.ErrorMessages.ServerError})
		return
	}

	// Update password
	err = h.passwordResetRepo.UpdateUserPassword(user.ID, string(hashedPassword))
	if err != nil {
		log.Printf("Error updating password: %v", err)
		h.passwordResetRepo.CreatePasswordResetLog(&user.ID, user.Email, "reset_failed", ipAddress, userAgent, false, stringPtr(utils.ErrorMessages.ServerError))
		c.JSON(http.StatusInternalServerError, gin.H{"message": utils.ErrorMessages.ServerError})
		return
	}

	// Mark token as used
	h.passwordResetRepo.MarkResetTokenAsUsed(token.ID)

	// Invalidate all other reset tokens for this user
	h.passwordResetRepo.InvalidateAllUserTokens(user.ID)

	// Log successful reset
	h.passwordResetRepo.CreatePasswordResetLog(&user.ID, user.Email, "reset_success", ipAddress, userAgent, true, nil)

	// Send confirmation email
	err = h.emailService.SendPasswordResetSuccessEmail(user.Email, user.FirstName)
	if err != nil {
		log.Printf("Error sending confirmation email: %v", err)
	}

	c.JSON(http.StatusOK, models.ResetPasswordResponse{
		Message: "Password reset successfully. Please log in with your new password.",
		Success: true,
	})
}

// Helper functions
func hashValue(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func stringPtr(s string) *string {
	return &s
}
