package handlers

import (
	"log"
	"net/http"
	"sterling-hms-backend/internal/config"
	"sterling-hms-backend/internal/models"
	"sterling-hms-backend/internal/repositories"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ReceptionistHandler struct {
	receptionistRepo *repositories.ReceptionistRepository
	userRepo         *repositories.UserRepository
	appointmentRepo  *repositories.AppointmentRepository
	cfg              *config.Config
}

func NewReceptionistHandler(receptionistRepo *repositories.ReceptionistRepository, userRepo *repositories.UserRepository,
	appointmentRepo *repositories.AppointmentRepository, cfg *config.Config) *ReceptionistHandler {
	return &ReceptionistHandler{
		receptionistRepo: receptionistRepo,
		userRepo:         userRepo,
		appointmentRepo:  appointmentRepo,
		cfg:              cfg,
	}
}

// ReceptionistLogin handles receptionist authentication
func (h *ReceptionistHandler) ReceptionistLogin(c *gin.Context) {
	var req models.ReceptionistLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid credentials",
			"success": false,
		})
		return
	}

	log.Printf("Receptionist login attempt for email: %s", req.Email)

	// Find receptionist by email
	receptionist, err := h.receptionistRepo.FindByEmail(req.Email)
	if err != nil {
		log.Printf("Receptionist not found for email %s: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
			"success": false,
		})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(receptionist.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Printf("Password verification failed for receptionist %s: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
			"success": false,
		})
		return
	}

	// Check if receptionist is active
	if !receptionist.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Account is inactive",
			"success": false,
		})
		return
	}

	// Generate JWT token
	token, err := h.generateReceptionistToken(receptionist)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Prepare response
	receptionistResponse := &models.ReceptionistUser{
		ID:         receptionist.ID,
		Email:      receptionist.Email,
		Name:       receptionist.Name,
		Phone:      receptionist.Phone,
		Department: receptionist.Department,
		Role:       receptionist.Role,
		CreatedAt:  receptionist.CreatedAt,
		UpdatedAt:  receptionist.UpdatedAt,
		IsActive:   receptionist.IsActive,
	}

	c.JSON(http.StatusOK, models.ReceptionistLoginResponse{
		Message:      "Login successful",
		Success:      true,
		Receptionist: receptionistResponse,
		Token:        token,
	})
}

// RegisterPatient handles patient registration by receptionist
func (h *ReceptionistHandler) RegisterPatient(c *gin.Context) {
	var req models.RegisterPatientRequest

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
		log.Printf("Error checking email existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Email already registered",
			"success": false,
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"success": false,
		})
		return
	}

	// Parse date of birth
	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid date of birth format (use YYYY-MM-DD)",
				"success": false,
			})
			return
		}
		dateOfBirth = &dob
	}

	// Create patient and patient record
	patientRecord, err := h.receptionistRepo.CreatePatient(
		req.FirstName, req.LastName, req.Email, string(hashedPassword), req.Phone,
		dateOfBirth, req.Gender, req.BloodType, req.Address, req.City, req.State,
		req.PostalCode, req.Country, req.Allergies, req.MedicalConditions,
		req.CurrentMedications, req.EmergencyContactName, req.EmergencyContactPhone,
	)
	if err != nil {
		log.Printf("Error creating patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error registering patient",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Patient registered successfully",
		"success": true,
		"patient": patientRecord,
	})
}

// BookAppointmentByReceptionist books appointment on behalf of patient
func (h *ReceptionistHandler) BookAppointmentByReceptionist(c *gin.Context) {
	var req struct {
		PatientID       int    `json:"patientId" binding:"required"`
		DoctorID        int    `json:"doctorId" binding:"required"`
		AppointmentDate string `json:"appointmentDate" binding:"required"`
		TimeSlot        string `json:"timeSlot" binding:"required"`
		Reason          string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"success": false,
		})
		return
	}

	// Get patient info
	user, err := h.userRepo.FindByID(req.PatientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Patient not found",
			"success": false,
		})
		return
	}

	// Get patient record for phone
	patientRecord, err := h.receptionistRepo.GetPatientRecord(req.PatientID)
	if err != nil {
		log.Printf("Error getting patient record: %v", err)
		// Continue without phone if not found
		patientRecord = &models.PatientRecord{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}
	}

	// Book appointment
	appointment, err := h.receptionistRepo.BookAppointmentByReceptionist(
		req.PatientID, req.DoctorID, req.AppointmentDate, req.TimeSlot, req.Reason,
		user.FirstName, user.LastName, user.Email, patientRecord.Phone,
	)
	if err != nil {
		log.Printf("Error booking appointment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error booking appointment",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Appointment booked successfully (pending approval)",
		"success":     true,
		"appointment": appointment,
	})
}

