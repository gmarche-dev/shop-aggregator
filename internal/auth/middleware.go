package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Storer interface {
	EnsureValidToken(context.Context, string) (uuid.UUID, error)
}

func Middleware(a Storer) gin.HandlerFunc {

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token not provided"})
			c.Abort()
			return
		}

		id, err := a.EnsureValidToken(c, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization token"})
			c.Abort()
			return
		}

		// Token is valid, add user ID to the context
		c.Set("userID", id.String())
		c.Next()
	}
}
