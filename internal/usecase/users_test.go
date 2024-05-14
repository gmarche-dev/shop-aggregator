package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/usecase"
	"shop-aggregator/internal/utils"
	"testing"
)

func TestUsers_CreateOrUpdateUser(t *testing.T) {
	ctx := context.Background()
	mockUserStorer := NewUsersStorer(t)

	expectedUser := &model.User{
		Login:    "login",
		Password: "password",
		Email:    "email",
	}

	expectedError := errors.New("random error")

	u := usecase.NewUsers(mockUserStorer)

	t.Run("GetUserByEmail error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByEmail(ctx, expectedUser.Email).Return(nil, expectedError).Once()
		assert.ErrorIs(t, model.ErrUserError, u.CreateOrUpdateUser(ctx, expectedUser))
	})

	t.Run("user exists", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByEmail(ctx, expectedUser.Email).Return(&model.User{}, nil).Once()
		assert.Equal(t, fmt.Errorf("user exist for email %s", expectedUser.Email), u.CreateOrUpdateUser(ctx, expectedUser))
	})

	t.Run("GetUserByLogin error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByEmail(ctx, expectedUser.Email).Return(nil, nil).Once()
		mockUserStorer.EXPECT().GetUserByLogin(ctx, expectedUser.Login).Return(nil, expectedError).Once()
		assert.ErrorIs(t, model.ErrUserError, u.CreateOrUpdateUser(ctx, expectedUser))
	})

	t.Run("GetUserByLogin user exists error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByEmail(ctx, expectedUser.Email).Return(nil, nil).Once()
		mockUserStorer.EXPECT().GetUserByLogin(ctx, expectedUser.Login).Return(&model.User{}, nil).Once()
		assert.Equal(t, fmt.Errorf("user exist for login %s", expectedUser.Login), u.CreateOrUpdateUser(ctx, expectedUser))
	})

	t.Run("GetUserByLogin user exists error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByEmail(ctx, expectedUser.Email).Return(nil, nil).Once()
		mockUserStorer.EXPECT().GetUserByLogin(ctx, expectedUser.Login).Return(&model.User{}, nil).Once()
		assert.Equal(t, fmt.Errorf("user exist for login %s", expectedUser.Login), u.CreateOrUpdateUser(ctx, expectedUser))
	})

	t.Run("Upsert error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByEmail(ctx, expectedUser.Email).Return(nil, nil).Once()
		mockUserStorer.EXPECT().GetUserByLogin(ctx, expectedUser.Login).Return(nil, nil).Once()
		mockUserStorer.EXPECT().Upsert(ctx, mock.Anything).Return(expectedError).Once()
		assert.ErrorIs(t, model.ErrUserError, u.CreateOrUpdateUser(ctx, expectedUser))
	})

	t.Run("no error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByEmail(ctx, expectedUser.Email).Return(nil, nil).Once()
		mockUserStorer.EXPECT().GetUserByLogin(ctx, expectedUser.Login).Return(nil, nil).Once()
		mockUserStorer.EXPECT().Upsert(ctx, mock.Anything).Run(func(_a0 context.Context, _a1 *model.User) {
			require.True(t, utils.CheckPasswordHash(expectedUser.Password, _a1.HashPassword))
		}).Return(nil).Once()
		assert.NoError(t, u.CreateOrUpdateUser(ctx, expectedUser))
	})
}

