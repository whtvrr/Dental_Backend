package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/internal/delivery/http/response"
	"github.com/whtvrr/Dental_Backend/internal/domain/entities"
	"github.com/whtvrr/Dental_Backend/internal/usecases"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userUseCase *usecases.UserUseCase
}

type CreateUserRequest struct {
	FullName    string  `json:"full_name" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Address     *string `json:"address,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"` // Accepts DD.MM.YYYY format
}

type UpdateUserRequest struct {
	Email       *string            `json:"email,omitempty"`
	Role        *entities.UserRole `json:"role,omitempty"`
	FullName    *string            `json:"full_name,omitempty"`
	PhoneNumber *string            `json:"phone_number,omitempty"`
	Address     *string            `json:"address,omitempty"`
	Gender      *string            `json:"gender,omitempty"`
	BirthDate   *string            `json:"birth_date,omitempty"` // Accepts DD.MM.YYYY format
}

func NewUserHandler(userUseCase *usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUser godoc
// @Summary Create a new client user
// @Description Create a new client user with the provided details. Users created via this endpoint are automatically set as clients. BirthDate should be in DD.MM.YYYY format (e.g., 24.06.2003). CreatedAt and UpdatedAt are automatically set.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	// Convert request to user entity - automatically set as client
	user := entities.User{
		Role:        entities.RoleClient,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Gender:      req.Gender,
	}

	// Parse birth date if provided
	if req.BirthDate != nil && strings.TrimSpace(*req.BirthDate) != "" {
		birthDate, err := time.Parse("02.01.2006", strings.TrimSpace(*req.BirthDate))
		if err != nil {
			c.JSON(http.StatusBadRequest, response.BadRequest("invalid birth date format, expected DD.MM.YYYY (e.g., 24.06.2003)"))
			return
		}
		user.BirthDate = &birthDate
	}

	if err := h.userUseCase.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Created("User created successfully", user))
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 404 {object} response.StandardResponse "Not Found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid user id"))
		return
	}

	user, err := h.userUseCase.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "user not found"))
		return
	}

	c.JSON(http.StatusOK, response.OK("User retrieved successfully", user))
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update a user with the provided details. BirthDate should be in DD.MM.YYYY format (e.g., 24.06.2003). UpdatedAt is automatically set.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body UpdateUserRequest true "Updated user data"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 404 {object} response.StandardResponse "Not Found"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid user id"))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	// Get existing user to preserve CreatedAt and other fields
	existingUser, err := h.userUseCase.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "user not found"))
		return
	}

	// Update only provided fields
	user := *existingUser
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.Address != nil {
		user.Address = req.Address
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}

	// Parse birth date if provided
	if req.BirthDate != nil && strings.TrimSpace(*req.BirthDate) != "" {
		birthDate, err := time.Parse("02.01.2006", strings.TrimSpace(*req.BirthDate))
		if err != nil {
			c.JSON(http.StatusBadRequest, response.BadRequest("invalid birth date format, expected DD.MM.YYYY (e.g., 24.06.2003)"))
			return
		}
		user.BirthDate = &birthDate
	}

	if err := h.userUseCase.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("User updated successfully", user))
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid user id"))
		return
	}

	if err := h.userUseCase.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("User deleted successfully", nil))
}

// ListUsers godoc
// @Summary List all users with pagination
// @Description Get a paginated list of all users
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
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

	users, err := h.userUseCase.ListUsers(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Users retrieved successfully", gin.H{"users": users}))
}

// GetDoctors godoc
// @Summary Get all doctors with pagination
// @Description Get a paginated list of all users with doctor role
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users/doctors [get]
func (h *UserHandler) GetDoctors(c *gin.Context) {
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

	doctors, err := h.userUseCase.GetDoctorsWithPagination(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Doctors retrieved successfully", gin.H{"doctors": doctors}))
}

// GetClients godoc
// @Summary Get all clients with pagination
// @Description Get a paginated list of all users with client role
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users/clients [get]
func (h *UserHandler) GetClients(c *gin.Context) {
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

	clients, err := h.userUseCase.GetClientsWithPagination(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Clients retrieved successfully", gin.H{"clients": clients}))
}

// GetStaff godoc
// @Summary Get all staff members with pagination
// @Description Get a paginated list of all users with doctor and receptionist roles
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users/staff [get]
func (h *UserHandler) GetStaff(c *gin.Context) {
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

	staff, err := h.userUseCase.GetStaffWithPagination(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Staff retrieved successfully", gin.H{"staff": staff}))
}
