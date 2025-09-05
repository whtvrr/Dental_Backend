package requests

import "github.com/whtvrr/Dental_Backend/internal/domain/entities"

type CreateStatusRequest struct {
	Title       string                `json:"title" binding:"required"`
	Type        entities.StatusType   `json:"type" binding:"required"`
	Code        string                `json:"code" binding:"required"`
	Description *string               `json:"description,omitempty"`
	Color       *string               `json:"color,omitempty"`
}

type UpdateStatusRequest struct {
	Title       string                `json:"title" binding:"required"`
	Type        entities.StatusType   `json:"type" binding:"required"`
	Code        string                `json:"code" binding:"required"`
	Description *string               `json:"description,omitempty"`
	Color       *string               `json:"color,omitempty"`
	IsActive    bool                  `json:"is_active"`
}