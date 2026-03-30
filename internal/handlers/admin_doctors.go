package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sterling-hms-backend/internal/config"
	"sterling-hms-backend/internal/repositories"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminDoctorHandler struct {
	doctorRepo *repositories.DoctorRepository
	userRepo   *repositories.UserRepository
	cfg        *config.Config
}

func NewAdminDoctorHandler(doctorRepo *repositories.DoctorRepository, userRepo *repositories.UserRepository, cfg *config.Config) *AdminDoctorHandler {
	return &AdminDoctorHandler{
		doctorRepo: doctorRepo,
		userRepo:   userRepo,
		cfg:        cfg,
	}
}

// ListDoctors retrieves all doctors with pagination and search
func (h *AdminDoctorHandler) ListDoctors(c *gin.Context) {
	page := 1
	limit := 10
	search := ""

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	search = c.Query("search")

	offset := (page - 1) * limit

	// Build query with direct database access
	query := `
		SELECT d.id, d.user_id, u.full_name, u.email, u.phone, 
		       d.specialization, d.qualification, d.registration_number,
		       d.experience_years, d.consultation_fee, d.department_id,
		       d.available_days, d.start_time, d.end_time, d.slot_duration_minutes,
		       u.is_active, u.created_at, dept.name as department_name
		FROM doctors d
		LEFT JOIN users u ON d.user_id = u.id
		LEFT JOIN departments dept ON d.department_id = dept.id
	`

	args := []interface{}{}

	// Only return doctors with user_id set (linked doctors)
	query += ` WHERE d.user_id IS NOT NULL`

	// Add search condition if provided
	if search != "" {
		query += ` AND (u.full_name ILIKE $1 OR d.specialization ILIKE $1 OR d.registration_number ILIKE $1)`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY d.created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	log.Printf("[ADMIN_DOCTORS] Query: %s", query)
	log.Printf("[ADMIN_DOCTORS] Args: %v", args)

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying doctors: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch doctors",
			"success": false,
		})
		return
	}
	defer rows.Close()

	type DoctorList struct {
		ID                  string    `json:"id"`
		UserID              int       `json:"user_id"`
		FullName            string    `json:"full_name"`
		Email               string    `json:"email"`
		Phone               string    `json:"phone"`
		Specialization      string    `json:"specialization"`
		Qualification       string    `json:"qualification"`
		RegistrationNumber  string    `json:"registration_number"`
		ExperienceYears     int       `json:"experience_years"`
		ConsultationFee     float64   `json:"consultation_fee"`
		DepartmentID        string    `json:"department_id"`
		DepartmentName      string    `json:"department_name"`
		AvailableDays       string    `json:"available_days"`
		StartTime           string    `json:"start_time"`
		EndTime             string    `json:"end_time"`
		SlotDurationMinutes int       `json:"slot_duration_minutes"`
		IsActive            bool      `json:"is_active"`
		CreatedAt           time.Time `json:"created_at"`
	}

	enrichedDoctors := make([]DoctorList, 0)
	for rows.Next() {
		doctor := DoctorList{}
		var userID sql.NullInt64
		var fullName sql.NullString
		var email sql.NullString
		var phone sql.NullString
		var specialization sql.NullString
		var qualification sql.NullString
		var registrationNumber sql.NullString
		var experienceYears sql.NullInt64
		var consultationFee sql.NullFloat64
		var deptID sql.NullString
		var availableDays sql.NullString
		var startTime sql.NullString
		var endTime sql.NullString
		var slotDuration sql.NullInt64
		var isActive sql.NullBool
		var createdAt sql.NullTime
		var deptName sql.NullString

		err := rows.Scan(
			&doctor.ID, &userID, &fullName, &email, &phone,
			&specialization, &qualification, &registrationNumber,
			&experienceYears, &consultationFee, &deptID,
			&availableDays, &startTime, &endTime, &slotDuration,
			&isActive, &createdAt, &deptName,
		)
		if err != nil {
			log.Printf("Error scanning doctor row: %v", err)
			continue
		}

		// Convert nullable types to their non-null equivalents
		if userID.Valid {
			doctor.UserID = int(userID.Int64)
		}
		if fullName.Valid {
			doctor.FullName = fullName.String
		}
		if email.Valid {
			doctor.Email = email.String
		}
		if phone.Valid {
			doctor.Phone = phone.String
		}
		if specialization.Valid {
			doctor.Specialization = specialization.String
		}
		if qualification.Valid {
			doctor.Qualification = qualification.String
		}
		if registrationNumber.Valid {
			doctor.RegistrationNumber = registrationNumber.String
		}
		if experienceYears.Valid {
			doctor.ExperienceYears = int(experienceYears.Int64)
		}
		if consultationFee.Valid {
			doctor.ConsultationFee = consultationFee.Float64
		}
		if deptID.Valid {
			doctor.DepartmentID = deptID.String
		}
		if deptName.Valid {
			doctor.DepartmentName = deptName.String
		}
		if availableDays.Valid {
			doctor.AvailableDays = availableDays.String
		}
		if startTime.Valid {
			doctor.StartTime = startTime.String
		}
		if endTime.Valid {
			doctor.EndTime = endTime.String
		}
		if slotDuration.Valid {
			doctor.SlotDurationMinutes = int(slotDuration.Int64)
		}
		if isActive.Valid {
			doctor.IsActive = isActive.Bool
		}
		if createdAt.Valid {
			doctor.CreatedAt = createdAt.Time
		}

		enrichedDoctors = append(enrichedDoctors, doctor)
	}

	// Get total count (only count doctors with user_id set)
	countQuery := `SELECT COUNT(*) FROM doctors d LEFT JOIN users u ON d.user_id = u.id WHERE d.user_id IS NOT NULL`
	if search != "" {
		countQuery += ` AND (u.full_name ILIKE $1 OR d.specialization ILIKE $1 OR d.registration_number ILIKE $1)`
	}

	var total int
	var countArgs []interface{}
	if search != "" {
		countArgs = append(countArgs, "%"+search+"%")
	}

	err = config.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		log.Printf("Error counting doctors: %v", err)
		total = 0
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"message":     "Doctors retrieved successfully",
		"success":     true,
		"data":        enrichedDoctors,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// CreateDoctor creates a new doctor
