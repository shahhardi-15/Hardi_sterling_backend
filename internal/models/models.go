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

// Doctor Model
type Doctor struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Specialization string    `json:"specialization"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	IsAvailable    bool      `json:"isAvailable"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Appointment Model
type Appointment struct {
	ID              int       `json:"id"`
	PatientID       int       `json:"patientId"`
	DoctorID        int       `json:"doctorId"`
	AppointmentDate string    `json:"appointmentDate"` // YYYY-MM-DD format
	TimeSlot        string    `json:"timeSlot"`
	Reason          string    `json:"reason"`
	Status          string    `json:"status"` // scheduled, completed, cancelled, no-show
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Doctor          *Doctor   `json:"doctor,omitempty"`
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
