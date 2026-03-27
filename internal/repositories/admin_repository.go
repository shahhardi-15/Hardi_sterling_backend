package repositories

import (
	"database/sql"
	"errors"
	"sterling-hms-backend/internal/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// FindByEmail gets admin user by email
func (r *AdminRepository) FindByEmail(email string) (*models.AdminUser, error) {
	admin := &models.AdminUser{}
	err := r.db.QueryRow(
		`SELECT id, email, password_hash, name, role, created_at, updated_at, is_active
		 FROM admin_users WHERE email = $1`,
		email,
	).Scan(&admin.ID, &admin.Email, &admin.PasswordHash, &admin.Name, &admin.Role, &admin.CreatedAt, &admin.UpdatedAt, &admin.IsActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}

	return admin, nil
}

// FindByID gets admin user by ID
func (r *AdminRepository) FindByID(id int) (*models.AdminUser, error) {
	admin := &models.AdminUser{}
	err := r.db.QueryRow(
		`SELECT id, email, name, role, created_at, updated_at, is_active
		 FROM admin_users WHERE id = $1`,
		id,
	).Scan(&admin.ID, &admin.Email, &admin.Name, &admin.Role, &admin.CreatedAt, &admin.UpdatedAt, &admin.IsActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}

	return admin, nil
}

// GetDashboardStats returns aggregated stats for the dashboard
func (r *AdminRepository) GetDashboardStats() (*models.AdminDashboardStats, error) {
	stats := &models.AdminDashboardStats{}

	// Count total patients
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&stats.TotalPatients)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Count total appointments
	err = r.db.QueryRow(`SELECT COUNT(*) FROM appointments`).Scan(&stats.TotalAppointments)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Count total doctors
	err = r.db.QueryRow(`SELECT COUNT(*) FROM doctors`).Scan(&stats.TotalDoctors)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Count total staff (admin users)
	err = r.db.QueryRow(`SELECT COUNT(*) FROM admin_users WHERE is_active = true`).Scan(&stats.TotalStaff)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return stats, nil
}

// LogAdminAction logs admin actions for audit trail
func (r *AdminRepository) LogAdminAction(adminID int, action string, resourceType string, resourceID *int, details string, ipAddress string, userAgent string) error {
	_, err := r.db.Exec(
		`INSERT INTO admin_audit_logs (admin_id, action, resource_type, resource_id, details, ip_address, user_agent, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
		adminID, action, resourceType, resourceID, details, ipAddress, userAgent,
	)
	return err
}

// EmailExists checks if an admin email exists
func (r *AdminRepository) EmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM admin_users WHERE email = $1`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
