package handlers

import (
	"log"
	"net/http"
	"sterling-hms-backend/internal/config"
	"sterling-hms-backend/internal/models"
	"sterling-hms-backend/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	adminRepo *repositories.AdminRepository
	cfg       *config.Config
}

func NewAdminHandler(adminRepo *repositories.AdminRepository, cfg *config.Config) *AdminHandler {
	return &AdminHandler{
		adminRepo: adminRepo,
		cfg:       cfg,
	}
}

// AdminLogin handles admin authentication
func (h *AdminHandler) AdminLogin(c *gin.Context) {
	var req models.AdminLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid credentials",
			"success": false,
		})
		return
	}

	log.Printf("Admin login attempt for email: %s", req.Email)

	// Find admin by email
	admin, err := h.adminRepo.FindByEmail(req.Email)
	if err != nil {
		log.Printf("Admin not found for email %s: %v", req.Email, err)
		// Generic error message to prevent email enumeration
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
			"success": false,
		})
		return
	}

	log.Printf("Admin found: ID=%d, Email=%s, Active=%v", admin.ID, admin.Email, admin.IsActive)
	log.Printf("Password hash from DB: %s", admin.PasswordHash[:20]+"...")

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Printf("Password verification failed for admin %s: %v", req.Email, err)
		// Generic error message to prevent password enumeration
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
			"success": false,
		})
		return
	}

	log.Printf("Password verified successfully for admin %s", req.Email)

	// Check if admin is active
	if !admin.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Account is inactive",
			"success": false,
		})
		return
	}

	// Generate JWT token
	token, err := h.generateAdminToken(admin)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Prepare admin response (without sensitive data)
	adminResponse := &models.AdminUser{
		ID:        admin.ID,
		Email:     admin.Email,
		Name:      admin.Name,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
		IsActive:  admin.IsActive,
	}

	c.JSON(http.StatusOK, models.AdminLoginResponse{
		Message: "Login successful",
		Success: true,
		Admin:   adminResponse,
		Token:   token,
	})
}

// GetDashboardStats returns dashboard statistics
func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	// Admin ID is set by admin middleware
	adminID, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"success": false,
		})
		return
	}

	// Get stats from repository
	stats, err := h.adminRepo.GetDashboardStats()
	if err != nil {
		log.Printf("Error fetching stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch statistics",
			"success": false,
		})
		return
	}

	// Log admin action
	adminIDInt := adminID.(int)
	h.adminRepo.LogAdminAction(
		adminIDInt,
		"view_dashboard",
		"",
		nil,
		"",
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, models.AdminDashboardResponse{
		Message: "Statistics retrieved successfully",
		Success: true,
		Stats:   *stats,
	})
}

// AdminLogout handles admin logout
func (h *AdminHandler) AdminLogout(c *gin.Context) {
	adminID, exists := c.Get("adminID")
	if exists {
		adminIDInt := adminID.(int)
		h.adminRepo.LogAdminAction(
			adminIDInt,
			"logout",
			"",
			nil,
			"",
			c.ClientIP(),
			c.GetHeader("User-Agent"),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
		"success": true,
	})
}

// generateAdminToken generates a JWT token for admin
func (h *AdminHandler) generateAdminToken(admin *models.AdminUser) (string, error) {
	parsedDuration, err := time.ParseDuration(h.cfg.JWTExpire)
	if err != nil {
		parsedDuration = 24 * time.Hour // Default 24 hours for admin
	}

	claims := models.AdminClaims{
		ID:    admin.ID,
		Email: admin.Email,
		Role:  admin.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(parsedDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
