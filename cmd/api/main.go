package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/auth"
	"github.com/whtvrr/Dental_Backend/internal/config"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/handlers"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/middleware"
	"github.com/whtvrr/Dental_Backend/internal/infrastructure/database"
	"github.com/whtvrr/Dental_Backend/internal/infrastructure/persistence/mongodb"
	"github.com/whtvrr/Dental_Backend/internal/usecases"
	logger "github.com/whtvrr/Dental_Backend/pkg"
)

func main() {
	var cfg *config.Config
	cfg, err := config.LoadConfig("config/config.local.yaml")
	if err != nil {
		cfg, err = config.LoadConfig("config/config.yaml")
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	l := logger.GetLogger()
	l.Info("Starting Dental Clinic Management System...")

	db, err := database.NewMongoConnection(cfg)
	if err != nil {
		l.Fatalf("Failed to connect to database: %v", err)
	}
	l.Info("Successfully connected to MongoDB")

	userRepo := mongodb.NewUserRepository(db)
	appointmentRepo := mongodb.NewAppointmentRepository(db)
	statusRepo := mongodb.NewStatusRepository(db)
	complaintRepo := mongodb.NewComplaintRepository(db)
	formulaRepo := mongodb.NewFormulaRepository(db)

	// Parse token TTL
	accessTokenTTL, err := time.ParseDuration(cfg.Auth.TokenTTL)
	if err != nil {
		l.Fatalf("Invalid access token TTL: %v", err)
	}
	refreshTokenTTL := accessTokenTTL * 7 // Refresh token valid for 7x access token duration

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, accessTokenTTL, refreshTokenTTL)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(userRepo, formulaRepo)
	appointmentUseCase := usecases.NewAppointmentUseCase(appointmentRepo, userRepo, formulaRepo)
	statusUseCase := usecases.NewStatusUseCase(statusRepo)
	complaintUseCase := usecases.NewComplaintUseCase(complaintRepo)
	formulaUseCase := usecases.NewFormulaUseCase(formulaRepo)
	authUseCase := usecases.NewAuthUseCase(userRepo, jwtManager)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	userHandler := handlers.NewUserHandler(userUseCase)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentUseCase)
	statusHandler := handlers.NewStatusHandler(statusUseCase)
	complaintHandler := handlers.NewComplaintHandler(complaintUseCase)
	formulaHandler := handlers.NewFormulaHandler(formulaUseCase)

	handlerGroup := &http.Handlers{
		Auth:        authHandler,
		User:        userHandler,
		Appointment: appointmentHandler,
		Status:      statusHandler,
		Complaint:   complaintHandler,
		Formula:     formulaHandler,
	}

	// Setup router
	r := gin.Default()
	l.Info("Router created")

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	http.SetupRoutes(r, handlerGroup, authMiddleware)
	l.Info("Routes configured")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Dental Clinic Management System is running",
		})
	})

	l.Infof("Starting server on port %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		l.Fatalf("Failed to start server: %v", err)
	}
}