// GetPendingAppointments gets all pending appointments for approval
func (h *ReceptionistHandler) GetPendingAppointments(c *gin.Context) {
	page := 1
	limit := 10

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

	appointments, total, err := h.receptionistRepo.GetPendingAppointments(page, limit)
	if err != nil {
		log.Printf("Error getting pending appointments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching appointments",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Pending appointments retrieved",
		"success":      true,
		"appointments": appointments,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

// ApproveAppointment approves a pending appointment
func (h *ReceptionistHandler) ApproveAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	receptionistID := c.GetInt("receptionistId")

	id, err := strconv.Atoi(appointmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid appointment ID",
			"success": false,
		})
		return
	}

	// Check if appointment exists
	appointment, err := h.receptionistRepo.GetAppointmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Appointment not found",
			"success": false,
		})
		return
	}

	// Check if already approved or rejected
	if appointment.ApprovalStatus != "pending" {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Appointment is already " + appointment.ApprovalStatus,
			"success": false,
		})
		return
	}

	// Approve appointment
	err = h.receptionistRepo.ApproveAppointment(id, receptionistID)
	if err != nil {
		log.Printf("Error approving appointment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error approving appointment",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment approved successfully",
		"success": true,
	})
}

// RejectAppointment rejects a pending appointment
func (h *ReceptionistHandler) RejectAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	receptionistID := c.GetInt("receptionistId")

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"success": false,
		})
		return
	}

	id, err := strconv.Atoi(appointmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid appointment ID",
			"success": false,
		})
		return
	}

	// Check if appointment exists
	appointment, err := h.receptionistRepo.GetAppointmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Appointment not found",
			"success": false,
		})
		return
	}

	// Check if already approved or rejected
	if appointment.ApprovalStatus != "pending" {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Appointment is already " + appointment.ApprovalStatus,
			"success": false,
		})
		return
	}

	// Reject appointment
	err = h.receptionistRepo.RejectAppointment(id, receptionistID, req.Reason)
	if err != nil {
		log.Printf("Error rejecting appointment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error rejecting appointment",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment rejected",
		"success": true,
	})
}

// GetPatientRecords gets all registered patients
func (h *ReceptionistHandler) GetPatientRecords(c *gin.Context) {
	page := 1
	limit := 10

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

	patients, total, err := h.receptionistRepo.GetAllPatients(page, limit)
	if err != nil {
		log.Printf("Error getting patient records: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching patient records",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Patient records retrieved",
		"success":  true,
		"patients": patients,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// GetDashboardStats gets receptionist dashboard statistics
func (h *ReceptionistHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.receptionistRepo.GetReceptionistDashboardStats()
	if err != nil {
		log.Printf("Error getting dashboard stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching dashboard stats",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, models.ReceptionistDashboardResponse{
		Message: "Dashboard stats retrieved",
		Success: true,
		Stats:   *stats,
	})
}

// generateReceptionistToken generates a JWT token for receptionist
func (h *ReceptionistHandler) generateReceptionistToken(receptionist *models.ReceptionistUser) (string, error) {
	claims := models.ReceptionistClaims{
		ID:    receptionist.ID,
		Email: receptionist.Email,
		Role:  receptionist.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sterling-hms",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

// ReceptionistLogout handles receptionist logout
func (h *ReceptionistHandler) ReceptionistLogout(c *gin.Context) {
	// Since we're using stateless JWT, logout is just a client-side action
	// but we can still return a success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
		"success": true,
	})
}