func (h *AdminDoctorHandler) CreateDoctor(c *gin.Context) {
	var req struct {
		FullName           string   `json:"full_name" binding:"required"`
		Email              string   `json:"email" binding:"required,email"`
		Password           string   `json:"password" binding:"required,min=8"`
		Phone              string   `json:"phone" binding:"required"`
		Specialization     string   `json:"specialization" binding:"required"`
		Qualification      string   `json:"qualification" binding:"required"`
		RegistrationNumber string   `json:"registration_number" binding:"required"`
		ExperienceYears    int      `json:"experience_years" binding:"min=0"`
		ConsultationFee    float64  `json:"consultation_fee" binding:"min=0"`
		DepartmentID       string   `json:"department_id" binding:"required"`
		AvailableDays      []string `json:"available_days" binding:"required,min=1"`
		StartTime          string   `json:"start_time" binding:"required"`
		EndTime            string   `json:"end_time" binding:"required"`
		SlotDurationMin    int      `json:"slot_duration_minutes" binding:"required,oneof=10 15 20 30"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"success": false,
		})
		return
	}

	// Validate phone is exactly 10 digits
	if len(strings.TrimSpace(req.Phone)) != 10 || !isNumeric(req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Phone must be exactly 10 digits",
			"success": false,
		})
		return
	}

	// Check if email already exists
	exists, err := h.userRepo.EmailExists(req.Email)
	if err != nil {
		log.Printf("Error checking email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"message": "This email is already registered",
			"success": false,
		})
		return
	}

	// Check if registration number already exists
	regExists, err := h.doctorRepo.RegistrationNumberExists(req.RegistrationNumber)
	if err != nil {
		log.Printf("Error checking registration number: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	if regExists {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Registration number already exists",
			"success": false,
		})
		return
	}

	// Validate end time is after start time
	if req.EndTime <= req.StartTime {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "End Time must be after Start Time",
			"success": false,
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Start transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Insert into users table
	userID := 0
	err = tx.QueryRow(
		`INSERT INTO users (full_name, email, password, phone, role, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		 RETURNING id`,
		req.FullName,
		req.Email,
		string(hashedPassword),
		req.Phone,
		"doctor",
		true,
	).Scan(&userID)

	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Convert available days to comma-separated string
	availableDaysStr := strings.Join(req.AvailableDays, ",")

	// Insert into doctors table
	doctorID := ""
	err = tx.QueryRow(
		`INSERT INTO doctors (user_id, specialization, qualification, registration_number, 
		 experience_years, consultation_fee, department_id, available_days, 
		 start_time, end_time, slot_duration_minutes, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		 RETURNING id`,
		userID,
		req.Specialization,
		req.Qualification,
		req.RegistrationNumber,
		req.ExperienceYears,
		req.ConsultationFee,
		req.DepartmentID,
		availableDaysStr,
		req.StartTime,
		req.EndTime,
		req.SlotDurationMin,
	).Scan(&doctorID)

	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting doctor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Doctor created successfully",
		"success": true,
		"data": gin.H{
			"id":                    doctorID,
			"user_id":               userID,
			"full_name":             req.FullName,
			"email":                 req.Email,
			"phone":                 req.Phone,
			"specialization":        req.Specialization,
			"qualification":         req.Qualification,
			"registration_number":   req.RegistrationNumber,
			"experience_years":      req.ExperienceYears,
			"consultation_fee":      req.ConsultationFee,
			"department_id":         req.DepartmentID,
			"available_days":        req.AvailableDays,
			"start_time":            req.StartTime,
			"end_time":              req.EndTime,
			"slot_duration_minutes": req.SlotDurationMin,
			"is_active":             true,
		},
	})
}

// GetDoctor retrieves a single doctor by ID
func (h *AdminDoctorHandler) GetDoctor(c *gin.Context) {
	doctorID := c.Param("id")

	query := `
		SELECT d.id, d.user_id, u.full_name, u.email, u.phone,
		       d.specialization, d.qualification, d.registration_number,
		       d.experience_years, d.consultation_fee, d.department_id,
		       d.available_days, d.start_time, d.end_time, d.slot_duration_minutes,
		       u.is_active, u.created_at, dept.name as department_name
		FROM doctors d
		LEFT JOIN users u ON d.user_id = u.id
		LEFT JOIN departments dept ON d.department_id = dept.id
		WHERE d.id = $1
	`

	type DoctorDetail struct {
		ID                  string    `json:"id"`
		UserID              int       `json:"user_id"`
		FullName            string    `json:"full_name"`
		Email               string    `json:"email"`
		Phone               string    `json:"phone"`
		Specialization      string    `json:"specialization"`
		Qualification       string    `json:"qualification"`
		RegistrationNumber  string    `json:"registration_number"`
		ExperienceYears     int       `json:"experience_years"`
		ConsultationFee     float64   `json:"consultation_fee"`
		DepartmentID        string    `json:"department_id"`
		DepartmentName      string    `json:"department_name"`
		AvailableDays       string    `json:"available_days"`
		StartTime           string    `json:"start_time"`
		EndTime             string    `json:"end_time"`
		SlotDurationMinutes int       `json:"slot_duration_minutes"`
		IsActive            bool      `json:"is_active"`
		CreatedAt           time.Time `json:"created_at"`
	}

	doctor := DoctorDetail{}
	var deptID sql.NullString
	var deptName sql.NullString

	err := config.DB.QueryRow(query, doctorID).Scan(
		&doctor.ID, &doctor.UserID, &doctor.FullName, &doctor.Email, &doctor.Phone,
		&doctor.Specialization, &doctor.Qualification, &doctor.RegistrationNumber,
		&doctor.ExperienceYears, &doctor.ConsultationFee, &deptID,
		&doctor.AvailableDays, &doctor.StartTime, &doctor.EndTime, &doctor.SlotDurationMinutes,
		&doctor.IsActive, &doctor.CreatedAt, &deptName,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Doctor not found",
			"success": false,
		})
		return
	}

	if err != nil {
		log.Printf("Error querying doctor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch doctor",
			"success": false,
		})
		return
	}

	if deptID.Valid {
		doctor.DepartmentID = deptID.String
	}
	if deptName.Valid {
		doctor.DepartmentName = deptName.String
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor retrieved successfully",
		"success": true,
		"data":    doctor,
	})
}

// UpdateDoctor updates doctor information
func (h *AdminDoctorHandler) UpdateDoctor(c *gin.Context) {
	doctorID := c.Param("id")

	var req struct {
		FullName           string   `json:"full_name" binding:"required"`
		Phone              string   `json:"phone" binding:"required"`
		Specialization     string   `json:"specialization" binding:"required"`
		Qualification      string   `json:"qualification" binding:"required"`
		RegistrationNumber string   `json:"registration_number" binding:"required"`
		ExperienceYears    int      `json:"experience_years" binding:"min=0"`
		ConsultationFee    float64  `json:"consultation_fee" binding:"min=0"`
		DepartmentID       string   `json:"department_id" binding:"required"`
		AvailableDays      []string `json:"available_days" binding:"required,min=1"`
		StartTime          string   `json:"start_time" binding:"required"`
		EndTime            string   `json:"end_time" binding:"required"`
		SlotDurationMin    int      `json:"slot_duration_minutes" binding:"required,oneof=10 15 20 30"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"success": false,
		})
		return
	}

	// Validate phone is exactly 10 digits
	if len(strings.TrimSpace(req.Phone)) != 10 || !isNumeric(req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Phone must be exactly 10 digits",
			"success": false,
		})
		return
	}

	// Validate end time is after start time
	if req.EndTime <= req.StartTime {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "End Time must be after Start Time",
			"success": false,
		})
		return
	}

	// Get doctor to find user_id
	var userID int
	err := config.DB.QueryRow(`SELECT user_id FROM doctors WHERE id = $1`, doctorID).Scan(&userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Doctor not found",
			"success": false,
		})
		return
	}
	if err != nil {
		log.Printf("Error finding doctor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Start transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Update users table (only full_name and phone)
	_, err = tx.Exec(
		`UPDATE users SET full_name = $1, phone = $2, updated_at = NOW() WHERE id = $3`,
		req.FullName,
		req.Phone,
		userID,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Convert available days
	availableDaysStr := strings.Join(req.AvailableDays, ",")

	// Update doctors table
	_, err = tx.Exec(
		`UPDATE doctors SET 
		 specialization = $1, qualification = $2, registration_number = $3,
		 experience_years = $4, consultation_fee = $5, department_id = $6,
		 available_days = $7, start_time = $8, end_time = $9, 
		 slot_duration_minutes = $10, updated_at = NOW()
		 WHERE id = $11`,
		req.Specialization,
		req.Qualification,
		req.RegistrationNumber,
		req.ExperienceYears,
		req.ConsultationFee,
		req.DepartmentID,
		availableDaysStr,
		req.StartTime,
		req.EndTime,
		req.SlotDurationMin,
		doctorID,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating doctor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor updated successfully",
		"success": true,
	})
}

// UpdateDoctorStatus updates doctor's active status
func (h *AdminDoctorHandler) UpdateDoctorStatus(c *gin.Context) {
	doctorID := c.Param("id")

	var req struct {
		IsActive bool `json:"is_active" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"success": false,
		})
		return
	}

	// Get doctor to find user_id
	var userID int
	err := config.DB.QueryRow(`SELECT user_id FROM doctors WHERE id = $1`, doctorID).Scan(&userID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Doctor not found",
			"success": false,
		})
		return
	}
	if err != nil {
		log.Printf("Error finding doctor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Update is_active in users table
	_, err = config.DB.Exec(
		`UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2`,
		req.IsActive,
		userID,
	)
	if err != nil {
		log.Printf("Error updating doctor status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	status := "activated"
	if !req.IsActive {
		status = "deactivated"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Doctor %s successfully", status),
		"success": true,
	})
}

// Helper function to check if string is numeric
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
