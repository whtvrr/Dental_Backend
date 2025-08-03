package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *ComplaintHandler) CreateComplaint(c *gin.Context) {
	var complaint entities.Complaint
	if err := c.ShouldBindJSON(&complaint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.complaintUseCase.CreateComplaint(c.Request.Context(), &complaint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, complaint)
}

func (h *ComplaintHandler) GetComplaint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid complaint id"})
		return
	}

	complaint, err := h.complaintUseCase.GetComplaint(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "complaint not found"})
		return
	}

	c.JSON(http.StatusOK, complaint)
}

func (h *ComplaintHandler) UpdateComplaint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid complaint id"})
		return
	}

	var complaint entities.Complaint
	if err := c.ShouldBindJSON(&complaint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	complaint.ID = id

	if err := h.complaintUseCase.UpdateComplaint(c.Request.Context(), &complaint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, complaint)
}

func (h *ComplaintHandler) DeleteComplaint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid complaint id"})
		return
	}

	if err := h.complaintUseCase.DeleteComplaint(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "complaint deleted successfully"})
}

func (h *ComplaintHandler) ListComplaints(c *gin.Context) {
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

	complaints, err := h.complaintUseCase.ListComplaints(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"complaints": complaints})
}

func (h *ComplaintHandler) GetActiveComplaints(c *gin.Context) {
	complaints, err := h.complaintUseCase.GetActiveComplaints(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"complaints": complaints})
}