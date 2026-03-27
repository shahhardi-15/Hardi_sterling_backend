package utils

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuditLog handles audit logging for security events
type AuditLog struct {
	// Empty struct - methods will handle logging to database
}

// NewAuditLog creates a new audit logger
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// GetIPAddress extracts client IP from request
func (al *AuditLog) GetIPAddress(c *gin.Context) string {
	// Check for X-Forwarded-For header first (for proxies/load balancers)
	forwardedFor := c.Request.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// X-Forwarded-For can have multiple IPs, take the first one
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check for X-Real-IP header
	realIP := c.Request.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}

	return ip
}

// GetUserAgent extracts user agent from request
func (al *AuditLog) GetUserAgent(c *gin.Context) string {
	return c.Request.Header.Get("User-Agent")
}

// LogAction logs an action with all relevant information
// This is typically called from the repository layer
func (al *AuditLog) LogAction(c *gin.Context, email string, userID *int, action string, success bool, errorMessage *string) {
	ipAddress := al.GetIPAddress(c)
	userAgent := al.GetUserAgent(c)

	// Log to a file or service for monitoring
	logMessage := "AUDIT: action=" + action + " email=" + email + " ip=" + ipAddress + " success="
	if success {
		logMessage += "true"
	} else {
		logMessage += "false"
	}
	if errorMessage != nil {
		logMessage += " error=" + *errorMessage
	}
	logMessage += " userAgent=" + userAgent

	// In production, you'd send this to a logging service (e.g., ELK, Datadog, etc.)
	// For now, just logging to stdout (structured logging can be added)
	// log.Println(logMessage)
}

// SecurityConstants defines common security scenarios
const (
	ActionForgotPasswordRequest = "forgot_password_request"
	ActionPasswordResetSuccess  = "reset_success"
	ActionPasswordResetFailed   = "reset_failed"
)

// ErrorMessages provides generic error messages to prevent user enumeration
var ErrorMessages = struct {
	GenericError     string
	InvalidEmail     string
	InvalidToken     string
	InvalidPassword  string
	WeakPassword     string
	PasswordMismatch string
	LockedOut        string
	RateLimited      string
	ServerError      string
}{
	GenericError:     "An error occurred. Please try again.",
	InvalidEmail:     "If an account exists with this email, you will receive a password reset link.",
	InvalidToken:     "Invalid or expired reset link. Please request a new password reset.",
	InvalidPassword:  "Invalid password format. Please ensure it meets all requirements.",
	WeakPassword:     "Password is too weak. Please use a stronger password.",
	PasswordMismatch: "Passwords do not match.",
	LockedOut:        "Too many failed attempts. Please try again later.",
	RateLimited:      "Too many requests. Please try again later.",
	ServerError:      "A server error occurred. Our team has been notified. Please try again later.",
}

// SecurityHeaders returns middleware to add security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
