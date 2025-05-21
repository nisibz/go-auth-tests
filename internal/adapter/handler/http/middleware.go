package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nisibz/go-auth-tests/internal/core/service"
	"github.com/nisibz/go-auth-tests/internal/core/util"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload_user"
)

func AuthMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is not provided"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authorizationTypeBearer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unsupported authorization type: " + authType})
			return
		}

		accessToken := fields[1]
		userID, err := util.ValidateToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			return
		}

		user, err := userService.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found or unauthorized"})
			return
		}

		user.Password = ""

		c.Set(authorizationPayloadKey, user)
		c.Next()
	}
}
