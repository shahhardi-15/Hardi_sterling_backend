package repositories

import (
	"database/sql"
	"errors"
	"sterling-hms-backend/internal/models"
)

type DoctorRepository struct {
	db *sql.DB
}

func NewDoctorRepository(db *sql.DB) *DoctorRepository {
	return &DoctorRepository{db: db}
}

// FindByEmail retrieves a doctor by email
func (r *DoctorRepository) FindByEmail(email string) (*models.DoctorUser, error) {
	doctor := &models.DoctorUser{}
	err := r.db.QueryRow(
		`SELECT id, email, name, password_hash, specialization, phone, role, is_active, created_at, updated_at
		 FROM doctor_users WHERE email = $1`,
		email,
	).Scan(&doctor.ID, &doctor.Email, &doctor.Name, &doctor.PasswordHash, &doctor.Specialization,
		&doctor.Phone, &doctor.Role, &doctor.IsActive, &doctor.CreatedAt, &doctor.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("doctor not found")
		}
		return nil, err
	}

	return doctor, nil
}

// FindByID retrieves a doctor by ID
func (r *DoctorRepository) FindByID(id int) (*models.DoctorUser, error) {
	doctor := &models.DoctorUser{}
	err := r.db.QueryRow(
		`SELECT id, email, name, password_hash, specialization, phone, role, is_active, created_at, updated_at
		 FROM doctor_users WHERE id = $1`,
		id,
	).Scan(&doctor.ID, &doctor.Email, &doctor.Name, &doctor.PasswordHash, &doctor.Specialization,
		&doctor.Phone, &doctor.Role, &doctor.IsActive, &doctor.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("doctor not found")
		}
		return nil, err
	}

	return doctor, nil
}

// GetAssignedPatients retrieves all patients assigned to a doctor
func (r *DoctorRepository) GetAssignedPatients(doctorID int) ([]models.PatientInfo, error) {
	query := `
		SELECT DISTINCT u.id, CONCAT(u.first_name, ' ', u.last_name) as name, 
		       COALESCE(pr.phone, '') as phone,
		       u.email,
		       COALESCE(pr.medical_conditions, '') as medical_history,
		       COALESCE(pr.address, '') as address,
		       COALESCE(pr.blood_type, '') as blood_type
		FROM doctor_patient_assignment dpa
		JOIN users u ON dpa.patient_id = u.id
		LEFT JOIN patient_records pr ON u.id = pr.user_id
		WHERE dpa.doctor_id = $1 AND dpa.is_active = true
		ORDER BY u.first_name, u.last_name
	`

	rows, err := r.db.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.PatientInfo
	for rows.Next() {
		var patient models.PatientInfo
		err := rows.Scan(
			&patient.ID, &patient.Name, &patient.Phone, &patient.Email,
			&patient.MedicalHistory, &patient.Address, &patient.BloodType,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, rows.Err()
}

// GetAppointments retrieves all appointments for a doctor
func (r *DoctorRepository) GetAppointments(doctorID int) ([]models.AppointmentInfoForDoctor, error) {
	query := `
		SELECT a.id, CONCAT(u.first_name, ' ', u.last_name) as patient_name, a.patient_id,
		       a.appointment_date, a.time_slot, a.reason, a.status, a.notes,
		       u.phone, u.email
		FROM appointments a
		JOIN users u ON a.patient_id = u.id
		WHERE a.doctor_id = $1
		ORDER BY a.appointment_date DESC, a.time_slot DESC
	`

	rows, err := r.db.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.AppointmentInfoForDoctor
	for rows.Next() {
		var apt models.AppointmentInfoForDoctor
		err := rows.Scan(
			&apt.ID, &apt.PatientName, &apt.PatientID, &apt.AppointmentDate, &apt.TimeSlot,
			&apt.Reason, &apt.Status, &apt.Notes, &apt.PatientPhone, &apt.PatientEmail,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, apt)
	}

	return appointments, rows.Err()
}

// CheckAppointmentOwnership checks if an appointment belongs to a doctor
func (r *DoctorRepository) CheckAppointmentOwnership(appointmentID int, doctorID int) (bool, error) {
	query := `SELECT COUNT(*) FROM appointments WHERE id = $1 AND doctor_id = $2`

	var count int
	err := r.db.QueryRow(query, appointmentID, doctorID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateAppointmentStatus updates the status of an appointment
func (r *DoctorRepository) UpdateAppointmentStatus(appointmentID int, status string, notes string) error {
	query := `
		UPDATE appointments 
		SET status = $1, notes = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, notes, appointmentID)
	return err
}

// GetAppointmentByID retrieves a single appointment
func (r *DoctorRepository) GetAppointmentByID(appointmentID int) (*models.AppointmentInfoForDoctor, error) {
	query := `
		SELECT a.id, CONCAT(u.first_name, ' ', u.last_name) as patient_name, a.patient_id,
		       a.appointment_date, a.time_slot, a.reason, a.status, a.notes,
		       u.phone, u.email
		FROM appointments a
		JOIN users u ON a.patient_id = u.id
		WHERE a.id = $1
	`

	var apt models.AppointmentInfoForDoctor
	err := r.db.QueryRow(query, appointmentID).Scan(
		&apt.ID, &apt.PatientName, &apt.PatientID, &apt.AppointmentDate, &apt.TimeSlot,
		&apt.Reason, &apt.Status, &apt.Notes, &apt.PatientPhone, &apt.PatientEmail,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}

	return &apt, nil
}

// GetDashboardStats retrieves dashboard statistics for a doctor
func (r *DoctorRepository) GetDashboardStats(doctorID int) (*models.DoctorDashboardStats, error) {
	stats := &models.DoctorDashboardStats{}

	// Total patients
	err := r.db.QueryRow(
		`SELECT COUNT(DISTINCT patient_id) FROM doctor_patient_assignment WHERE doctor_id = $1 AND is_active = true`,
		doctorID,
	).Scan(&stats.TotalPatients)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Total appointments
	err = r.db.QueryRow(
		`SELECT COUNT(*) FROM appointments WHERE doctor_id = $1`,
		doctorID,
	).Scan(&stats.TotalAppointments)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Completed appointments
	err = r.db.QueryRow(
		`SELECT COUNT(*) FROM appointments WHERE doctor_id = $1 AND status = 'completed'`,
		doctorID,
	).Scan(&stats.CompletedAppointments)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Upcoming appointments (scheduled for future dates)
	err = r.db.QueryRow(
		`SELECT COUNT(*) FROM appointments WHERE doctor_id = $1 AND status = 'scheduled' AND appointment_date >= CURRENT_DATE`,
		doctorID,
	).Scan(&stats.UpcomingAppointments)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return stats, nil
}

// RegistrationNumberExists checks if a registration number already exists
func (r *DoctorRepository) RegistrationNumberExists(regNumber string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM doctors WHERE registration_number = $1`,
		regNumber,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
