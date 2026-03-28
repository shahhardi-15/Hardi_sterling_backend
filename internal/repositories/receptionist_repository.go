package repositories

import (
	"database/sql"
	"errors"
	"time"

	"sterling-hms-backend/internal/models"
)

type ReceptionistRepository struct {
	db *sql.DB
}

func NewReceptionistRepository(db *sql.DB) *ReceptionistRepository {
	return &ReceptionistRepository{db: db}
}

// FindByEmail gets receptionist user by email
func (r *ReceptionistRepository) FindByEmail(email string) (*models.ReceptionistUser, error) {
	receptionist := &models.ReceptionistUser{}
	err := r.db.QueryRow(
		`SELECT id, email, password_hash, name, phone, department, role, created_at, updated_at, is_active
		 FROM receptionist_users WHERE email = $1`,
		email,
	).Scan(&receptionist.ID, &receptionist.Email, &receptionist.PasswordHash, &receptionist.Name,
		&receptionist.Phone, &receptionist.Department, &receptionist.Role, &receptionist.CreatedAt,
		&receptionist.UpdatedAt, &receptionist.IsActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("receptionist not found")
		}
		return nil, err
	}

	return receptionist, nil
}

// FindByID gets receptionist user by ID
func (r *ReceptionistRepository) FindByID(id int) (*models.ReceptionistUser, error) {
	receptionist := &models.ReceptionistUser{}
	err := r.db.QueryRow(
		`SELECT id, email, name, phone, department, role, created_at, updated_at, is_active
		 FROM receptionist_users WHERE id = $1`,
		id,
	).Scan(&receptionist.ID, &receptionist.Email, &receptionist.Name, &receptionist.Phone,
		&receptionist.Department, &receptionist.Role, &receptionist.CreatedAt, &receptionist.UpdatedAt,
		&receptionist.IsActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("receptionist not found")
		}
		return nil, err
	}

	return receptionist, nil
}

