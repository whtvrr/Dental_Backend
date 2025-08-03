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