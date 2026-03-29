package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        int        `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	LastLogin *time.Time `json:"lastLogin"`
	IsActive  bool       `json:"isActive"`
}

type SignUpRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user"`
	Token   string `json:"token"`
}

type UserResponse struct {
	User *User `json:"user"`
}

type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Password Reset Models
type PasswordResetToken struct {
	ID        int        `json:"id"`
	UserID    int        `json:"userId"`
	TokenHash string     `json:"-"` // Not exposed to client
	ExpiresAt time.Time  `json:"expiresAt"`
	UsedAt    *time.Time `json:"usedAt"`
	IsUsed    bool       `json:"isUsed"`
	CreatedAt time.Time  `json:"createdAt"`
}

type PasswordResetLog struct {
	ID           int       `json:"id"`
	UserID       *int      `json:"userId"`
	Email        string    `json:"email"`
	Action       string    `json:"action"` // forgot_password_request, reset_success, reset_failed
	IPAddress    string    `json:"ipAddress"`
	UserAgent    string    `json:"userAgent"`
	Success      bool      `json:"success"`
	ErrorMessage *string   `json:"errorMessage"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Request/Response DTOs for Password Reset
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ResetPasswordRequest struct {
	ResetToken string `json:"resetToken" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// For email links
type PasswordResetLink struct {
	Token     string
	ExpiresAt time.Time
}

// Specialization Model
type Specialization struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Doctor Model
type Doctor struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Specialization  string    `json:"specialization"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	ExperienceYears int       `json:"experienceYears"`
	Qualification   string    `json:"qualification"`
	Address         string    `json:"address"`
	IsAvailable     bool      `json:"isAvailable"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// Appointment Model
type Appointment struct {
	ID               int        `json:"id"`
	PatientID        int        `json:"patientId"`
	DoctorID         int        `json:"doctorId"`
	AppointmentDate  string     `json:"appointmentDate"` // YYYY-MM-DD format
	TimeSlot         string     `json:"timeSlot"`
	Reason           string     `json:"reason"`
	Status           string     `json:"status"`         // scheduled, completed, cancelled, no-show
	ApprovalStatus   string     `json:"approvalStatus"` // pending, approved, rejected
	Notes            string     `json:"notes"`
	ApprovedBy       *int       `json:"approvedBy,omitempty"`
	ApprovedAt       *time.Time `json:"approvedAt,omitempty"`
	RejectionReason  string     `json:"rejectionReason,omitempty"`
	PatientFirstName string     `json:"patientFirstName,omitempty"`
	PatientLastName  string     `json:"patientLastName,omitempty"`
	PatientEmail     string     `json:"patientEmail,omitempty"`
	PatientPhone     string     `json:"patientPhone,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	Doctor           *Doctor    `json:"doctor,omitempty"`
}

// AppointmentSlot Model
type AppointmentSlot struct {
	ID          int       `json:"id"`
	DoctorID    int       `json:"doctorId"`
	SlotDate    string    `json:"slotDate"` // YYYY-MM-DD format
	TimeSlot    string    `json:"timeSlot"`
	IsAvailable bool      `json:"isAvailable"`
	CreatedAt   time.Time `json:"createdAt"`
	Doctor      *Doctor   `json:"doctor,omitempty"`
}

// Patient Profile Model (extends User with medical info)
type PatientProfile struct {
	ID             int    `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	MedicalHistory string `json:"medicalHistory"`
}

// Request DTOs for Appointments
type BookAppointmentRequest struct {
	DoctorID        int    `json:"doctorId" binding:"required"`
	AppointmentDate string `json:"appointmentDate" binding:"required"` // YYYY-MM-DD
	TimeSlot        string `json:"timeSlot" binding:"required"`
	Reason          string `json:"reason" binding:"required"`
	Notes           string `json:"notes"`
}

type UpdateAppointmentRequest struct {
	Status string `json:"status" binding:"required"` // completed, cancelled, no-show
	Notes  string `json:"notes"`
}

// Response DTOs
type AppointmentResponse struct {
	Message     string       `json:"message"`
	Appointment *Appointment `json:"appointment"`
}

type AppointmentsListResponse struct {
	Message      string        `json:"message"`
	Appointments []Appointment `json:"appointments"`
	Total        int           `json:"total"`
}

type AvailableSlotsResponse struct {
	Message string            `json:"message"`
	Slots   []AppointmentSlot `json:"slots"`
}

type SpecializationsResponse struct {
	Message         string           `json:"message"`
	Specializations []Specialization `json:"specializations"`
}

type DoctorsResponse struct {
	Message string   `json:"message"`
	Doctors []Doctor `json:"doctors"`
}

type PatientProfileResponse struct {
	Message string          `json:"message"`
	Profile *PatientProfile `json:"profile"`
}

// Admin Models
type AdminUser struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	IsActive     bool      `json:"isActive"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Message string     `json:"message"`
	Success bool       `json:"success"`
	Admin   *AdminUser `json:"admin"`
	Token   string     `json:"token"`
}

type AdminDashboardStats struct {
	TotalPatients     int `json:"totalPatients"`
	TotalAppointments int `json:"totalAppointments"`
	TotalDoctors      int `json:"totalDoctors"`
	TotalStaff        int `json:"totalStaff"`
}

type AdminDashboardResponse struct {
	Message string              `json:"message"`
	Success bool                `json:"success"`
	Stats   AdminDashboardStats `json:"stats"`
}

type AdminClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// Receptionist Models
type ReceptionistUser struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Department   string    `json:"department"`
	Role         string    `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	IsActive     bool      `json:"isActive"`
}

type ReceptionistLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ReceptionistLoginResponse struct {
	Message      string            `json:"message"`
	Success      bool              `json:"success"`
	Receptionist *ReceptionistUser `json:"receptionist"`
	Token        string            `json:"token"`
}

type ReceptionistClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// Patient Record Model
type PatientRecord struct {
	ID                    int        `json:"id"`
	UserID                int        `json:"userId"`
	FirstName             string     `json:"firstName"`
	LastName              string     `json:"lastName"`
	Email                 string     `json:"email"`
	Phone                 string     `json:"phone"`
	DateOfBirth           *time.Time `json:"dateOfBirth"`
	Gender                string     `json:"gender"`
	BloodType             string     `json:"bloodType"`
	Address               string     `json:"address"`
	City                  string     `json:"city"`
	State                 string     `json:"state"`
	PostalCode            string     `json:"postalCode"`
	Country               string     `json:"country"`
	Allergies             string     `json:"allergies"`
	MedicalConditions     string     `json:"medicalConditions"`
	CurrentMedications    string     `json:"currentMedications"`
	EmergencyContactName  string     `json:"emergencyContactName"`
	EmergencyContactPhone string     `json:"emergencyContactPhone"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
}

// Register Patient Request (used by receptionist)
type RegisterPatientRequest struct {
	FirstName             string `json:"firstName" binding:"required"`
	LastName              string `json:"lastName" binding:"required"`
	Email                 string `json:"email" binding:"required,email"`
	Password              string `json:"password" binding:"required,min=8"`
	Phone                 string `json:"phone" binding:"required"`
	DateOfBirth           string `json:"dateOfBirth"` // YYYY-MM-DD
	Gender                string `json:"gender"`
	BloodType             string `json:"bloodType"`
	Address               string `json:"address"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	PostalCode            string `json:"postalCode"`
	Country               string `json:"country"`
	Allergies             string `json:"allergies"`
	MedicalConditions     string `json:"medicalConditions"`
	CurrentMedications    string `json:"currentMedications"`
	EmergencyContactName  string `json:"emergencyContactName"`
	EmergencyContactPhone string `json:"emergencyContactPhone"`
}

// Appointment Approval Request
type ApproveAppointmentRequest struct {
	Status string `json:"status" binding:"required"` // approved or rejected
	Reason string `json:"reason"`                    // For rejection
}

// Patient Registration Response
type PatientRegistrationResponse struct {
	Message string         `json:"message"`
	Success bool           `json:"success"`
	Patient *PatientRecord `json:"patient"`
	User    *User          `json:"user"`
}

// Pending Appointments Response
type PendingAppointmentsResponse struct {
	Message      string        `json:"message"`
	Appointments []Appointment `json:"appointments"`
	Total        int           `json:"total"`
}

// Receptionist Dashboard Stats
type ReceptionistDashboardStats struct {
	TotalPatients        int `json:"totalPatients"`
	PendingAppointments  int `json:"pendingAppointments"`
	ApprovedAppointments int `json:"approvedAppointments"`
	RejectedAppointments int `json:"rejectedAppointments"`
}

type ReceptionistDashboardResponse struct {
	Message string                     `json:"message"`
	Success bool                       `json:"success"`
	Stats   ReceptionistDashboardStats `json:"stats"`
}

// Doctor Models
type DoctorUser struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	PasswordHash   string    `json:"-"`
	Specialization string    `json:"specialization"`
	Phone          string    `json:"phone"`
	Role           string    `json:"role"`
	IsActive       bool      `json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type DoctorLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type DoctorLoginResponse struct {
	Message  string      `json:"message"`
	Success  bool        `json:"success"`
	Doctor   *DoctorUser `json:"doctor"`
	Token    string      `json:"token"`
	Role     string      `json:"role"`
	DoctorID int         `json:"doctorId"`
}

type DoctorClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// Patient info for doctor
type PatientInfo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	MedicalHistory string `json:"medicalHistory"`
	Address        string `json:"address"`
	BloodType      string `json:"bloodType"`
}

// Appointment info for doctor
type AppointmentInfoForDoctor struct {
	ID              int    `json:"id"`
	PatientName     string `json:"patientName"`
	PatientID       int    `json:"patientId"`
	AppointmentDate string `json:"appointmentDate"`
	TimeSlot        string `json:"timeSlot"`
	Reason          string `json:"reason"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
	PatientPhone    string `json:"patientPhone"`
	PatientEmail    string `json:"patientEmail"`
}

type DoctorPatientsResponse struct {
	Message  string        `json:"message"`
	Success  bool          `json:"success"`
	Patients []PatientInfo `json:"patients"`
	Total    int           `json:"total"`
}

type DoctorAppointmentsResponse struct {
	Message      string                     `json:"message"`
	Success      bool                       `json:"success"`
	Appointments []AppointmentInfoForDoctor `json:"appointments"`
	Total        int                        `json:"total"`
}

type UpdateAppointmentStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

type UpdateAppointmentStatusResponse struct {
	Message     string                    `json:"message"`
	Success     bool                      `json:"success"`
	Appointment *AppointmentInfoForDoctor `json:"appointment"`
}

type DoctorDashboardStats struct {
	TotalPatients         int `json:"totalPatients"`
	TotalAppointments     int `json:"totalAppointments"`
	CompletedAppointments int `json:"completedAppointments"`
	UpcomingAppointments  int `json:"upcomingAppointments"`
}

type DoctorDashboardResponse struct {
	Message string               `json:"message"`
	Success bool                 `json:"success"`
	Stats   DoctorDashboardStats `json:"stats"`
}
