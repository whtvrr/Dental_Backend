package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/requests"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/response"
	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnamnesisHandler struct {
	anamnesisUseCase *usecases.AnamnesisUseCase
}

func NewAnamnesisHandler(anamnesisUseCase *usecases.AnamnesisUseCase) *AnamnesisHandler {
	return &AnamnesisHandler{
		anamnesisUseCase: anamnesisUseCase,
	}
}

// CreateAnamnesis godoc
// @Summary Create a new anamnesis
// @Description Create a new anamnesis record
// @Tags anamnesis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param anamnesis body requests.CreateAnamnesisRequest true "Anamnesis data"
// @Success 201 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /anamnesis [post]
func (h *AnamnesisHandler) CreateAnamnesis(c *gin.Context) {
	var req requests.CreateAnamnesisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	anamnesis := entities.Anamnesis{
		Text: req.Text,
	}

	if err := h.anamnesisUseCase.CreateAnamnesis(c.Request.Context(), &anamnesis); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Created("Anamnesis created successfully", anamnesis))
}

// GetAnamnesis godoc
// @Summary Get anamnesis by ID
// @Description Get a specific anamnesis by its ID
// @Tags anamnesis
// @Produce json
// @Security BearerAuth
// @Param id path string true "Anamnesis ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 404 {object} response.StandardResponse "Not Found"
// @Router /anamnesis/{id} [get]
func (h *AnamnesisHandler) GetAnamnesis(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid anamnesis id"))
		return
	}

	anamnesis, err := h.anamnesisUseCase.GetAnamnesis(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "anamnesis not found"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Anamnesis retrieved successfully", anamnesis))
}

// UpdateAnamnesis godoc
// @Summary Update an existing anamnesis
// @Description Update an anamnesis with the provided details
// @Tags anamnesis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Anamnesis ID"
// @Param anamnesis body requests.UpdateAnamnesisRequest true "Updated anamnesis data"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /anamnesis/{id} [put]
func (h *AnamnesisHandler) UpdateAnamnesis(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid anamnesis id"))
		return
	}

	var req requests.UpdateAnamnesisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	anamnesis := entities.Anamnesis{
		ID:   id,
		Text: req.Text,
	}

	if err := h.anamnesisUseCase.UpdateAnamnesis(c.Request.Context(), &anamnesis); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Anamnesis updated successfully", anamnesis))
}

// DeleteAnamnesis godoc
// @Summary Delete an anamnesis
// @Description Delete an anamnesis by its ID
// @Tags anamnesis
// @Produce json
// @Security BearerAuth
// @Param id path string true "Anamnesis ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /anamnesis/{id} [delete]
func (h *AnamnesisHandler) DeleteAnamnesis(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid anamnesis id"))
		return
	}

	if err := h.anamnesisUseCase.DeleteAnamnesis(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Anamnesis deleted successfully", nil))
}

// ListAnamnesis godoc
// @Summary List all anamnesis with pagination
// @Description Get a paginated list of all anamnesis records
// @Tags anamnesis
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /anamnesis [get]
func (h *AnamnesisHandler) ListAnamnesis(c *gin.Context) {
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

	anamnesises, err := h.anamnesisUseCase.ListAnamnesis(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Anamnesis retrieved successfully", gin.H{"anamnesis": anamnesises}))
}
