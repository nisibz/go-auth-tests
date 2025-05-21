package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nisibz/go-auth-tests/internal/core/domain"
	"github.com/nisibz/go-auth-tests/internal/core/port"
)

type UserHandler struct {
	userService port.UserService
}

func NewUserHandler(us port.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "user not found: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

type ListUsersQuery struct {
	Limit  int64 `form:"limit,default=10"`
	Offset int64 `form:"offset,default=0"`
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var query ListUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		if errMessages := FormatValidationErrors(err); len(errMessages) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMessages})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters: " + err.Error()})
		return
	}

	users, err := h.userService.ListUsers(c.Request.Context(), query.Limit, query.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"email"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userValue, exists := c.Get(authorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	userFromContext, _ := userValue.(*domain.User)
	userID := userFromContext.ID.Hex()

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errMessages := FormatValidationErrors(err); len(errMessages) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMessages})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user: " + err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
