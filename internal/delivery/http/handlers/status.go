package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Success 201 {object} entities.Status
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /statuses [post]
func (h *StatusHandler) CreateStatus(c *gin.Context) {
	var status entities.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.statusUseCase.CreateStatus(c.Request.Context(), &status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, status)
}

// GetStatus godoc
// @Summary Get status by ID
// @Description Get a specific status by its ID
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Status ID"
// @Success 200 {object} entities.Status
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /statuses/{id} [get]
func (h *StatusHandler) GetStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status id"})
		return
	}

	status, err := h.statusUseCase.GetStatus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "status not found"})
		return
	}

	c.JSON(http.StatusOK, status)
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
// @Success 200 {object} entities.Status
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /statuses/{id} [put]
func (h *StatusHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status id"})
		return
	}

	var status entities.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status.ID = id

	if err := h.statusUseCase.UpdateStatus(c.Request.Context(), &status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// DeleteStatus godoc
// @Summary Delete a status
// @Description Delete a status by its ID
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Status ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /statuses/{id} [delete]
func (h *StatusHandler) DeleteStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status id"})
		return
	}

	if err := h.statusUseCase.DeleteStatus(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status deleted successfully"})
}

// ListStatuses godoc
// @Summary List all statuses with pagination
// @Description Get a paginated list of all statuses
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} map[string][]entities.Status
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /statuses [get]
func (h *StatusHandler) ListStatuses(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	statuses, err := h.statusUseCase.ListStatuses(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statuses": statuses})
}

// GetStatusesByType godoc
// @Summary Get statuses by type
// @Description Get all active statuses filtered by type
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Param type path string true "Status type"
// @Success 200 {object} map[string][]entities.Status
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /statuses/type/{type} [get]
func (h *StatusHandler) GetStatusesByType(c *gin.Context) {
	statusType := entities.StatusType(c.Param("type"))
	
	statuses, err := h.statusUseCase.GetActiveStatusesByType(c.Request.Context(), statusType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statuses": statuses})
}

// GetDiagnosisStatuses godoc
// @Summary Get all diagnosis statuses
// @Description Get all statuses of diagnosis type
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]entities.Status
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /statuses/diagnosis [get]
func (h *StatusHandler) GetDiagnosisStatuses(c *gin.Context) {
	statuses, err := h.statusUseCase.GetDiagnosisStatuses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statuses": statuses})
}

// GetTreatmentStatuses godoc
// @Summary Get all treatment statuses
// @Description Get all statuses of treatment type
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]entities.Status
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /statuses/treatment [get]
func (h *StatusHandler) GetTreatmentStatuses(c *gin.Context) {
	statuses, err := h.statusUseCase.GetTreatmentStatuses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statuses": statuses})
}

// GetToothStatuses godoc
// @Summary Get all tooth statuses
// @Description Get all statuses of tooth type
// @Tags statuses
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]entities.Status
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /statuses/tooth [get]
func (h *StatusHandler) GetToothStatuses(c *gin.Context) {
	statuses, err := h.statusUseCase.GetToothStatuses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statuses": statuses})
}