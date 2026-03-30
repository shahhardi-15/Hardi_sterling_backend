package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"sterling-hms-backend/internal/config"
	"sterling-hms-backend/internal/handlers"
	"sterling-hms-backend/internal/middleware"
	"sterling-hms-backend/internal/repositories"
	"sterling-hms-backend/internal/utils"
)

func main() {
	// Load environment variables
	godotenv.Load(".env")

	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	err := cfg.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.DB.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(config.DB)
	adminRepo := repositories.NewAdminRepository(config.DB)
	receptionistRepo := repositories.NewReceptionistRepository(config.DB)
	appointmentRepo := repositories.NewAppointmentRepository(config.DB)
	doctorRepo := repositories.NewDoctorRepository(config.DB)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, cfg)
	adminHandler := handlers.NewAdminHandler(adminRepo, cfg)
	receptionistHandler := handlers.NewReceptionistHandler(receptionistRepo, userRepo, appointmentRepo, cfg)
	appointmentHandler := handlers.NewAppointmentHandler(cfg)
	doctorHandler := handlers.NewDoctorHandler(doctorRepo, cfg)

	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Setup Gin router (avoid duplicate middleware warning by using gin.New())
	router := gin.New()

	// Add Logger and Recovery middleware explicitly
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Set trusted proxies
	router.SetTrustedProxies([]string{"127.0.0.1", "localhost", "::1"})

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Security headers middleware
	router.Use(utils.SecurityHeaders())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running"})
	})

	// Auth routes
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/signup", authHandler.SignUp)
		authGroup.POST("/signin", authHandler.SignIn)
		authGroup.GET("/me", middleware.AuthMiddleware(cfg), authHandler.GetCurrentUser)
		authGroup.GET("/users", authHandler.GetAllUsers)

		// Password reset routes
		authGroup.POST("/forgot-password", authHandler.ForgotPassword)
		authGroup.POST("/reset-password", authHandler.ResetPassword)
	}

	// Patient routes (protected)
	patientGroup := router.Group("/api/patient")
	patientGroup.Use(middleware.AuthMiddleware(cfg))
	{
		patientGroup.GET("/profile", appointmentHandler.GetPatientProfile)
	}

	// Appointment routes (protected)
	appointmentGroup := router.Group("/api/appointments")
	appointmentGroup.Use(middleware.AuthMiddleware(cfg))
	{
		appointmentGroup.GET("/history", appointmentHandler.GetAppointmentHistory)
		appointmentGroup.POST("/book", appointmentHandler.BookAppointment)
		appointmentGroup.DELETE("/:id", appointmentHandler.CancelAppointment)
	}

	// Doctor routes (public)
	doctorGroup := router.Group("/api/doctors")
	{
		doctorGroup.GET("", appointmentHandler.GetDoctors)
		doctorGroup.GET("/available-slots", appointmentHandler.GetAvailableSlots)
		doctorGroup.GET("/by-specialization", appointmentHandler.GetDoctorsBySpecialization)
	}

	// Specialization routes (public)
	specializationGroup := router.Group("/api/specializations")
	{
		specializationGroup.GET("", appointmentHandler.GetSpecializations)
	}

	// Admin routes
	adminGroup := router.Group("/api/admin")
	{
		// Public admin routes
		adminGroup.POST("/login", adminHandler.AdminLogin)

		// Protected admin routes (require admin token)
		adminProtected := adminGroup.Group("")
		adminProtected.Use(middleware.AdminAuthMiddleware(cfg))
		{
			adminProtected.GET("/dashboard/stats", adminHandler.GetDashboardStats)
			adminProtected.POST("/logout", adminHandler.AdminLogout)
		}
	}

	// Receptionist routes
	receptionistGroup := router.Group("/api/receptionist")
	{
		// Public receptionist routes
		receptionistGroup.POST("/login", receptionistHandler.ReceptionistLogin)

		// Protected receptionist routes (require receptionist token)
		receptionistProtected := receptionistGroup.Group("")
		receptionistProtected.Use(middleware.ReceptionistAuthMiddleware(cfg))
		{
			// Dashboard
			receptionistProtected.GET("/dashboard/stats", receptionistHandler.GetDashboardStats)
			receptionistProtected.POST("/logout", receptionistHandler.ReceptionistLogout)

			// Patient management
			receptionistProtected.POST("/patients/register", receptionistHandler.RegisterPatient)
			receptionistProtected.GET("/patients", receptionistHandler.GetPatientRecords)

			// Appointment management
			receptionistProtected.POST("/appointments/book", receptionistHandler.BookAppointmentByReceptionist)
			receptionistProtected.GET("/appointments/pending", receptionistHandler.GetPendingAppointments)
			receptionistProtected.POST("/appointments/:id/approve", receptionistHandler.ApproveAppointment)
			receptionistProtected.POST("/appointments/:id/reject", receptionistHandler.RejectAppointment)
		}
	}

	// Doctor routes
	doctorAuthGroup := router.Group("/api/doctor")
	{
		// Protected doctor routes (require doctor token)
		doctorProtected := doctorAuthGroup.Group("")
		doctorProtected.Use(middleware.DoctorAuthMiddleware(cfg))
		{
			// Dashboard
			doctorProtected.GET("/dashboard/stats", doctorHandler.GetDashboardStats)

			// Patient management
			doctorProtected.GET("/patients", doctorHandler.GetAssignedPatients)

			// Appointment management
			doctorProtected.GET("/appointments", doctorHandler.GetAppointments)
			doctorProtected.PUT("/appointments/:appointmentId/status", doctorHandler.UpdateAppointmentStatus)

			// Logout
			doctorProtected.POST("/logout", doctorHandler.DoctorLogout)

			// Profile
			doctorProtected.GET("/profile", doctorHandler.GetProfile)
		}
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Route not found"})
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server is running on http://localhost:%s", cfg.Port)
	log.Printf("Environment: %s", cfg.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
