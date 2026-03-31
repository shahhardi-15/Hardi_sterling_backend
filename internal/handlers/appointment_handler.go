package handlers

import (
	"net/http"
	"sterling-hms-backend/internal/config"
	"sterling-hms-backend/internal/models"
	"sterling-hms-backend/internal/repositories"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	appointmentRepo *repositories.AppointmentRepository
}

func NewAppointmentHandler(cfg *config.Config) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentRepo: repositories.NewAppointmentRepository(config.DB),
	}
}

// GetPatientProfile retrieves the patient's profile information
func (h *AppointmentHandler) GetPatientProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	patientID := userID.(int)

	profile, err := h.appointmentRepo.GetPatientProfile(patientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Patient profile not found"})
		return
	}

	c.JSON(http.StatusOK, models.PatientProfileResponse{
		Message: "Patient profile retrieved successfully",
		Profile: profile,
	})
}

// GetAppointmentHistory retrieves appointment history for the patient
func (h *AppointmentHandler) GetAppointmentHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	patientID := userID.(int)

	// Get pagination parameters
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	appointments, total, err := h.appointmentRepo.GetAppointmentHistory(patientID, limitNum, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve appointment history"})
		return
	}

	if appointments == nil {
		appointments = []models.Appointment{}
	}

	c.JSON(http.StatusOK, models.AppointmentsListResponse{
		Message:      "Appointment history retrieved successfully",
		Appointments: appointments,
		Total:        total,
	})
}

// GetAvailableSlots retrieves available appointment slots
func (h *AppointmentHandler) GetAvailableSlots(c *gin.Context) {
	doctorID := c.Query("doctorId")
	if doctorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "doctorId query parameter is required"})
		return
	}

	doctorIDInt, err := strconv.Atoi(doctorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid doctorId"})
		return
	}

	// Calculate date range: today to next 30 days
	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().AddDate(0, 0, 30).Format("2006-01-02")

	slots, err := h.appointmentRepo.GetAvailableSlots(doctorIDInt, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve available slots"})
		return
	}

	if slots == nil {
		slots = []models.AppointmentSlot{}
	}

	c.JSON(http.StatusOK, models.AvailableSlotsResponse{
		Message: "Available slots retrieved successfully",
		Slots:   slots,
	})
}

// GetDoctors retrieves all available doctors
func (h *AppointmentHandler) GetDoctors(c *gin.Context) {
	doctors, err := h.appointmentRepo.GetDoctors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve doctors"})
		return
	}

	if doctors == nil {
		doctors = []models.Doctor{}
	}

	c.JSON(http.StatusOK, models.DoctorsResponse{
		Message: "Doctors retrieved successfully",
		Doctors: doctors,
	})
}

// GetSpecializations retrieves all available specializations
func (h *AppointmentHandler) GetSpecializations(c *gin.Context) {
	specializations, err := h.appointmentRepo.GetSpecializations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve specializations"})
		return
	}

	if specializations == nil {
		specializations = []models.Specialization{}
	}

	c.JSON(http.StatusOK, models.SpecializationsResponse{
		Message:         "Specializations retrieved successfully",
		Specializations: specializations,
	})
}

// GetDoctorsBySpecialization retrieves doctors for a specific specialization
func (h *AppointmentHandler) GetDoctorsBySpecialization(c *gin.Context) {
	specialization := c.Query("specialization")
	if specialization == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "specialization query parameter is required"})
		return
	}

	doctors, err := h.appointmentRepo.GetDoctorsBySpecialization(specialization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve doctors"})
		return
	}

	if doctors == nil {
		doctors = []models.Doctor{}
	}

	c.JSON(http.StatusOK, models.DoctorsResponse{
		Message: "Doctors retrieved successfully",
		Doctors: doctors,
	})
}

// BookAppointment creates a new appointment
func (h *AppointmentHandler) BookAppointment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	patientID := userID.(int)

	var req models.BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	// Validate appointment date format
	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Check if appointment date is at least today (not in the past)
	// Compare only the date part, not the exact time
	today := time.Now().Truncate(24 * time.Hour)
	appointmentDateTruncated := appointmentDate.Truncate(24 * time.Hour)

	if appointmentDateTruncated.Before(today) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Appointment date must be today or in the future"})
		return
	}

	// Validate doctor exists
	doctor, err := h.appointmentRepo.GetDoctorByID(req.DoctorID)
	if err != nil || doctor == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid doctor ID"})
		return
	}

	if !doctor.IsAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Selected doctor is not available"})
		return
	}

	// Check slot availability
	isAvailable, err := h.appointmentRepo.CheckSlotAvailability(req.DoctorID, req.AppointmentDate, req.TimeSlot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to check slot availability"})
		return
	}

	if !isAvailable {
		c.JSON(http.StatusConflict, gin.H{"message": "Selected time slot is not available"})
		return
	}

	// Create appointment
	appointment, err := h.appointmentRepo.CreateAppointment(
		patientID,
		req.DoctorID,
		req.AppointmentDate,
		req.TimeSlot,
		req.Reason,
		req.Notes,
	)

	if err != nil {
		// Check if this is a duplicate key error (slot already booked)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			c.JSON(http.StatusConflict, gin.H{"message": "This time slot is no longer available. Please select another slot."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create appointment"})
		return
	}

	c.JSON(http.StatusCreated, models.AppointmentResponse{
		Message:     "Appointment request submitted. Waiting for admin approval.",
		Appointment: appointment,
	})
}

// CancelAppointment cancels an existing appointment
func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	patientID := userID.(int)

	appointmentIDStr := c.Param("id")
	appointmentID, err := strconv.Atoi(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid appointment ID"})
		return
	}

	// Verify appointment belongs to patient
	appointment, err := h.appointmentRepo.GetAppointmentByID(appointmentID)
	if err != nil || appointment == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Appointment not found"})
		return
	}

	if appointment.PatientID != patientID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You can only cancel your own appointments"})
		return
	}

	// Cancel the appointment
	err = h.appointmentRepo.CancelAppointment(appointmentID, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to cancel appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled successfully"})
}
