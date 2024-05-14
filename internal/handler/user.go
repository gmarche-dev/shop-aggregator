package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/model/request"
	"shop-aggregator/internal/model/response"
)

type UserUsecase interface {
	CreateOrUpdateUser(ctx context.Context, m *model.User) error
	UpdatePassword(ctx context.Context, up *model.UpdatePassword) error
	UpdateEmail(ctx context.Context, ue *model.UpdateEmail) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type User struct {
	UserUsecase UserUsecase
}

func NewUser(uu UserUsecase) *User {
	return &User{
		UserUsecase: uu,
	}
}

func (u *User) CreateUser(c *gin.Context) {
	var cu request.CreateUser
	if err := c.ShouldBindJSON(&cu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := model.User{
		Login:    cu.Login,
		Email:    cu.Email,
		Password: cu.Password,
	}

	if err := u.UserUsecase.CreateOrUpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func (u *User) UpdatePassword(c *gin.Context) {
	var up request.UpdatePassword
	if err := c.ShouldBindJSON(&up); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	updatePassword := model.UpdatePassword{
		ID:          uuid.MustParse(id.(string)),
		Password:    up.Password,
		OldPassword: up.OldPassword,
	}

	if err := u.UserUsecase.UpdatePassword(c.Request.Context(), &updatePassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func (u *User) UpdateEmail(c *gin.Context) {
	var ue request.UpdateEmail
	if err := c.ShouldBindJSON(&ue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	updateEmail := model.UpdateEmail{
		ID:    uuid.MustParse(id.(string)),
		Email: ue.Email,
	}

	if err := u.UserUsecase.UpdateEmail(c.Request.Context(), &updateEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email updated"})
}

func (u *User) GetUser(c *gin.Context) {
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	user, err := u.UserUsecase.GetUserByID(c.Request.Context(), uuid.MustParse(id.(string)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewUserFromModel(*user))
}
