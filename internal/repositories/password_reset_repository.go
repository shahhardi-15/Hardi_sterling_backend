package repositories

import (
	"database/sql"
	"errors"
	"sterling-hms-backend/internal/models"
	"time"
)

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// CreatePasswordResetToken creates a new password reset token
func (r *PasswordResetRepository) CreatePasswordResetToken(userID int, tokenHash string, expiresAt time.Time) (*models.PasswordResetToken, error) {
	token := &models.PasswordResetToken{}
	err := r.db.QueryRow(
		`INSERT INTO password_reset_tokens (user_id, token_hash, expires_at, is_used, created_at)
		 VALUES ($1, $2, $3, false, NOW())
		 RETURNING id, user_id, token_hash, expires_at, used_at, is_used, created_at`,
		userID, tokenHash, expiresAt,
	).Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.IsUsed, &token.CreatedAt)

	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetPasswordResetTokenByHash gets token by hash
func (r *PasswordResetRepository) GetPasswordResetTokenByHash(tokenHash string) (*models.PasswordResetToken, error) {
	token := &models.PasswordResetToken{}
	err := r.db.QueryRow(
		`SELECT id, user_id, token_hash, expires_at, used_at, is_used, created_at
		 FROM password_reset_tokens
		 WHERE token_hash = $1 AND is_used = false AND expires_at > NOW()`,
		tokenHash,
	).Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.IsUsed, &token.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid or expired reset token")
		}
		return nil, err
	}

	return token, nil
}

// MarkResetTokenAsUsed marks token as used after successful password reset
func (r *PasswordResetRepository) MarkResetTokenAsUsed(tokenID int) error {
	_, err := r.db.Exec(
		`UPDATE password_reset_tokens SET is_used = true, used_at = NOW() WHERE id = $1`,
		tokenID,
	)
	return err
}

// UpdateUserPassword updates user password
func (r *PasswordResetRepository) UpdateUserPassword(userID int, hashedPassword string) error {
	_, err := r.db.Exec(
		`UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`,
		hashedPassword, userID,
	)
	return err
}

// CreatePasswordResetLog logs password reset actions
func (r *PasswordResetRepository) CreatePasswordResetLog(userID *int, email, action, ipAddress, userAgent string, success bool, errorMessage *string) error {
	_, err := r.db.Exec(
		`INSERT INTO password_reset_logs (user_id, email, action, ip_address, user_agent, success, error_message, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
		userID, email, action, ipAddress, userAgent, success, errorMessage,
	)
	return err
}

// CountResetAttempts counts password reset attempts by email in the last hour
func (r *PasswordResetRepository) CountResetAttempts(email string, hours int) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM password_reset_logs
		 WHERE email = $1 AND action = 'forgot_password_request' AND created_at > NOW() - INTERVAL '1 hour' * $2`,
		email, hours,
	).Scan(&count)
	return count, err
}

// InvalidateAllUserTokens invalidates all password reset tokens for a user
func (r *PasswordResetRepository) InvalidateAllUserTokens(userID int) error {
	_, err := r.db.Exec(
		`UPDATE password_reset_tokens SET is_used = true WHERE user_id = $1 AND is_used = false`,
		userID,
	)
	return err
}

// CleanupExpiredTokens cleans up expired password reset tokens (maintenance task)
func (r *PasswordResetRepository) CleanupExpiredTokens() error {
	_, err := r.db.Exec(
		`DELETE FROM password_reset_tokens WHERE expires_at < NOW()`,
	)
	return err
}
