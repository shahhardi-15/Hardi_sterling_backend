package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// RateLimiter tracks requests for rate limiting
type RateLimiter struct {
	maxRequests int           // Maximum requests allowed
	windowSize  time.Duration // Time window for the requests
}

// NewRateLimiter creates a new rate limiter
// maxRequests: limit requests within windowSize
func NewRateLimiter(maxRequests int, windowSize time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		windowSize:  windowSize,
	}
}

// IsRateLimited checks if request exceeds rate limit
// requestCount: number of requests in the current window
func (rl *RateLimiter) IsRateLimited(requestCount int) bool {
	return requestCount >= rl.maxRequests
}

// GetRetryAfterSeconds returns seconds to wait before next request
func (rl *RateLimiter) GetRetryAfterSeconds(requestCount int) int {
	if requestCount < rl.maxRequests {
		return 0
	}
	return int(rl.windowSize.Seconds())
}

// TokenGenerator generates secure tokens
type TokenGenerator struct {
	length int
}

// NewTokenGenerator creates a new token generator
func NewTokenGenerator(length int) *TokenGenerator {
	if length <= 0 {
		length = 32
	}
	return &TokenGenerator{length: length}
}

// GenerateToken generates a random hex token
func (tg *TokenGenerator) GenerateToken() (string, error) {
	bytes := make([]byte, tg.length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// PasswordResetFlow encapsulates password reset flow configuration
type PasswordResetFlow struct {
	ResetTokenExpiry   time.Duration
	MaxResetRequests   int // Per hour
	ResetRequestWindow time.Duration
}

// NewPasswordResetFlow creates default password reset flow configuration
func NewPasswordResetFlow() *PasswordResetFlow {
	return &PasswordResetFlow{
		ResetTokenExpiry:   1 * time.Hour,
		MaxResetRequests:   5,
		ResetRequestWindow: 1 * time.Hour,
	}
}

// GetResetTokenExpiryTime returns when reset token will expire
func (prf *PasswordResetFlow) GetResetTokenExpiryTime() time.Time {
	return time.Now().Add(prf.ResetTokenExpiry)
}

// IsWithinResetRequestWindow checks if within reset request time window
func (prf *PasswordResetFlow) IsWithinResetRequestWindow(timestamp time.Time) bool {
	return time.Since(timestamp) < prf.ResetRequestWindow
}
