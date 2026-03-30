package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"sterling-hms-backend/internal/models"
	"time"
)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

// ListPatients retrieves all patients with pagination and search
func (r *PatientRepository) ListPatients(page, limit int, search string) ([]*models.PatientRecord, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT pr.id, pr.user_id, pr.first_name, pr.last_name, pr.email, pr.phone,
		       pr.date_of_birth, pr.gender, pr.blood_type, pr.address, pr.city, pr.state,
		       pr.postal_code, pr.country, pr.allergies, pr.medical_conditions, 
		       pr.current_medications, pr.emergency_contact_name, pr.emergency_contact_phone,
		       pr.created_at, pr.updated_at, u.uhid, u.is_active,
		       COALESCE(pr.registration_date, pr.created_at) as registration_date
		FROM patient_records pr
		LEFT JOIN users u ON pr.user_id = u.id
	`

	args := []interface{}{}

	// Add search condition if provided
	if search != "" {
		query += ` WHERE pr.first_name ILIKE $1 OR pr.phone ILIKE $1 OR u.uhid ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY pr.created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var patients []*models.PatientRecord
	for rows.Next() {
		patient := &models.PatientRecord{}
		err := rows.Scan(
			&patient.ID, &patient.UserID, &patient.FirstName, &patient.LastName, &patient.Email,
			&patient.Phone, &patient.DateOfBirth, &patient.Gender, &patient.BloodType, &patient.Address,
			&patient.City, &patient.State, &patient.PostalCode, &patient.Country, &patient.Allergies,
			&patient.MedicalConditions, &patient.CurrentMedications, &patient.EmergencyContactName,
			&patient.EmergencyContactPhone, &patient.CreatedAt, &patient.UpdatedAt, nil, nil, nil,
		)
		// For now, ignore UHID and is_active, we'll fetch those from users table separately
		if err != nil {
			continue
		}
		patients = append(patients, patient)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM patient_records pr LEFT JOIN users u ON pr.user_id = u.id`
	if search != "" {
		countQuery += ` WHERE pr.first_name ILIKE $1 OR pr.phone ILIKE $1 OR u.uhid ILIKE $1`
	}

	var total int
	var countArgs []interface{}
	if search != "" {
		countArgs = append(countArgs, "%"+search+"%")
	}

	err = r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

// GetPatientWithUser retrieves a patient with user details
func (r *PatientRepository) GetPatientWithUser(patientID int) (*models.PatientRecord, string, bool, error) {
	patient := &models.PatientRecord{}
	var uhid *string
	var isActive bool

	err := r.db.QueryRow(`
		SELECT pr.id, pr.user_id, pr.first_name, pr.last_name, pr.email, pr.phone,
		       pr.date_of_birth, pr.gender, pr.blood_type, pr.address, pr.city, pr.state,
		       pr.postal_code, pr.country, pr.allergies, pr.medical_conditions,
		       pr.current_medications, pr.emergency_contact_name, pr.emergency_contact_phone,
		       pr.created_at, pr.updated_at, u.uhid, u.is_active,
		       COALESCE(pr.registration_date, pr.created_at) as registration_date
		FROM patient_records pr
		LEFT JOIN users u ON pr.user_id = u.id
		WHERE pr.id = $1`,
		patientID,
	).Scan(
		&patient.ID, &patient.UserID, &patient.FirstName, &patient.LastName, &patient.Email,
		&patient.Phone, &patient.DateOfBirth, &patient.Gender, &patient.BloodType, &patient.Address,
		&patient.City, &patient.State, &patient.PostalCode, &patient.Country, &patient.Allergies,
		&patient.MedicalConditions, &patient.CurrentMedications, &patient.EmergencyContactName,
		&patient.EmergencyContactPhone, &patient.CreatedAt, &patient.UpdatedAt, &uhid, &isActive, nil,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", false, errors.New("patient not found")
		}
		return nil, "", false, err
	}

	uhidStr := ""
	if uhid != nil {
		uhidStr = *uhid
	}

	return patient, uhidStr, isActive, nil
}

// CreatePatient creates a new patient record
// Note: User creation should be done in the handler with transaction
func (r *PatientRepository) CreatePatient(
	userID int,
	firstName, lastName, email, phone string,
	dateOfBirth *time.Time,
	gender, bloodType, address, city, state, postalCode, country,
	allergies, medicalConditions, currentMedications, emergencyContactName, emergencyContactPhone string,
) (*models.PatientRecord, error) {

	patient := &models.PatientRecord{}

	err := r.db.QueryRow(`
		INSERT INTO patient_records (
			user_id, first_name, last_name, email, phone, date_of_birth, gender,
			blood_type, address, city, state, postal_code, country, allergies,
			medical_conditions, current_medications, emergency_contact_name, emergency_contact_phone,
			created_at, updated_at, registration_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, NOW(), NOW(), NOW())
		RETURNING id, user_id, first_name, last_name, email, phone, date_of_birth, gender,
		          blood_type, address, city, state, postal_code, country, allergies,
		          medical_conditions, current_medications, emergency_contact_name, emergency_contact_phone,
		          created_at, updated_at`,
		userID, firstName, lastName, email, phone, dateOfBirth, gender,
		bloodType, address, city, state, postalCode, country, allergies,
		medicalConditions, currentMedications, emergencyContactName, emergencyContactPhone,
	).Scan(
		&patient.ID, &patient.UserID, &patient.FirstName, &patient.LastName, &patient.Email,
		&patient.Phone, &patient.DateOfBirth, &patient.Gender, &patient.BloodType, &patient.Address,
		&patient.City, &patient.State, &patient.PostalCode, &patient.Country, &patient.Allergies,
		&patient.MedicalConditions, &patient.CurrentMedications, &patient.EmergencyContactName,
		&patient.EmergencyContactPhone, &patient.CreatedAt, &patient.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return patient, nil
}

// UpdatePatient updates patient information
func (r *PatientRepository) UpdatePatient(
	patientID int,
	firstName, lastName, phone string,
	dateOfBirth *time.Time,
	gender, address string,
) error {

	_, err := r.db.Exec(`
		UPDATE patient_records
		SET first_name = $1, last_name = $2, phone = $3, date_of_birth = $4,
		    gender = $5, address = $6, updated_at = NOW()
		WHERE id = $7`,
		firstName, lastName, phone, dateOfBirth, gender, address, patientID,
	)

	return err
}

// UpdatePatientStatus updates the is_active status of a patient
func (r *PatientRepository) UpdatePatientStatus(userID int, isActive bool) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET is_active = $1, updated_at = NOW()
		WHERE id = $2`,
		isActive, userID,
	)
	return err
}

// GenerateUHID generates a new UHID for a patient
func (r *PatientRepository) GenerateUHID() (string, error) {
	year := time.Now().Year()

	var maxSeq int
	err := r.db.QueryRow(
		`SELECT COALESCE(MAX(CAST(SUBSTRING(uhid, 10) AS INTEGER)), 0) 
		 FROM users WHERE uhid LIKE $1`,
		fmt.Sprintf("SHMS-%d-%%", year),
	).Scan(&maxSeq)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	nextSeq := maxSeq + 1
	uhid := fmt.Sprintf("SHMS-%d-%05d", year, nextSeq)

	return uhid, nil
}
