package repositories

import (
	"database/sql"
	"sterling-hms-backend/internal/models"
)

type AppointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

// GetPatientProfile retrieves patient profile information
func (r *AppointmentRepository) GetPatientProfile(userID int) (*models.PatientProfile, error) {
	query := `
		SELECT id, first_name, last_name, email
		FROM users
		WHERE id = $1
	`

	var profile models.PatientProfile
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
	)

	if err != nil {
		return nil, err
	}

	// Phone, Address, and MedicalHistory are not stored in users table
	// They would need to be stored in a separate patient_records table
	// For now, these fields remain empty strings

	return &profile, nil
}

// GetAppointmentHistory retrieves all appointments for a patient
func (r *AppointmentRepository) GetAppointmentHistory(patientID int, limit int, offset int) ([]models.Appointment, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM appointments WHERE patient_id = $1`
	var total int
	err := r.db.QueryRow(countQuery, patientID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT 
			a.id, a.patient_id, a.doctor_id, a.appointment_date, a.time_slot,
			a.reason, a.status, a.notes, a.created_at, a.updated_at,
			d.id, d.name, d.specialization, d.email, d.phone, d.experience_years, d.qualification, d.address, d.is_available, d.created_at, d.updated_at
		FROM appointments a
		JOIN doctors d ON a.doctor_id = d.id
		WHERE a.patient_id = $1
		ORDER BY a.appointment_date DESC, a.time_slot DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, patientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var apt models.Appointment
		var doc models.Doctor

		err := rows.Scan(
			&apt.ID, &apt.PatientID, &apt.DoctorID, &apt.AppointmentDate, &apt.TimeSlot,
			&apt.Reason, &apt.Status, &apt.Notes, &apt.CreatedAt, &apt.UpdatedAt,
			&doc.ID, &doc.Name, &doc.Specialization, &doc.Email, &doc.Phone, &doc.ExperienceYears, &doc.Qualification, &doc.Address, &doc.IsAvailable, &doc.CreatedAt, &doc.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		apt.Doctor = &doc
		appointments = append(appointments, apt)
	}

	return appointments, total, rows.Err()
}

// GetAvailableSlots retrieves available appointment slots for a specific doctor
func (r *AppointmentRepository) GetAvailableSlots(doctorID int, startDate string, endDate string) ([]models.AppointmentSlot, error) {
	query := `
		SELECT id, doctor_id, slot_date, time_slot, is_available, created_at
		FROM appointment_slots
		WHERE doctor_id = $1 
			AND slot_date >= $2 
			AND slot_date <= $3
			AND is_available = true
		ORDER BY slot_date ASC, time_slot ASC
	`

	rows, err := r.db.Query(query, doctorID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []models.AppointmentSlot
	for rows.Next() {
		var slot models.AppointmentSlot
		err := rows.Scan(
			&slot.ID, &slot.DoctorID, &slot.SlotDate, &slot.TimeSlot, &slot.IsAvailable, &slot.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	return slots, rows.Err()
}

// GetDoctors retrieves all available doctors
func (r *AppointmentRepository) GetDoctors() ([]models.Doctor, error) {
	query := `
		SELECT id, name, specialization, email, phone, experience_years, qualification, address, is_available, created_at, updated_at
		FROM doctors
		WHERE is_available = true
		ORDER BY specialization, name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var doc models.Doctor
		err := rows.Scan(
			&doc.ID, &doc.Name, &doc.Specialization, &doc.Email, &doc.Phone,
			&doc.ExperienceYears, &doc.Qualification, &doc.Address, &doc.IsAvailable, &doc.CreatedAt, &doc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doc)
	}

	return doctors, rows.Err()
}

