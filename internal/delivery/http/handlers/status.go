package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/response"
	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusHandler struct {
	statusUseCase *usecases.StatusUseCase
}

func NewStatusHandler(statusUseCase *usecases.StatusUseCase) *StatusHandler {
	return &StatusHandler{
		statusUseCase: statusUseCase,
	}
}

// CreateStatus godoc
// @Summary Create a new status
// @Description Create a new dental status
// @Tags statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status body entities.Status true "Status data"
// @Success 201 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /statuses [post]
func (h *StatusHandler) CreateStatus(c *gin.Context) {
	var status entities.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	if err := h.statusUseCase.CreateStatus(c.Request.Context(), &status); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Created("Status created successfully", status))
}

// GetStatus godoc
// @Summary Get status by ID
// @Description Get a specific status by its ID
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Status ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 404 {object} response.StandardResponse "Not Found"
// @Router /statuses/{id} [get]
func (h *StatusHandler) GetStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid status id"))
		return
	}

	status, err := h.statusUseCase.GetStatus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "status not found"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Status retrieved successfully", status))
}

// UpdateStatus godoc
// @Summary Update an existing status
// @Description Update a status with the provided details
// @Tags statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Status ID"
// @Param status body entities.Status true "Updated status data"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /statuses/{id} [put]
func (h *StatusHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid status id"))
		return
	}

	var status entities.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}
	status.ID = id

	if err := h.statusUseCase.UpdateStatus(c.Request.Context(), &status); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Status updated successfully", status))
}

// DeleteStatus godoc
// @Summary Delete a status
// @Description Delete a status by its ID
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Status ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /statuses/{id} [delete]
func (h *StatusHandler) DeleteStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid status id"))
		return
	}

	if err := h.statusUseCase.DeleteStatus(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Status deleted successfully", nil))
}

// ListStatuses godoc
// @Summary List all statuses with pagination
// @Description Get a paginated list of all statuses
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /statuses [get]
func (h *StatusHandler) ListStatuses(c *gin.Context) {
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

	statuses, err := h.statusUseCase.ListStatuses(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Statuses retrieved successfully", gin.H{"statuses": statuses}))
}

// GetStatusesByType godoc
// @Summary Get statuses by type
// @Description Get all active statuses filtered by type
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param type path string true "Status type"
// @Success 200 {object} response.StandardResponse
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /statuses/type/{type} [get]
func (h *StatusHandler) GetStatusesByType(c *gin.Context) {
	statusType := entities.StatusType(c.Param("type"))
	
	statuses, err := h.statusUseCase.GetActiveStatusesByType(c.Request.Context(), statusType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Statuses retrieved successfully", gin.H{"statuses": statuses}))
}

// GetDiagnosisStatuses godoc
// @Summary Get all diagnosis statuses with pagination
// @Description Get paginated list of diagnosis statuses
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /statuses/diagnosis [get]
func (h *StatusHandler) GetDiagnosisStatuses(c *gin.Context) {
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

	statuses, err := h.statusUseCase.GetDiagnosisStatusesWithPagination(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Diagnosis statuses retrieved successfully", gin.H{"statuses": statuses}))
}

// GetTreatmentStatuses godoc
// @Summary Get all treatment statuses with pagination
// @Description Get paginated list of treatment statuses
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /statuses/treatment [get]
func (h *StatusHandler) GetTreatmentStatuses(c *gin.Context) {
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

	statuses, err := h.statusUseCase.GetTreatmentStatusesWithPagination(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Treatment statuses retrieved successfully", gin.H{"statuses": statuses}))
}

// GetToothStatuses godoc
// @Summary Get all tooth statuses with pagination
// @Description Get paginated list of tooth statuses
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /statuses/tooth [get]
func (h *StatusHandler) GetToothStatuses(c *gin.Context) {
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

	statuses, err := h.statusUseCase.GetToothStatusesWithPagination(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Tooth statuses retrieved successfully", gin.H{"statuses": statuses}))
}