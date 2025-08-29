package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FormulaHandler struct {
	formulaUseCase *usecases.FormulaUseCase
}

func NewFormulaHandler(formulaUseCase *usecases.FormulaUseCase) *FormulaHandler {
	return &FormulaHandler{
		formulaUseCase: formulaUseCase,
	}
}

// GetFormula godoc
// @Summary Get formula by ID
// @Description Get a specific dental formula by its ID
// @Tags formulas
// @Produce json
// @Security BearerAuth
// @Param id path string true "Formula ID"
// @Success 200 {object} entities.Formula
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /formulas/{id} [get]
func (h *FormulaHandler) GetFormula(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid formula id"})
		return
	}

	formula, err := h.formulaUseCase.GetFormula(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "formula not found"})
		return
	}

	c.JSON(http.StatusOK, formula)
}

// GetFormulaByUserID godoc
// @Summary Get formula by user ID
// @Description Get a user's dental formula by their user ID
// @Tags formulas
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Success 200 {object} entities.Formula
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /formulas/user/{userId} [get]
func (h *FormulaHandler) GetFormulaByUserID(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	formula, err := h.formulaUseCase.GetFormulaByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "formula not found"})
		return
	}

	c.JSON(http.StatusOK, formula)
}