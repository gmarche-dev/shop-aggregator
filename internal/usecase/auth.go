package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/utils"
)

type AuthStorer interface {
	Upsert(context.Context, uuid.UUID, string) error
	Logout(context.Context, uuid.UUID) error
}

type AuthUserStorer interface {
	GetUserByLogin(context.Context, string) (*model.User, error)
}

type Auth struct {
	AuthStorer     AuthStorer
	AuthUserStorer AuthUserStorer
}

func NewAuth(a AuthStorer, au AuthUserStorer) *Auth {
	return &Auth{
		AuthStorer:     a,
		AuthUserStorer: au,
	}
}

func (a *Auth) Login(ctx context.Context, login, password string) (string, error) {
	user, err := a.AuthUserStorer.GetUserByLogin(ctx, login)
	if err != nil {
		log.Error().Caller().Err(err)
		return "", model.ErrUserError
	}

	if user == nil {
		return "", model.ErrUserNotFound
	}

	if !utils.CheckPasswordHash(password, user.HashPassword) {
		return "", model.ErrPasswordError
	}

	token, err := utils.GenerateToken(128)
	if err != nil {
		log.Error().Caller().Err(err)
		return "", model.ErrUserError
	}

	err = a.AuthStorer.Upsert(ctx, user.ID, token)
	if err != nil {
		log.Error().Caller().Err(err)
		return "", model.ErrUserError
	}

	return token, nil
}

func (a *Auth) Logout(ctx context.Context, id uuid.UUID) error {
	err := a.AuthStorer.Logout(ctx, id)
	if err != nil {
		log.Error().Caller().Err(err)
		return model.ErrUserNotFound
	}
	return nil
}