// CreateAppointment creates a new appointment
func (r *AppointmentRepository) CreateAppointment(patientID int, doctorID int, appointmentDate string, timeSlot string, reason string, notes string) (*models.Appointment, error) {
	query := `
		INSERT INTO appointments (patient_id, doctor_id, appointment_date, time_slot, reason, notes, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'scheduled', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, patient_id, doctor_id, appointment_date, time_slot, reason, status, notes, created_at, updated_at
	`

	var apt models.Appointment
	err := r.db.QueryRow(query, patientID, doctorID, appointmentDate, timeSlot, reason, notes).Scan(
		&apt.ID, &apt.PatientID, &apt.DoctorID, &apt.AppointmentDate, &apt.TimeSlot,
		&apt.Reason, &apt.Status, &apt.Notes, &apt.CreatedAt, &apt.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Mark slot as unavailable
	_ = r.MarkSlotUnavailable(doctorID, appointmentDate, timeSlot)

	return &apt, nil
}

// CancelAppointment cancels an appointment (only if status is scheduled)
func (r *AppointmentRepository) CancelAppointment(appointmentID int, patientID int) error {
	query := `
		UPDATE appointments 
		SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND patient_id = $2 AND status = 'scheduled'
	`

	result, err := r.db.Exec(query, appointmentID, patientID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Get appointment details to free up the slot
	apt, err := r.GetAppointmentByID(appointmentID)
	if err == nil && apt != nil {
		_ = r.MarkSlotAvailable(apt.DoctorID, apt.AppointmentDate, apt.TimeSlot)
	}

	return nil
}

// GetAppointmentByID retrieves a single appointment
func (r *AppointmentRepository) GetAppointmentByID(appointmentID int) (*models.Appointment, error) {
	query := `
		SELECT 
			a.id, a.patient_id, a.doctor_id, a.appointment_date, a.time_slot,
			a.reason, a.status, a.notes, a.created_at, a.updated_at,
			d.id, d.name, d.specialization, d.email, d.phone, d.experience_years, d.qualification, d.address, d.is_available, d.created_at, d.updated_at
		FROM appointments a
		JOIN doctors d ON a.doctor_id = d.id
		WHERE a.id = $1
	`

	var apt models.Appointment
	var doc models.Doctor

	err := r.db.QueryRow(query, appointmentID).Scan(
		&apt.ID, &apt.PatientID, &apt.DoctorID, &apt.AppointmentDate, &apt.TimeSlot,
		&apt.Reason, &apt.Status, &apt.Notes, &apt.CreatedAt, &apt.UpdatedAt,
		&doc.ID, &doc.Name, &doc.Specialization, &doc.Email, &doc.Phone, &doc.ExperienceYears, &doc.Qualification, &doc.Address, &doc.IsAvailable, &doc.CreatedAt, &doc.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	apt.Doctor = &doc
	return &apt, nil
}

// CheckSlotAvailability checks if a specific slot is available
func (r *AppointmentRepository) CheckSlotAvailability(doctorID int, appointmentDate string, timeSlot string) (bool, error) {
	query := `
		SELECT is_available 
		FROM appointment_slots 
		WHERE doctor_id = $1 AND slot_date = $2 AND time_slot = $3
	`

	var isAvailable bool
	err := r.db.QueryRow(query, doctorID, appointmentDate, timeSlot).Scan(&isAvailable)

	if err != nil {
		if err == sql.ErrNoRows {
			// Slot doesn't exist - check if any slots exist for this doctor and date
			var count int
			checkQuery := `SELECT COUNT(*) FROM appointment_slots WHERE doctor_id = $1 AND slot_date = $2`
			r.db.QueryRow(checkQuery, doctorID, appointmentDate).Scan(&count)

			if count == 0 {
				// No slots at all for this date - this might be normal for past dates
				// But it could also mean no slots were created for this date
				return false, nil
			}
			// Slots exist for this date, but not this specific time
			return false, nil
		}
		return false, err
	}

	return isAvailable, nil
}

// MarkSlotUnavailable marks a slot as unavailable
func (r *AppointmentRepository) MarkSlotUnavailable(doctorID int, slotDate string, timeSlot string) error {
	query := `
		UPDATE appointment_slots 
		SET is_available = false
		WHERE doctor_id = $1 AND slot_date = $2 AND time_slot = $3
	`

	_, err := r.db.Exec(query, doctorID, slotDate, timeSlot)
	return err
}

// MarkSlotAvailable marks a slot as available
func (r *AppointmentRepository) MarkSlotAvailable(doctorID int, slotDate string, timeSlot string) error {
	query := `
		UPDATE appointment_slots 
		SET is_available = true
		WHERE doctor_id = $1 AND slot_date = $2 AND time_slot = $3
	`

	_, err := r.db.Exec(query, doctorID, slotDate, timeSlot)
	return err
}

// GetDoctorByID retrieves a single doctor with full details
func (r *AppointmentRepository) GetDoctorByID(doctorID int) (*models.Doctor, error) {
	query := `
		SELECT id, name, specialization, email, phone, experience_years, qualification, address, is_available, created_at, updated_at
		FROM doctors
		WHERE id = $1
	`

	var doc models.Doctor
	err := r.db.QueryRow(query, doctorID).Scan(
		&doc.ID, &doc.Name, &doc.Specialization, &doc.Email, &doc.Phone,
		&doc.ExperienceYears, &doc.Qualification, &doc.Address, &doc.IsAvailable, &doc.CreatedAt, &doc.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &doc, nil
}

// GetSpecializations retrieves all unique specializations
func (r *AppointmentRepository) GetSpecializations() ([]models.Specialization, error) {
	query := `
		SELECT DISTINCT specialization
		FROM doctors
		WHERE is_available = true
		ORDER BY specialization ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specializations []models.Specialization
	id := 1
	for rows.Next() {
		var spec models.Specialization
		err := rows.Scan(&spec.Name)
		if err != nil {
			return nil, err
		}
		spec.ID = id
		id++
		specializations = append(specializations, spec)
	}

	return specializations, rows.Err()
}

// GetDoctorsBySpecialization retrieves all doctors for a specific specialization
func (r *AppointmentRepository) GetDoctorsBySpecialization(specialization string) ([]models.Doctor, error) {
	query := `
		SELECT id, name, specialization, email, phone, experience_years, qualification, address, is_available, created_at, updated_at
		FROM doctors
		WHERE specialization = $1 AND is_available = true
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query, specialization)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var doc models.Doctor
		err := rows.Scan(
			&doc.ID, &doc.Name, &doc.Specialization, &doc.Email, &doc.Phone,
			&doc.ExperienceYears, &doc.Qualification, &doc.Address, &doc.IsAvailable, &doc.CreatedAt, &doc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doc)
	}

	return doctors, rows.Err()
}
