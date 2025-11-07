package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/response"
	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentCreateRequest struct {
	DateTime        time.Time          `json:"date_time" binding:"required"`
	DoctorID        primitive.ObjectID `json:"doctor_id" binding:"required"`
	ClientID        primitive.ObjectID `json:"client_id" binding:"required"`
	DurationMinutes *int               `json:"duration_minutes,omitempty"`
	Status          string             `json:"status,omitempty"`
	Comment         *string            `json:"comment,omitempty"`
}

type AppointmentHandler struct {
	appointmentUseCase *usecases.AppointmentUseCase
}

func NewAppointmentHandler(appointmentUseCase *usecases.AppointmentUseCase) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentUseCase: appointmentUseCase,
	}
}

// CreateAppointment godoc
// @Summary Create a new appointment
// @Description Create a new appointment with the provided details. Duration defaults to 30 minutes if not provided.
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param appointment body AppointmentCreateRequest true "Appointment data"
// @Success 201 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /appointments [post]
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req AppointmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	// Set default duration to 30 minutes if not provided
	durationMinutes := 30
	if req.DurationMinutes != nil {
		durationMinutes = *req.DurationMinutes
	}

	// Convert request to appointment entity
	appointment := &entities.Appointment{
		DateTime:        req.DateTime,
		DoctorID:        req.DoctorID,
		ClientID:        req.ClientID,
		DurationMinutes: durationMinutes,
		Status:          entities.AppointmentStatus(req.Status),
		Comment:         req.Comment,
	}

	if err := h.appointmentUseCase.CreateAppointment(c.Request.Context(), appointment); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Created("Appointment created successfully", appointment))
}

// GetAppointment godoc
// @Summary Get appointment by ID
// @Description Get a specific appointment by its ID
// @Tags appointments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 404 {object} response.StandardResponse "Not Found"
// @Router /appointments/{id} [get]
func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid appointment id"))
		return
	}

	appointment, err := h.appointmentUseCase.GetAppointment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "appointment not found"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Appointment retrieved successfully", appointment))
}

// UpdateAppointment godoc
// @Summary Update an existing appointment
// @Description Update an appointment with the provided details
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Param appointment body entities.Appointment true "Updated appointment data"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /appointments/{id} [put]
func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid appointment id"))
		return
	}

	var appointment entities.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}
	appointment.ID = id

	if err := h.appointmentUseCase.UpdateAppointment(c.Request.Context(), &appointment); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Appointment updated successfully", appointment))
}

// CompleteAppointment godoc
// @Summary Complete an appointment with medical data
// @Description Mark an appointment as completed and update with medical information
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Param medicalData body usecases.AppointmentMedicalData true "Medical data for appointment"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /appointments/{id}/complete [post]
func (h *AppointmentHandler) CompleteAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid appointment id"))
		return
	}

	var medicalData usecases.AppointmentMedicalData
	if err := c.ShouldBindJSON(&medicalData); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	if err := h.appointmentUseCase.CompleteAppointment(c.Request.Context(), id, &medicalData); err != nil {
		// Check if it's a validation error
		if usecases.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Appointment completed successfully", nil))
}

// CancelAppointment godoc
// @Summary Cancel an appointment
// @Description Mark an appointment as canceled
// @Tags appointments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /appointments/{id}/cancel [post]
func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid appointment id"))
		return
	}

	if err := h.appointmentUseCase.CancelAppointment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Appointment canceled successfully", nil))
}

// DeleteAppointment godoc
// @Summary Delete an appointment
// @Description Delete an appointment by its ID
// @Tags appointments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /appointments/{id} [delete]
func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid appointment id"))
		return
	}

	if err := h.appointmentUseCase.DeleteAppointment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Appointment deleted successfully", nil))
}

// ListAppointments godoc
// @Summary List all appointments with pagination
// @Description Get a paginated list of all appointments
// @Tags appointments
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /appointments [get]
func (h *AppointmentHandler) ListAppointments(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid offset"))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid limit"))
		return
	}

	appointments, err := h.appointmentUseCase.ListAppointments(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Appointments retrieved successfully", gin.H{"appointments": appointments}))
}

// GetDoctorAppointments godoc
// @Summary Get appointments for a specific doctor
// @Description Get all appointments for a doctor within a date range
// @Tags appointments
// @Produce json
// @Security BearerAuth
// @Param doctorId path string true "Doctor ID"
// @Param from query string true "Start date (DD.MM.YYYY)"
// @Param to query string true "End date (DD.MM.YYYY)"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /appointments/doctor/{doctorId} [get]
func (h *AppointmentHandler) GetDoctorAppointments(c *gin.Context) {
	doctorIDStr := c.Param("doctorId")
	doctorID, err := primitive.ObjectIDFromHex(doctorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid doctor id"))
		return
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := time.Parse("02.01.2006", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid from date format (DD.MM.YYYY)"))
		return
	}

	to, err := time.Parse("02.01.2006", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid to date format (DD.MM.YYYY)"))
		return
	}

	to = to.Add(24*time.Hour - time.Nanosecond)

	appointments, err := h.appointmentUseCase.GetDoctorAppointments(c.Request.Context(), doctorID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Doctor appointments retrieved successfully", gin.H{"appointments": appointments}))
}

// GetClientAppointments godoc
// @Summary Get appointments for a specific client
// @Description Get all appointments for a specific client
// @Tags appointments
// @Produce json
// @Security BearerAuth
// @Param clientId path string true "Client ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /appointments/client/{clientId} [get]
func (h *AppointmentHandler) GetClientAppointments(c *gin.Context) {
	clientIDStr := c.Param("clientId")
	clientID, err := primitive.ObjectIDFromHex(clientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid client id"))
		return
	}

	appointments, err := h.appointmentUseCase.GetClientAppointments(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Client appointments retrieved successfully", gin.H{"appointments": appointments}))
}
