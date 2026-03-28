package middleware

import (
	"errors"
	"net/http"
	"sterling-hms-backend/internal/config"
	"sterling-hms-backend/internal/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok || !token.Valid {
			c.Abort()
			return
		}

		c.Set("userID", claims.ID)
		c.Next()
	}
}

func GenerateToken(cfg *config.Config, user *models.User) (string, error) {
	parsedDuration, err := time.ParseDuration(cfg.JWTExpire)
	if err != nil {
		parsedDuration = 7 * 24 * time.Hour
	}

	claims := models.Claims{
		ID:    user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(parsedDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

// AdminAuthMiddleware checks for valid admin JWT token
func AdminAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &models.AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*models.AdminClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// Check if user role is admin
		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("adminID", claims.ID)
		c.Set("adminEmail", claims.Email)
		c.Next()
	}
}

// ReceptionistAuthMiddleware checks for valid receptionist JWT token
func ReceptionistAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &models.ReceptionistClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*models.ReceptionistClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// Check if user role is receptionist
		if claims.Role != "receptionist" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("receptionistID", claims.ID)
		c.Set("receptionistEmail", claims.Email)
		c.Next()
	}
}

// DoctorAuthMiddleware checks for valid doctor JWT token
func DoctorAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &models.DoctorClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*models.DoctorClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// Check if user role is doctor
		if claims.Role != "doctor" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("doctorID", claims.ID)
		c.Set("doctorEmail", claims.Email)
		c.Set("doctorName", claims.Name)
		c.Next()
	}
}

// GenerateDoctorToken generates a JWT token for a doctor
func GenerateDoctorToken(cfg *config.Config, doctor *models.DoctorUser) (string, error) {
	parsedDuration, err := time.ParseDuration(cfg.JWTExpire)
	if err != nil {
		parsedDuration = 7 * 24 * time.Hour
	}

	claims := models.DoctorClaims{
		ID:    doctor.ID,
		Email: doctor.Email,
		Name:  doctor.Name,
		Role:  "doctor",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(parsedDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}
