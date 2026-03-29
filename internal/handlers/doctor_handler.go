package handlers

import (
	"log"
	"net/http"
	"sterling-hms-backend/internal/config"
	"sterling-hms-backend/internal/models"
	"sterling-hms-backend/internal/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type DoctorHandler struct {
	doctorRepo *repositories.DoctorRepository
	cfg        *config.Config
}

func NewDoctorHandler(doctorRepo *repositories.DoctorRepository, cfg *config.Config) *DoctorHandler {
	return &DoctorHandler{
		doctorRepo: doctorRepo,
		cfg:        cfg,
	}
}

// GetAssignedPatients returns all patients assigned to a doctor
func (h *DoctorHandler) GetAssignedPatients(c *gin.Context) {
	doctorID, exists := c.Get("doctorID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Doctor ID not found in token"})
		return
	}

	patients, err := h.doctorRepo.GetAssignedPatients(doctorID.(int))
	if err != nil {
		log.Printf("Error fetching assigned patients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch patients"})
		return
	}

	c.JSON(http.StatusOK, models.DoctorPatientsResponse{
		Message:  "Assigned patients retrieved successfully",
		Success:  true,
		Patients: patients,
		Total:    len(patients),
	})
}

// GetAppointments returns all appointments for a doctor
func (h *DoctorHandler) GetAppointments(c *gin.Context) {
	doctorID, exists := c.Get("doctorID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Doctor ID not found in token"})
		return
	}

	appointments, err := h.doctorRepo.GetAppointments(doctorID.(int))
	if err != nil {
		log.Printf("Error fetching appointments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch appointments"})
		return
	}

	c.JSON(http.StatusOK, models.DoctorAppointmentsResponse{
		Message:      "Appointments retrieved successfully",
		Success:      true,
		Appointments: appointments,
		Total:        len(appointments),
	})
}

// UpdateAppointmentStatus updates the status of an appointment
func (h *DoctorHandler) UpdateAppointmentStatus(c *gin.Context) {
	doctorID, exists := c.Get("doctorID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Doctor ID not found in token"})
		return
	}

	appointmentIDStr := c.Param("appointmentId")
	appointmentID, err := strconv.Atoi(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid appointment ID"})
		return
	}

	// Verify appointment belongs to this doctor
	owns, err := h.doctorRepo.CheckAppointmentOwnership(appointmentID, doctorID.(int))
	if err != nil {
		log.Printf("Error checking appointment ownership: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have permission to update this appointment"})
		return
	}

	var req models.UpdateAppointmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"upcoming":  true,
		"completed": true,
		"cancelled": true,
		"no-show":   true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid appointment status"})
		return
	}

	// Update appointment
	err = h.doctorRepo.UpdateAppointmentStatus(appointmentID, req.Status, req.Notes)
	if err != nil {
		log.Printf("Error updating appointment status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update appointment status"})
		return
	}

	// Fetch updated appointment
	appointment, err := h.doctorRepo.GetAppointmentByID(appointmentID)
	if err != nil {
		log.Printf("Error fetching updated appointment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Appointment updated but failed to retrieve"})
		return
	}

	c.JSON(http.StatusOK, models.UpdateAppointmentStatusResponse{
		Message:     "Appointment status updated successfully",
		Success:     true,
		Appointment: appointment,
	})
}

// GetDashboardStats returns dashboard statistics for a doctor
func (h *DoctorHandler) GetDashboardStats(c *gin.Context) {
	doctorID, exists := c.Get("doctorID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Doctor ID not found in token"})
		return
	}

	stats, err := h.doctorRepo.GetDashboardStats(doctorID.(int))
	if err != nil {
		log.Printf("Error fetching dashboard stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch statistics"})
		return
	}

	c.JSON(http.StatusOK, models.DoctorDashboardResponse{
		Message: "Dashboard statistics retrieved successfully",
		Success: true,
		Stats:   *stats,
	})
}

// GetProfile returns the profile of a doctor
func (h *DoctorHandler) GetProfile(c *gin.Context) {
	doctorID, exists := c.Get("doctorID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Doctor ID not found in token"})
		return
	}

	doctor, err := h.doctorRepo.FindByID(doctorID.(int))
	if err != nil {
		log.Printf("Error fetching doctor profile: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Doctor profile not found"})
		return
	}

	// Don't send password hash
	doctor.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor profile retrieved successfully",
		"success": true,
		"doctor":  doctor,
	})
}

// VerifyDoctorCredentials verifies doctor email and password (used in auth handler)
func (h *DoctorHandler) VerifyDoctorCredentials(email string, password string) (*models.DoctorUser, error) {
	doctor, err := h.doctorRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(doctor.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return doctor, nil
}