func TestUsers_UpdatePassword(t *testing.T) {
	ctx := context.Background()
	mockUserStorer := NewUsersStorer(t)

	expectedUpdatePassword := &model.UpdatePassword{
		ID:          uuid.New(),
		OldPassword: "password",
		Password:    "new-password",
	}

	expectedOldHash, err := utils.HashPassword("password")
	require.NoError(t, err)

	expectedUser := &model.User{
		HashPassword: expectedOldHash,
	}

	expectedError := errors.New("random error")

	u := usecase.NewUsers(mockUserStorer)

	t.Run("GetUserByID error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUpdatePassword.ID).Return(nil, expectedError).Once()
		assert.ErrorIs(t, model.ErrUserError, u.UpdatePassword(ctx, expectedUpdatePassword))
	})

	t.Run("user not found error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUpdatePassword.ID).Return(nil, nil).Once()
		assert.ErrorIs(t, model.ErrUserNotFound, u.UpdatePassword(ctx, expectedUpdatePassword))
	})

	t.Run("CheckPasswordHash invalid old password", func(t *testing.T) {
		eup := expectedUpdatePassword
		eup.OldPassword = ""
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUpdatePassword.ID).Return(expectedUser, nil).Once()
		assert.ErrorIs(t, model.ErrOldPasswordError, u.UpdatePassword(ctx, expectedUpdatePassword))
	})

	t.Run("UpdatePassword error", func(t *testing.T) {
		eup := expectedUpdatePassword
		eup.OldPassword = ""
		mockUserStorer.EXPECT().GetUserByID(ctx, eup.ID).Return(expectedUser, nil).Once()
		assert.ErrorIs(t, model.ErrOldPasswordError, u.UpdatePassword(ctx, eup))
	})

	t.Run("UpdatePassword error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUpdatePassword.ID).Return(expectedUser, nil).Once()
		mockUserStorer.EXPECT().UpdatePassword(ctx, expectedUpdatePassword.ID, mock.Anything).Return(expectedError).Once()
		assert.ErrorIs(t, model.ErrUserError, u.UpdatePassword(ctx, expectedUpdatePassword))
	})

	t.Run("no error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUpdatePassword.ID).Return(expectedUser, nil).Once()
		mockUserStorer.EXPECT().UpdatePassword(ctx, expectedUpdatePassword.ID, mock.Anything).Run(func(_a0 context.Context, _a1 uuid.UUID, _a2 string) {
			require.True(t, utils.CheckPasswordHash(expectedUpdatePassword.Password, _a2))
		}).Return(nil).Once()
		assert.Nil(t, u.UpdatePassword(ctx, expectedUpdatePassword))
	})
}

func TestUsers_UpdateEmail(t *testing.T) {
	ctx := context.Background()
	mockUserStorer := NewUsersStorer(t)

	expectedUpdateEmail := &model.UpdateEmail{
		ID:    uuid.New(),
		Email: "new@email.com",
	}

	expectedUser := &model.User{
		Email: "test@test.com",
	}

	expectedError := errors.New("random error")

	u := usecase.NewUsers(mockUserStorer)

	t.Run("GetUserByID error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUpdateEmail.ID).Return(nil, expectedError).Once()
		assert.ErrorIs(t, model.ErrUserError, u.UpdateEmail(ctx, expectedUpdateEmail))
	})

	t.Run("GetUserByID, user not found", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUpdateEmail.ID).Return(nil, nil).Once()
		assert.ErrorIs(t, model.ErrUserNotFound, u.UpdateEmail(ctx, expectedUpdateEmail))
	})

	t.Run("UpdateEmail, error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUpdateEmail.ID).Return(expectedUser, nil).Once()
		mockUserStorer.EXPECT().UpdateEmail(ctx, expectedUpdateEmail.ID, expectedUpdateEmail.Email).Return(expectedError).Once()
		assert.ErrorIs(t, model.ErrUserError, u.UpdateEmail(ctx, expectedUpdateEmail))
	})

	t.Run("no error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUpdateEmail.ID).Return(expectedUser, nil).Once()
		mockUserStorer.EXPECT().UpdateEmail(ctx, expectedUpdateEmail.ID, expectedUpdateEmail.Email).Return(nil).Once()
		assert.NoError(t, u.UpdateEmail(ctx, expectedUpdateEmail))
	})
}

func TestUsers_GetUserByID(t *testing.T) {
	ctx := context.Background()
	mockUserStorer := NewUsersStorer(t)

	expectedUser := &model.User{
		ID:    uuid.New(),
		Email: "test@test.com",
	}

	expectedError := errors.New("random error")

	u := usecase.NewUsers(mockUserStorer)

	t.Run("GetUserByID error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUser.ID).Return(nil, expectedError).Once()
		result, err := u.GetUserByID(ctx, expectedUser.ID)
		assert.ErrorIs(t, model.ErrUserError, err)
		assert.Nil(t, result)
	})

	t.Run("GetUserByID, user not found", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUser.ID).Return(nil, nil).Once()
		result, err := u.GetUserByID(ctx, expectedUser.ID)
		assert.ErrorIs(t, model.ErrUserNotFound, err)
		assert.Nil(t, result)
	})

	t.Run("no error", func(t *testing.T) {
		mockUserStorer.EXPECT().GetUserByID(ctx, expectedUser.ID).Return(expectedUser, nil).Once()
		result, err := u.GetUserByID(ctx, expectedUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})
}
