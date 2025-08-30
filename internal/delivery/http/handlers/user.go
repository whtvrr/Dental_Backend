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

type UserHandler struct {
	userUseCase *usecases.UserUseCase
}

func NewUserHandler(userUseCase *usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body entities.User true "User data"
// @Success 201 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
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
// @Description Update a user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body entities.User true "Updated user data"
// @Success 200 {object} response.StandardResponse
// @Failure 400 {object} response.StandardResponse "Bad Request"
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("invalid user id"))
		return
	}

	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}
	user.ID = id

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
// @Summary Get all doctors
// @Description Get a list of all users with doctor role
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.StandardResponse
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users/doctors [get]
func (h *UserHandler) GetDoctors(c *gin.Context) {
	doctors, err := h.userUseCase.GetDoctors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Doctors retrieved successfully", gin.H{"doctors": doctors}))
}

// GetClients godoc
// @Summary Get all clients
// @Description Get a list of all users with client role
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.StandardResponse
// @Failure 500 {object} response.StandardResponse "Internal Server Error"
// @Router /users/clients [get]
func (h *UserHandler) GetClients(c *gin.Context) {
	clients, err := h.userUseCase.GetClients(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Clients retrieved successfully", gin.H{"clients": clients}))
}
