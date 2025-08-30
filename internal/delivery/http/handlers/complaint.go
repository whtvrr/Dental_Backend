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

type ComplaintHandler struct {
	complaintUseCase *usecases.ComplaintUseCase
}

func NewComplaintHandler(complaintUseCase *usecases.ComplaintUseCase) *ComplaintHandler {
	return &ComplaintHandler{
		complaintUseCase: complaintUseCase,
	}
}

// CreateComplaint godoc
// @Summary Create a new complaint
// @Description Create a new patient complaint
// @Tags complaints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param complaint body entities.Complaint true "Complaint data"
// @Success 201 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /complaints [post]
func (h *ComplaintHandler) CreateComplaint(c *gin.Context) {
	var complaint entities.Complaint
	if err := c.ShouldBindJSON(&complaint); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	if err := h.complaintUseCase.CreateComplaint(c.Request.Context(), &complaint); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Created("Complaint created successfully", complaint))
}

// GetComplaint godoc
// @Summary Get complaint by ID
// @Description Get a specific complaint by its ID
// @Tags complaints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Complaint ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 404 {object} response.StandardResponse "Not Found"
// @Router /complaints/{id} [get]
func (h *ComplaintHandler) GetComplaint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid complaint id"))
		return
	}

	complaint, err := h.complaintUseCase.GetComplaint(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "complaint not found"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Complaint retrieved successfully", complaint))
}

// UpdateComplaint godoc
// @Summary Update an existing complaint
// @Description Update a complaint with the provided details
// @Tags complaints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Complaint ID"
// @Param complaint body entities.Complaint true "Updated complaint data"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /complaints/{id} [put]
func (h *ComplaintHandler) UpdateComplaint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid complaint id"))
		return
	}

	var complaint entities.Complaint
	if err := c.ShouldBindJSON(&complaint); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}
	complaint.ID = id

	if err := h.complaintUseCase.UpdateComplaint(c.Request.Context(), &complaint); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Complaint updated successfully", complaint))
}

// DeleteComplaint godoc
// @Summary Delete a complaint
// @Description Delete a complaint by its ID
// @Tags complaints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Complaint ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /complaints/{id} [delete]
func (h *ComplaintHandler) DeleteComplaint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid complaint id"))
		return
	}

	if err := h.complaintUseCase.DeleteComplaint(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Complaint deleted successfully", nil))
}

// ListComplaints godoc
// @Summary List all complaints with pagination
// @Description Get a paginated list of all complaints
// @Tags complaints
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /complaints [get]
func (h *ComplaintHandler) ListComplaints(c *gin.Context) {
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

	complaints, err := h.complaintUseCase.ListComplaints(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Complaints retrieved successfully", gin.H{"complaints": complaints}))
}

// GetActiveComplaints godoc
// @Summary Get all active complaints
// @Description Get a list of all active complaints
// @Tags complaints
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.StandardResponse
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /complaints/active [get]
func (h *ComplaintHandler) GetActiveComplaints(c *gin.Context) {
	complaints, err := h.complaintUseCase.GetActiveComplaints(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Active complaints retrieved successfully", gin.H{"complaints": complaints}))
}