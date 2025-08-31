package http

import (
	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/handlers"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/middleware"
)

type Handlers struct {
	Auth        *handlers.AuthHandler
	User        *handlers.UserHandler
	Appointment *handlers.AppointmentHandler
	Status      *handlers.StatusHandler
	Complaint   *handlers.ComplaintHandler
	Formula     *handlers.FormulaHandler
}

func SetupRoutes(router *gin.Engine, h *Handlers, authMiddleware *middleware.AuthMiddleware) {
	api := router.Group("/api/v1")

	// Public authentication routes
	auth := api.Group("/auth")
	{
		auth.POST("/signup", h.Auth.SignUp)
		auth.POST("/signin", h.Auth.SignIn)
		auth.POST("/refresh", h.Auth.RefreshToken)
	}

	// Protected routes (require authentication)
	protected := api.Group("/")
	protected.Use(authMiddleware.RequireAuth())
	{
		// Profile route
		protected.GET("/me", h.Auth.Me)

		// User routes with RBAC
		users := protected.Group("/users")
		{
			// Admin only: Full CRUD on all users
			users.POST("", authMiddleware.RequireAdmin(), h.User.CreateUser)
			users.DELETE("/:id", authMiddleware.RequireAdmin(), h.User.DeleteUser)
			users.GET("", authMiddleware.RequireAdmin(), h.User.ListUsers)
			
			// Doctors: Can view user details for patient cards
			// Receptionists: Can only view clients
			users.GET("/:id", authMiddleware.RequireStaff(), h.User.GetUser)
			users.PUT("/:id", authMiddleware.RequireStaff(), h.User.UpdateUser)
			users.GET("/doctors", authMiddleware.RequireStaff(), h.User.GetDoctors)
			users.GET("/staff", authMiddleware.RequireStaff(), h.User.GetStaff)
			users.GET("/clients", authMiddleware.RequireStaff(), h.User.GetClients)
		}

		// Appointment routes with RBAC
		appointments := protected.Group("/appointments")
		{
			// Receptionists: Can create, update, delete appointments
			appointments.POST("", authMiddleware.RequireDoctorOrReceptionist(), h.Appointment.CreateAppointment)
			appointments.PUT("/:id", authMiddleware.RequireDoctorOrReceptionist(), h.Appointment.UpdateAppointment)
			appointments.DELETE("/:id", authMiddleware.RequireDoctorOrReceptionist(), h.Appointment.DeleteAppointment)
			appointments.POST("/:id/cancel", authMiddleware.RequireDoctorOrReceptionist(), h.Appointment.CancelAppointment)
			
			// Doctors: Can complete appointments (add medical data)
			appointments.POST("/:id/complete", authMiddleware.RequireDoctor(), h.Appointment.CompleteAppointment)
			
			// All staff: Can view appointments
			appointments.GET("/:id", authMiddleware.RequireStaff(), h.Appointment.GetAppointment)
			appointments.GET("", authMiddleware.RequireStaff(), h.Appointment.ListAppointments)
			appointments.GET("/doctor/:doctorId", authMiddleware.RequireStaff(), h.Appointment.GetDoctorAppointments)
			appointments.GET("/client/:clientId", authMiddleware.RequireStaff(), h.Appointment.GetClientAppointments)
		}

		// Status routes with RBAC
		statuses := protected.Group("/statuses")
		{
			// Doctors: Can CRUD statuses (for diagnosis, treatment, tooth conditions)
			statuses.POST("", authMiddleware.RequireDoctor(), h.Status.CreateStatus)
			statuses.PUT("/:id", authMiddleware.RequireDoctor(), h.Status.UpdateStatus)
			statuses.DELETE("/:id", authMiddleware.RequireDoctor(), h.Status.DeleteStatus)
			
			// All staff: Can view statuses
			statuses.GET("/:id", authMiddleware.RequireStaff(), h.Status.GetStatus)
			statuses.GET("", authMiddleware.RequireStaff(), h.Status.ListStatuses)
			statuses.GET("/type/:type", authMiddleware.RequireStaff(), h.Status.GetStatusesByType)
			statuses.GET("/diagnosis", authMiddleware.RequireStaff(), h.Status.GetDiagnosisStatuses)
			statuses.GET("/treatment", authMiddleware.RequireStaff(), h.Status.GetTreatmentStatuses)
			statuses.GET("/tooth", authMiddleware.RequireStaff(), h.Status.GetToothStatuses)
		}

		// Complaint routes with RBAC
		complaints := protected.Group("/complaints")
		{
			// Doctors: Can CRUD complaints
			complaints.POST("", authMiddleware.RequireDoctor(), h.Complaint.CreateComplaint)
			complaints.PUT("/:id", authMiddleware.RequireDoctor(), h.Complaint.UpdateComplaint)
			complaints.DELETE("/:id", authMiddleware.RequireDoctor(), h.Complaint.DeleteComplaint)
			
			// All staff: Can view complaints
			complaints.GET("/:id", authMiddleware.RequireStaff(), h.Complaint.GetComplaint)
			complaints.GET("", authMiddleware.RequireStaff(), h.Complaint.ListComplaints)
		}

		// Formula routes with RBAC
		formulas := protected.Group("/formulas")
		{
			// Doctors: Can CRUD formulas (dental records)
			formulas.GET("/:id", authMiddleware.RequireDoctor(), h.Formula.GetFormula)
			formulas.GET("/user/:userId", authMiddleware.RequireDoctor(), h.Formula.GetFormulaByUserID)
		}
	}
}
