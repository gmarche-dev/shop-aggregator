package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/utils"
)

type UsersStorer interface {
	Upsert(context.Context, *model.User) error
	GetUserByEmail(context.Context, string) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdatePassword(context.Context, uuid.UUID, string) error
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	UpdateEmail(ctx context.Context, id uuid.UUID, email string) error
}

type Users struct {
	UsersStorer UsersStorer
}

func NewUsers(us UsersStorer) *Users {
	return &Users{
		UsersStorer: us,
	}
}

func (u *Users) CreateOrUpdateUser(ctx context.Context, m *model.User) error {
	checkUserEmailExist, err := u.UsersStorer.GetUserByEmail(ctx, m.Email)
	if err != nil {
		log.Error().Caller().Err(err)
		return model.ErrUserError
	}

	if checkUserEmailExist != nil {
		return fmt.Errorf("user exist for email %s", m.Email)
	}

	checkUserLoginExist, err := u.UsersStorer.GetUserByLogin(ctx, m.Login)
	if err != nil {
		log.Error().Caller().Err(err)
		return model.ErrUserError
	}

	if checkUserLoginExist != nil {
		return fmt.Errorf("user exist for login %s", m.Login)
	}

	m.HashPassword, err = utils.HashPassword(m.Password)
	if err != nil {
		log.Error().Caller().Err(err)
		return err
	}

	if err = u.UsersStorer.Upsert(ctx, m); err != nil {
		log.Error().Caller().Err(err)
		return model.ErrUserError
	}

	return err
}

func (u *Users) UpdatePassword(ctx context.Context, up *model.UpdatePassword) error {
	user, err := u.UsersStorer.GetUserByID(ctx, up.ID)
	if err != nil {
		log.Error().Caller().Err(err)
		return model.ErrUserError
	}

	if user == nil {
		return model.ErrUserNotFound
	}

	if !utils.CheckPasswordHash(up.OldPassword, user.HashPassword) {
		return model.ErrOldPasswordError
	}

	hashedPassword, err := utils.HashPassword(up.Password)
	if err != nil {
		log.Error().Caller().Err(err)
		return model.ErrUserError
	}

	err = u.UsersStorer.UpdatePassword(ctx, up.ID, hashedPassword)
	if err != nil {
		log.Error().Caller().Err(err)
		return model.ErrUserError
	}

	return nil
}

func (u *Users) UpdateEmail(ctx context.Context, ue *model.UpdateEmail) error {
	user, err := u.UsersStorer.GetUserByID(ctx, ue.ID)
	if err != nil {
		log.Error().Caller().Err(err)
		return model.ErrUserError
	}

	if user == nil {
		return model.ErrUserNotFound
	}

	err = u.UsersStorer.UpdateEmail(ctx, ue.ID, ue.Email)
	if err != nil {
		log.Error().Caller().Err(err)
		return model.ErrUserError
	}

	return nil
}

func (u *Users) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := u.UsersStorer.GetUserByID(ctx, id)
	if err != nil {
		log.Error().Caller().Err(err)
		return nil, model.ErrUserError
	}
	if user == nil {
		return nil, model.ErrUserNotFound
	}
	return user, nil
}
