package response

import "net/http"

// StandardResponse represents the standardized API response format
type StandardResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success creates a successful response
func Success(status int, message string, data interface{}) *StandardResponse {
	return &StandardResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// Error creates an error response
func Error(status int, message string) *StandardResponse {
	return &StandardResponse{
		Status:  status,
		Message: message,
	}
}

// OK creates a 200 OK response
func OK(message string, data interface{}) *StandardResponse {
	return Success(http.StatusOK, message, data)
}

// Created creates a 201 Created response
func Created(message string, data interface{}) *StandardResponse {
	return Success(http.StatusCreated, message, data)
}

// BadRequest creates a 400 Bad Request response
func BadRequest(message string) *StandardResponse {
	return Error(http.StatusBadRequest, message)
}

// Unauthorized creates a 401 Unauthorized response
func Unauthorized(message string) *StandardResponse {
	return Error(http.StatusUnauthorized, message)
}

// InternalServerError creates a 500 Internal Server Error response
func InternalServerError(message string) *StandardResponse {
	return Error(http.StatusInternalServerError, message)
}