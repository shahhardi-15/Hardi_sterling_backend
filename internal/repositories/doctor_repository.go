package repositories

import (
	"database/sql"
	"sterling-hms-backend/internal/models"
)

type DoctorRepository struct {
	db *sql.DB
}

func NewDoctorRepository(db *sql.DB) *DoctorRepository {
	return &DoctorRepository{db: db}
}

// FindByEmail finds a doctor by email
func (r *DoctorRepository) FindByEmail(email string) (*models.DoctorUser, error) {
	doctor := &models.DoctorUser{}

	err := r.db.QueryRow(
		"SELECT id, email, name, password_hash, specialization, phone, role, is_active, created_at, updated_at FROM doctor_users WHERE email = $1",
		email,
	).Scan(
		&doctor.ID,
		&doctor.Email,
		&doctor.Name,
		&doctor.PasswordHash,
		&doctor.Specialization,
		&doctor.Phone,
		&doctor.Role,
		&doctor.IsActive,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return doctor, nil
}

// FindByID finds a doctor by ID
func (r *DoctorRepository) FindByID(id int) (*models.DoctorUser, error) {
	doctor := &models.DoctorUser{}

	err := r.db.QueryRow(
		"SELECT id, email, name, password_hash, specialization, phone, role, is_active, created_at, updated_at FROM doctor_users WHERE id = $1",
		id,
	).Scan(
		&doctor.ID,
		&doctor.Email,
		&doctor.Name,
		&doctor.PasswordHash,
		&doctor.Specialization,
		&doctor.Phone,
		&doctor.Role,
		&doctor.IsActive,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return doctor, nil
}

// GetAssignedPatients gets all patients assigned to a doctor
func (r *DoctorRepository) GetAssignedPatients(doctorID int) ([]models.PatientInfo, error) {
	patients := []models.PatientInfo{}

	rows, err := r.db.Query(
		`SELECT 
			u.id, 
			CONCAT(u.first_name, ' ', u.last_name) as name, 
			COALESCE(pr.phone, '') as phone,
			u.email,
			COALESCE(pr.medical_conditions, '') as medical_history,
			COALESCE(pr.address, '') as address,
			COALESCE(pr.blood_type, '') as blood_type
		FROM doctor_patient_assignment dpa
		JOIN users u ON dpa.patient_id = u.id
		LEFT JOIN patient_records pr ON pr.user_id = u.id
		WHERE dpa.doctor_id = $1 AND dpa.is_active = true
		ORDER BY u.first_name, u.last_name`,
		doctorID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		patient := models.PatientInfo{}
		err := rows.Scan(
			&patient.ID,
			&patient.Name,
			&patient.Phone,
			&patient.Email,
			&patient.MedicalHistory,
			&patient.Address,
			&patient.BloodType,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, rows.Err()
}

// GetAppointments gets all appointments for a doctor
func (r *DoctorRepository) GetAppointments(doctorID int) ([]models.AppointmentInfoForDoctor, error) {
	appointments := []models.AppointmentInfoForDoctor{}

	rows, err := r.db.Query(
		`SELECT 
			a.id,
			CONCAT(u.first_name, ' ', u.last_name) as patient_name,
			a.patient_id,
			a.appointment_date,
			a.time_slot,
			a.reason,
			a.status,
			a.notes,
			COALESCE(a.patient_phone, '') as patient_phone,
			COALESCE(a.patient_email, '') as patient_email
		FROM appointments a
		JOIN users u ON a.patient_id = u.id
		WHERE a.doctor_id = $1
		ORDER BY a.appointment_date DESC, a.time_slot DESC`,
		doctorID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		appointment := models.AppointmentInfoForDoctor{}
		err := rows.Scan(
			&appointment.ID,
			&appointment.PatientName,
			&appointment.PatientID,
			&appointment.AppointmentDate,
			&appointment.TimeSlot,
			&appointment.Reason,
			&appointment.Status,
			&appointment.Notes,
			&appointment.PatientPhone,
			&appointment.PatientEmail,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	return appointments, rows.Err()
}

// GetAppointmentByID gets a specific appointment
func (r *DoctorRepository) GetAppointmentByID(appointmentID int) (*models.AppointmentInfoForDoctor, error) {
	appointment := &models.AppointmentInfoForDoctor{}

	err := r.db.QueryRow(
		`SELECT 
			a.id,
			CONCAT(u.first_name, ' ', u.last_name) as patient_name,
			a.patient_id,
			a.appointment_date,
			a.time_slot,
			a.reason,
			a.status,
			a.notes,
			COALESCE(a.patient_phone, '') as patient_phone,
			COALESCE(a.patient_email, '') as patient_email
		FROM appointments a
		JOIN users u ON a.patient_id = u.id
		WHERE a.id = $1`,
		appointmentID,
	).Scan(
		&appointment.ID,
		&appointment.PatientName,
		&appointment.PatientID,
		&appointment.AppointmentDate,
		&appointment.TimeSlot,
		&appointment.Reason,
		&appointment.Status,
		&appointment.Notes,
		&appointment.PatientPhone,
		&appointment.PatientEmail,
	)

	if err != nil {
		return nil, err
	}

	return appointment, nil
}

// CheckAppointmentOwnership verifies that an appointment belongs to a doctor
func (r *DoctorRepository) CheckAppointmentOwnership(appointmentID int, doctorID int) (bool, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM appointments WHERE id = $1 AND doctor_id = $2",
		appointmentID,
		doctorID,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateAppointmentStatus updates the status of an appointment
func (r *DoctorRepository) UpdateAppointmentStatus(appointmentID int, status string, notes string) error {
	_, err := r.db.Exec(
		"UPDATE appointments SET status = $1, notes = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
		status,
		notes,
		appointmentID,
	)
	return err
}

// GetDashboardStats gets statistics for doctor dashboard
func (r *DoctorRepository) GetDashboardStats(doctorID int) (*models.DoctorDashboardStats, error) {
	stats := &models.DoctorDashboardStats{}

	// Get total patients assigned to this doctor
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM doctor_patient_assignment WHERE doctor_id = $1 AND is_active = true",
		doctorID,
	).Scan(&stats.TotalPatients)
	if err != nil {
		return nil, err
	}

	// Get total appointments
	err = r.db.QueryRow(
		"SELECT COUNT(*) FROM appointments WHERE doctor_id = $1",
		doctorID,
	).Scan(&stats.TotalAppointments)
	if err != nil {
		return nil, err
	}

	// Get completed appointments
	err = r.db.QueryRow(
		"SELECT COUNT(*) FROM appointments WHERE doctor_id = $1 AND status = 'completed'",
		doctorID,
	).Scan(&stats.CompletedAppointments)
	if err != nil {
		return nil, err
	}

	// Get upcoming appointments
	err = r.db.QueryRow(
		"SELECT COUNT(*) FROM appointments WHERE doctor_id = $1 AND status = 'upcoming' OR (status = 'scheduled' AND appointment_date >= CURRENT_DATE)",
		doctorID,
	).Scan(&stats.UpcomingAppointments)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// AssignPatientToDoctor assigns a patient to a doctor
func (r *DoctorRepository) AssignPatientToDoctor(doctorID int, patientID int) error {
	_, err := r.db.Exec(
		`INSERT INTO doctor_patient_assignment (doctor_id, patient_id, is_active) 
		 VALUES ($1, $2, true)
		 ON CONFLICT (doctor_id, patient_id) DO UPDATE SET is_active = true`,
		doctorID,
		patientID,
	)
	return err
}
