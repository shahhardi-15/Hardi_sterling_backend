package utils

import (
	"errors"
	"unicode"
)

// PasswordValidator validates password strength requirements
type PasswordValidator struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumber    bool
	RequireSpecial   bool
}

// NewPasswordValidator creates a new password validator with default requirements
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumber:    true,
		RequireSpecial:   true,
	}
}

// ValidatePassword validates password against all requirements
func (pv *PasswordValidator) ValidatePassword(password string) error {
	// Check minimum length
	if len(password) < pv.MinLength {
		return errors.New("password must be at least 8 characters long")
	}

	// Check maximum length (prevent DoS)
	if len(password) > 128 {
		return errors.New("password must not exceed 128 characters")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	// Check each character
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case isSpecialChar(char):
			hasSpecial = true
		}
	}

	// Validate requirements
	if pv.RequireUppercase && !hasUpper {
		return errors.New("password must contain at least one uppercase letter (A-Z)")
	}

	if pv.RequireLowercase && !hasLower {
		return errors.New("password must contain at least one lowercase letter (a-z)")
	}

	if pv.RequireNumber && !hasNumber {
		return errors.New("password must contain at least one number (0-9)")
	}

	if pv.RequireSpecial && !hasSpecial {
		return errors.New("password must contain at least one special character (!@#$%^&*)")
	}

	// Check for commonly used weak passwords
	if isCommonPassword(password) {
		return errors.New("password is too common or weak")
	}

	return nil
}

// isSpecialChar checks if character is a special character
func isSpecialChar(char rune) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?/~`"
	for _, s := range specialChars {
		if char == s {
			return true
		}
	}
	return false
}

// isCommonPassword checks against common weak passwords
func isCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "123456", "12345678", "qwerty", "abc123",
		"111111", "1234567", "password123", "admin", "letmein",
		"welcome", "monkey", "dragon", "master", "sunshine",
	}

	for _, common := range commonPasswords {
		if password == common {
			return true
		}
	}

	return false
}

// ValidatePasswordMatch validates that two passwords match
func ValidatePasswordMatch(password, confirmPassword string) error {
	if password != confirmPassword {
		return errors.New("passwords do not match")
	}
	return nil
}

// GetPasswordStrength returns a strength score (1-5) for a password
func GetPasswordStrength(password string) int {
	score := 0

	// Length score
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}

	// Character variety score
	hasUpper, hasLower, hasNumber, hasSpecial := false, false, false, false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
		if isSpecialChar(char) {
			hasSpecial = true
		}
	}

	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score++
	}

	// Normalize to 1-5 range
	if score > 5 {
		score = 5
	}
	if score < 1 {
		score = 1
	}

	return score
}

// GeneratePasswordRequirements returns a formatted string of password requirements
func GeneratePasswordRequirements() string {
	return `Password must contain:
- At least 8 characters
- At least one uppercase letter (A-Z)
- At least one lowercase letter (a-z)
- At least one number (0-9)
- At least one special character (!@#$%^&*)`
}
