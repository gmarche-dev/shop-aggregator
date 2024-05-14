package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"shop-aggregator/internal/model/request"
	"shop-aggregator/internal/model/response"
)

type AuthUsecase interface {
	Login(context.Context, string, string) (string, error)
	Logout(context.Context, uuid.UUID) error
}

type Auth struct {
	AuthUsecase AuthUsecase
}

func NewAuth(au AuthUsecase) *Auth {
	return &Auth{
		AuthUsecase: au,
	}
}

func (a *Auth) Login(c *gin.Context) {
	var login request.Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.AuthUsecase.Login(c.Request.Context(), login.Login, login.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Login{Token: token})
}

func (a *Auth) Logout(c *gin.Context) {
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	err := a.AuthUsecase.Logout(c.Request.Context(), uuid.MustParse(id.(string)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