// CreatePatient creates a new patient user and patient record
func (r *ReceptionistRepository) CreatePatient(firstName, lastName, email, hashedPassword, phone string,
	dateOfBirth *time.Time, gender, bloodType, address, city, state, postalCode, country,
	allergies, medicalConditions, currentMedications, emergencyContactName, emergencyContactPhone string) (*models.PatientRecord, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create user
	var userID int
	err = tx.QueryRow(
		`INSERT INTO users (first_name, last_name, email, password, created_at, updated_at, is_active)
		 VALUES ($1, $2, $3, $4, NOW(), NOW(), true)
		 RETURNING id`,
		firstName, lastName, email, hashedPassword,
	).Scan(&userID)
	if err != nil {
		return nil, err
	}

	// Create patient record
	patientRecord := &models.PatientRecord{
		UserID:                   userID,
		FirstName:                firstName,
		LastName:                 lastName,
		Email:                    email,
		Phone:                    phone,
		DateOfBirth:              dateOfBirth,
		Gender:                   gender,
		BloodType:                bloodType,
		Address:                  address,
		City:                     city,
		State:                    state,
		PostalCode:               postalCode,
		Country:                  country,
		Allergies:                allergies,
		MedicalConditions:        medicalConditions,
		CurrentMedications:       currentMedications,
		EmergencyContactName:     emergencyContactName,
		EmergencyContactPhone:    emergencyContactPhone,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	err = tx.QueryRow(
		`INSERT INTO patient_records (user_id, first_name, last_name, email, phone, date_of_birth, gender, blood_type,
		 address, city, state, postal_code, country, allergies, medical_conditions, current_medications,
		 emergency_contact_name, emergency_contact_phone, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, NOW(), NOW())
		 RETURNING id`,
		userID, firstName, lastName, email, phone, dateOfBirth, gender, bloodType, address, city,
		state, postalCode, country, allergies, medicalConditions, currentMedications,
		emergencyContactName, emergencyContactPhone,
	).Scan(&patientRecord.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return patientRecord, nil
}

// GetPatientRecord gets patient record by user ID
func (r *ReceptionistRepository) GetPatientRecord(userID int) (*models.PatientRecord, error) {
	record := &models.PatientRecord{}
	err := r.db.QueryRow(
		`SELECT id, user_id, first_name, last_name, email, phone, date_of_birth, gender, blood_type,
		 address, city, state, postal_code, country, allergies, medical_conditions, current_medications,
		 emergency_contact_name, emergency_contact_phone, created_at, updated_at
		 FROM patient_records WHERE user_id = $1`,
		userID,
	).Scan(&record.ID, &record.UserID, &record.FirstName, &record.LastName, &record.Email, &record.Phone,
		&record.DateOfBirth, &record.Gender, &record.BloodType, &record.Address, &record.City, &record.State,
		&record.PostalCode, &record.Country, &record.Allergies, &record.MedicalConditions,
		&record.CurrentMedications, &record.EmergencyContactName, &record.EmergencyContactPhone,
		&record.CreatedAt, &record.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("patient record not found")
		}
		return nil, err
	}

	return record, nil
}

// UpdatePatientRecord updates patient medical information
func (r *ReceptionistRepository) UpdatePatientRecord(userID int, record *models.PatientRecord) error {
	_, err := r.db.Exec(
		`UPDATE patient_records SET
		 first_name = $2, last_name = $3, phone = $4, date_of_birth = $5, gender = $6, blood_type = $7,
		 address = $8, city = $9, state = $10, postal_code = $11, country = $12, allergies = $13,
		 medical_conditions = $14, current_medications = $15, emergency_contact_name = $16,
		 emergency_contact_phone = $17, updated_at = NOW()
		 WHERE user_id = $1`,
		userID, record.FirstName, record.LastName, record.Phone, record.DateOfBirth, record.Gender,
		record.BloodType, record.Address, record.City, record.State, record.PostalCode, record.Country,
		record.Allergies, record.MedicalConditions, record.CurrentMedications, record.EmergencyContactName,
		record.EmergencyContactPhone,
	)
	return err
}

// BookAppointmentByReceptionist books appointment on behalf of patient
func (r *ReceptionistRepository) BookAppointmentByReceptionist(patientID, doctorID int, appointmentDate, timeSlot, reason string,
	patientFirstName, patientLastName, patientEmail, patientPhone string) (*models.Appointment, error) {

	appointment := &models.Appointment{}
	err := r.db.QueryRow(
		`INSERT INTO appointments (patient_id, doctor_id, appointment_date, time_slot, reason, status,
		 approval_status, patient_first_name, patient_last_name, patient_email, patient_phone, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, 'scheduled', 'pending', $6, $7, $8, $9, NOW(), NOW())
		 RETURNING id, patient_id, doctor_id, appointment_date, time_slot, reason, status, approval_status,
		 notes, approved_by, approved_at, rejection_reason, patient_first_name, patient_last_name,
		 patient_email, patient_phone, created_at, updated_at`,
		patientID, doctorID, appointmentDate, timeSlot, reason, patientFirstName, patientLastName,
		patientEmail, patientPhone,
	).Scan(&appointment.ID, &appointment.PatientID, &appointment.DoctorID, &appointment.AppointmentDate,
		&appointment.TimeSlot, &appointment.Reason, &appointment.Status, &appointment.ApprovalStatus,
		&appointment.Notes, &appointment.ApprovedBy, &appointment.ApprovedAt, &appointment.RejectionReason,
		&appointment.PatientFirstName, &appointment.PatientLastName, &appointment.PatientEmail,
		&appointment.PatientPhone, &appointment.CreatedAt, &appointment.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return appointment, nil
}

// GetPendingAppointments gets all pending appointments for approval
func (r *ReceptionistRepository) GetPendingAppointments(page, limit int) ([]models.Appointment, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM appointments WHERE approval_status = 'pending'`,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get appointments with doctor info
	rows, err := r.db.Query(
		`SELECT a.id, a.patient_id, a.doctor_id, a.appointment_date, a.time_slot, a.reason, a.status,
		 a.approval_status, a.notes, a.approved_by, a.approved_at, a.rejection_reason, a.patient_first_name,
		 a.patient_last_name, a.patient_email, a.patient_phone, a.created_at, a.updated_at,
		 d.id, d.name, d.specialization, d.email, d.phone, d.experience_years, d.qualification, d.address
		 FROM appointments a
		 LEFT JOIN doctors d ON a.doctor_id = d.id
		 WHERE a.approval_status = 'pending'
		 ORDER BY a.created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		apt := models.Appointment{}
		doc := &models.Doctor{}
		err := rows.Scan(&apt.ID, &apt.PatientID, &apt.DoctorID, &apt.AppointmentDate, &apt.TimeSlot,
			&apt.Reason, &apt.Status, &apt.ApprovalStatus, &apt.Notes, &apt.ApprovedBy, &apt.ApprovedAt,
			&apt.RejectionReason, &apt.PatientFirstName, &apt.PatientLastName, &apt.PatientEmail,
			&apt.PatientPhone, &apt.CreatedAt, &apt.UpdatedAt, &doc.ID, &doc.Name, &doc.Specialization,
			&doc.Email, &doc.Phone, &doc.ExperienceYears, &doc.Qualification, &doc.Address)
		if err != nil {
			return nil, 0, err
		}
		apt.Doctor = doc
		appointments = append(appointments, apt)
	}

	return appointments, total, nil
}

// ApproveAppointment approves a pending appointment
func (r *ReceptionistRepository) ApproveAppointment(appointmentID int, receptionistID int) error {
	_, err := r.db.Exec(
		`UPDATE appointments SET approval_status = 'approved', approved_by = $2, approved_at = NOW(), updated_at = NOW()
		 WHERE id = $1`,
		appointmentID, receptionistID,
	)
	return err
}

// RejectAppointment rejects a pending appointment
func (r *ReceptionistRepository) RejectAppointment(appointmentID int, receptionistID int, reason string) error {
	_, err := r.db.Exec(
		`UPDATE appointments SET approval_status = 'rejected', approved_by = $2, approved_at = NOW(),
		 rejection_reason = $3, updated_at = NOW()
		 WHERE id = $1`,
		appointmentID, receptionistID, reason,
	)
	return err
}

// GetAppointmentByID gets a single appointment by ID
func (r *ReceptionistRepository) GetAppointmentByID(id int) (*models.Appointment, error) {
	appointment := &models.Appointment{}
	err := r.db.QueryRow(
		`SELECT id, patient_id, doctor_id, appointment_date, time_slot, reason, status, approval_status,
		 notes, approved_by, approved_at, rejection_reason, patient_first_name, patient_last_name,
		 patient_email, patient_phone, created_at, updated_at
		 FROM appointments WHERE id = $1`,
		id,
	).Scan(&appointment.ID, &appointment.PatientID, &appointment.DoctorID, &appointment.AppointmentDate,
		&appointment.TimeSlot, &appointment.Reason, &appointment.Status, &appointment.ApprovalStatus,
		&appointment.Notes, &appointment.ApprovedBy, &appointment.ApprovedAt, &appointment.RejectionReason,
		&appointment.PatientFirstName, &appointment.PatientLastName, &appointment.PatientEmail,
		&appointment.PatientPhone, &appointment.CreatedAt, &appointment.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}

	return appointment, nil
}

// GetAllPatients gets all registered patients
func (r *ReceptionistRepository) GetAllPatients(page, limit int) ([]models.PatientRecord, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM patient_records`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get patients
	rows, err := r.db.Query(
		`SELECT id, user_id, first_name, last_name, email, phone, date_of_birth, gender, blood_type,
		 address, city, state, postal_code, country, allergies, medical_conditions, current_medications,
		 emergency_contact_name, emergency_contact_phone, created_at, updated_at
		 FROM patient_records
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var patients []models.PatientRecord
	for rows.Next() {
		record := models.PatientRecord{}
		err := rows.Scan(&record.ID, &record.UserID, &record.FirstName, &record.LastName, &record.Email,
			&record.Phone, &record.DateOfBirth, &record.Gender, &record.BloodType, &record.Address,
			&record.City, &record.State, &record.PostalCode, &record.Country, &record.Allergies,
			&record.MedicalConditions, &record.CurrentMedications, &record.EmergencyContactName,
			&record.EmergencyContactPhone, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		patients = append(patients, record)
	}

	return patients, total, nil
}

// GetReceptionistDashboardStats gets dashboard statistics
func (r *ReceptionistRepository) GetReceptionistDashboardStats() (*models.ReceptionistDashboardStats, error) {
	stats := &models.ReceptionistDashboardStats{}

	// Count total patients
	err := r.db.QueryRow(`SELECT COUNT(*) FROM patient_records`).Scan(&stats.TotalPatients)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Count pending appointments
	err = r.db.QueryRow(`SELECT COUNT(*) FROM appointments WHERE approval_status = 'pending'`).Scan(&stats.PendingAppointments)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Count approved appointments
	err = r.db.QueryRow(`SELECT COUNT(*) FROM appointments WHERE approval_status = 'approved'`).Scan(&stats.ApprovedAppointments)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Count rejected appointments
	err = r.db.QueryRow(`SELECT COUNT(*) FROM appointments WHERE approval_status = 'rejected'`).Scan(&stats.RejectedAppointments)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return stats, nil
}
