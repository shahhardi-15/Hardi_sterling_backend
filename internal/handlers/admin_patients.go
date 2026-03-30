package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sterling-hms-backend/internal/config"
	"sterling-hms-backend/internal/repositories"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminPatientHandler struct {
	patientRepo *repositories.PatientRepository
	userRepo    *repositories.UserRepository
	cfg         *config.Config
}

func NewAdminPatientHandler(patientRepo *repositories.PatientRepository, userRepo *repositories.UserRepository, cfg *config.Config) *AdminPatientHandler {
	return &AdminPatientHandler{
		patientRepo: patientRepo,
		userRepo:    userRepo,
		cfg:         cfg,
	}
}

// ListPatients retrieves all patients with pagination and search
func (h *AdminPatientHandler) ListPatients(c *gin.Context) {
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
		SELECT pr.id, pr.user_id, pr.first_name, pr.last_name, pr.email, pr.phone,
		       pr.date_of_birth, pr.gender, pr.blood_type, pr.address, pr.city, pr.state,
		       pr.postal_code, pr.country, pr.allergies, pr.medical_conditions, 
		       pr.current_medications, pr.emergency_contact_name, pr.emergency_contact_phone,
		       pr.created_at, pr.updated_at, COALESCE(u.uhid, ''), u.is_active,
		       COALESCE(pr.registration_date, pr.created_at)
		FROM patient_records pr
		LEFT JOIN users u ON pr.user_id = u.id
	`

	args := []interface{}{}

	// Add search condition if provided
	if search != "" {
		query += ` WHERE pr.first_name ILIKE $1 OR pr.last_name ILIKE $1 OR pr.phone ILIKE $1 OR u.uhid ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY pr.created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying patients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch patients",
			"success": false,
		})
		return
	}
	defer rows.Close()

	type PatientList struct {
		ID                    int        `json:"id"`
		UserID                int        `json:"user_id"`
		UHID                  string     `json:"uhid"`
		FullName              string     `json:"full_name"`
		Email                 string     `json:"email"`
		Phone                 string     `json:"phone"`
		Gender                string     `json:"gender"`
		DateOfBirth           *time.Time `json:"date_of_birth"`
		BloodType             string     `json:"blood_type"`
		Address               string     `json:"address"`
		City                  string     `json:"city"`
		State                 string     `json:"state"`
		PostalCode            string     `json:"postal_code"`
		Country               string     `json:"country"`
		Allergies             string     `json:"allergies"`
		MedicalConditions     string     `json:"medical_conditions"`
		CurrentMedications    string     `json:"current_medications"`
		EmergencyContactName  string     `json:"emergency_contact_name"`
		EmergencyContactPhone string     `json:"emergency_contact_phone"`
		RegistrationDate      *time.Time `json:"registration_date"`
		IsActive              bool       `json:"is_active"`
		CreatedAt             time.Time  `json:"created_at"`
		UpdatedAt             time.Time  `json:"updated_at"`
	}

	var enrichedPatients []PatientList = []PatientList{} // Initialize as empty slice instead of nil
	for rows.Next() {
		var firstName, lastName string
		var uhid string
		var isActive bool
		var regDate *time.Time
		patient := PatientList{}

		err := rows.Scan(
			&patient.ID, &patient.UserID, &firstName, &lastName, &patient.Email, &patient.Phone,
			&patient.DateOfBirth, &patient.Gender, &patient.BloodType, &patient.Address,
			&patient.City, &patient.State, &patient.PostalCode, &patient.Country,
			&patient.Allergies, &patient.MedicalConditions, &patient.CurrentMedications,
			&patient.EmergencyContactName, &patient.EmergencyContactPhone,
			&patient.CreatedAt, &patient.UpdatedAt, &uhid, &isActive, &regDate,
		)
		if err != nil {
			log.Printf("Error scanning patient row: %v", err)
			continue
		}

		patient.UHID = uhid
		patient.IsActive = isActive
		patient.RegistrationDate = regDate
		patient.FullName = firstName
		if lastName != "" {
			patient.FullName = firstName + " " + lastName
		}

		enrichedPatients = append(enrichedPatients, patient)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM patient_records pr LEFT JOIN users u ON pr.user_id = u.id`
	if search != "" {
		countQuery += ` WHERE pr.first_name ILIKE $1 OR pr.last_name ILIKE $1 OR pr.phone ILIKE $1 OR u.uhid ILIKE $1`
	}

	var total int
	var countArgs []interface{}
	if search != "" {
		countArgs = append(countArgs, "%"+search+"%")
	}

	err = config.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		log.Printf("Error counting patients: %v", err)
		total = 0
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"message":     "Patients retrieved successfully",
		"success":     true,
		"data":        enrichedPatients,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// CreatePatient creates a new patient
func (h *AdminPatientHandler) CreatePatient(c *gin.Context) {
	var req struct {
		FullName    string `json:"full_name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
		Gender      string `json:"gender"`
		DateOfBirth string `json:"date_of_birth"`
		Phone       string `json:"phone" binding:"required"`
		Address     string `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
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

	// Parse date of birth
	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		t, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			dateOfBirth = &t
		}
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

	// Parse full name
	firstName := req.FullName
	lastName := ""
	// Simple split on first space
	for i, r := range req.FullName {
		if r == ' ' {
			firstName = req.FullName[:i]
			lastName = req.FullName[i+1:]
			break
		}
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
	defer tx.Rollback()

	// Generate UHID
	uhid, err := h.patientRepo.GenerateUHID()
	if err != nil {
		log.Printf("Error generating UHID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Create user
	var userID int
	err = tx.QueryRow(
		`INSERT INTO users (first_name, last_name, email, password, uhid, created_at, updated_at, is_active)
		 VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), true)
		 RETURNING id`,
		firstName, lastName, req.Email, string(hashedPassword), uhid,
	).Scan(&userID)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Create patient record
	var patientID int
	err = tx.QueryRow(
		`INSERT INTO patient_records (
			user_id, first_name, last_name, email, phone, date_of_birth, gender,
			address, created_at, updated_at, registration_date
		)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), NOW())
		 RETURNING id`,
		userID, firstName, lastName, req.Email, req.Phone, dateOfBirth, req.Gender, req.Address,
	).Scan(&patientID)

	if err != nil {
		log.Printf("Error creating patient record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Patient registered successfully",
		"success": true,
		"data": gin.H{
			"id":    patientID,
			"uhid":  uhid,
			"email": req.Email,
		},
	})
}

// GetPatient retrieves a single patient by ID
func (h *AdminPatientHandler) GetPatient(c *gin.Context) {
	patientIDStr := c.Param("id")
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid patient ID",
			"success": false,
		})
		return
	}

	var patient struct {
		ID                    int        `json:"id"`
		UserID                int        `json:"user_id"`
		UHID                  string     `json:"uhid"`
		FullName              string     `json:"full_name"`
		Email                 string     `json:"email"`
		Phone                 string     `json:"phone"`
		Gender                string     `json:"gender"`
		DateOfBirth           *time.Time `json:"date_of_birth"`
		BloodType             string     `json:"blood_type"`
		Address               string     `json:"address"`
		City                  string     `json:"city"`
		State                 string     `json:"state"`
		PostalCode            string     `json:"postal_code"`
		Country               string     `json:"country"`
		Allergies             string     `json:"allergies"`
		MedicalConditions     string     `json:"medical_conditions"`
		CurrentMedications    string     `json:"current_medications"`
		EmergencyContactName  string     `json:"emergency_contact_name"`
		EmergencyContactPhone string     `json:"emergency_contact_phone"`
		RegistrationDate      *time.Time `json:"registration_date"`
		IsActive              bool       `json:"is_active"`
		CreatedAt             time.Time  `json:"created_at"`
		UpdatedAt             time.Time  `json:"updated_at"`
	}

	var firstName, lastName string
	err = config.DB.QueryRow(`
		SELECT pr.id, pr.user_id, pr.first_name, pr.last_name, pr.email, pr.phone,
		       pr.date_of_birth, pr.gender, pr.blood_type, pr.address, pr.city, pr.state,
		       pr.postal_code, pr.country, pr.allergies, pr.medical_conditions,
		       pr.current_medications, pr.emergency_contact_name, pr.emergency_contact_phone,
		       pr.created_at, pr.updated_at, COALESCE(u.uhid, ''), COALESCE(u.is_active, false),
		       COALESCE(pr.registration_date, pr.created_at)
		FROM patient_records pr
		LEFT JOIN users u ON pr.user_id = u.id
		WHERE pr.id = $1`,
		patientID,
	).Scan(
		&patient.ID, &patient.UserID, &firstName, &lastName, &patient.Email, &patient.Phone,
		&patient.DateOfBirth, &patient.Gender, &patient.BloodType, &patient.Address,
		&patient.City, &patient.State, &patient.PostalCode, &patient.Country,
		&patient.Allergies, &patient.MedicalConditions, &patient.CurrentMedications,
		&patient.EmergencyContactName, &patient.EmergencyContactPhone,
		&patient.CreatedAt, &patient.UpdatedAt, &patient.UHID, &patient.IsActive,
		&patient.RegistrationDate,
	)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Patient not found",
				"success": false,
			})
		} else {
			log.Printf("Error fetching patient: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Server error",
				"success": false,
			})
		}
		return
	}

	patient.FullName = firstName
	if lastName != "" {
		patient.FullName = firstName + " " + lastName
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient retrieved successfully",
		"success": true,
		"data":    patient,
	})
}

// UpdatePatient updates patient information
func (h *AdminPatientHandler) UpdatePatient(c *gin.Context) {
	patientIDStr := c.Param("id")
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid patient ID",
			"success": false,
		})
		return
	}

	var req struct {
		FullName    string `json:"full_name"`
		Gender      string `json:"gender"`
		DateOfBirth string `json:"date_of_birth"`
		Phone       string `json:"phone"`
		Address     string `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"success": false,
		})
		return
	}

	// Verify patient exists
	var userID int
	var firstName, lastName string
	err = config.DB.QueryRow(`
		SELECT pr.user_id, pr.first_name, pr.last_name, 
		       COALESCE(pr.date_of_birth::text, ''), 
		       COALESCE(pr.phone, ''), 
		       COALESCE(pr.gender, ''),
		       COALESCE(pr.address, '')
		FROM patient_records pr
		WHERE pr.id = $1`,
		patientID,
	).Scan(&userID, &firstName, &lastName, nil, nil, nil, nil)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Patient not found",
				"success": false,
			})
		} else {
			log.Printf("Error checking patient existence: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Server error",
				"success": false,
			})
		}
		return
	}

	// Parse date of birth if provided
	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		t, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			dateOfBirth = &t
		}
	}

	// Use provided values or keep existing
	updateFullName := req.FullName
	if updateFullName == "" {
		updateFullName = firstName
		if lastName != "" {
			updateFullName = firstName + " " + lastName
		}
	}

	updatePhone := req.Phone
	updateGender := req.Gender
	updateAddress := req.Address

	// Parse full name for first/last
	updateFirstName := updateFullName
	updateLastName := ""
	for i, r := range updateFullName {
		if r == ' ' {
			updateFirstName = updateFullName[:i]
			updateLastName = updateFullName[i+1:]
			break
		}
	}

	// Update patient
	err = h.patientRepo.UpdatePatient(patientID, updateFirstName, updateLastName, updatePhone, dateOfBirth, updateGender, updateAddress)
	if err != nil {
		log.Printf("Error updating patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update patient",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient updated successfully",
		"success": true,
	})
}

// UpdatePatientStatus updates the active status of a patient
func (h *AdminPatientHandler) UpdatePatientStatus(c *gin.Context) {
	patientIDStr := c.Param("id")
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid patient ID",
			"success": false,
		})
		return
	}

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

	// Verify patient exists and get user ID
	var userID int
	err = config.DB.QueryRow(`
		SELECT pr.user_id FROM patient_records pr WHERE pr.id = $1`,
		patientID,
	).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Patient not found",
				"success": false,
			})
		} else {
			log.Printf("Error checking patient: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Server error",
				"success": false,
			})
		}
		return
	}

	// Update user status
	err = h.patientRepo.UpdatePatientStatus(userID, req.IsActive)
	if err != nil {
		log.Printf("Error updating patient status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update patient status",
			"success": false,
		})
		return
	}

	message := "Patient activated successfully"
	if !req.IsActive {
		message = "Patient deactivated successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"success": true,
	})
}
