package usecase_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/usecase"
	"shop-aggregator/internal/utils"
	"testing"
)

func TestAuth_Login(t *testing.T) {
	ctx := context.Background()
	mockUserStore := NewAuthUserStorer(t)
	mockAuthStore := NewAuthStorer(t)

	au := usecase.NewAuth(mockAuthStore, mockUserStore)
	expectedLogin := "login"
	expectedPassword := "password"
	expectedHashPassword, errHash := utils.HashPassword(expectedPassword)
	require.NoError(t, errHash)
	expectedUser := model.User{
		HashPassword: expectedHashPassword,
	}

	expectedError := errors.New("random error")

	t.Run("GetUserByLogin error", func(t *testing.T) {
		mockUserStore.EXPECT().GetUserByLogin(ctx, expectedLogin).Return(nil, expectedError).Once()
		_, err := au.Login(ctx, expectedLogin, expectedPassword)
		require.ErrorIs(t, model.ErrUserError, err)
	})

	t.Run("GetUserByLogin, user not found", func(t *testing.T) {
		mockUserStore.EXPECT().GetUserByLogin(ctx, expectedLogin).Return(nil, nil).Once()
		_, err := au.Login(ctx, expectedLogin, expectedPassword)
		require.ErrorIs(t, model.ErrUserNotFound, err)
	})

	t.Run("not same old password", func(t *testing.T) {
		mockUserStore.EXPECT().GetUserByLogin(ctx, expectedLogin).Return(&expectedUser, nil).Once()
		_, err := au.Login(ctx, expectedLogin, "")
		require.ErrorIs(t, model.ErrPasswordError, err)
	})

	t.Run("Upsert error", func(t *testing.T) {
		mockUserStore.EXPECT().GetUserByLogin(ctx, expectedLogin).Return(&expectedUser, nil).Once()
		mockAuthStore.EXPECT().Upsert(ctx, expectedUser.ID, mock.Anything).Run(func(_a0 context.Context, _a1 uuid.UUID, _a2 string) {
			require.NotEmpty(t, _a2)
		}).Return(expectedError).Once()
		_, err := au.Login(ctx, expectedLogin, expectedPassword)
		require.ErrorIs(t, model.ErrUserError, err)
	})

	t.Run("no error", func(t *testing.T) {
		mockUserStore.EXPECT().GetUserByLogin(ctx, expectedLogin).Return(&expectedUser, nil).Once()
		mockAuthStore.EXPECT().Upsert(ctx, expectedUser.ID, mock.Anything).Run(func(_a0 context.Context, _a1 uuid.UUID, _a2 string) {
			require.NotEmpty(t, _a2)
		}).Return(nil).Once()
		token, err := au.Login(ctx, expectedLogin, expectedPassword)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})
}

func TestAuth_Logout(t *testing.T) {
	ctx := context.Background()
	mockAuthStore := NewAuthStorer(t)
	au := usecase.NewAuth(mockAuthStore, nil)
	expectedID := uuid.New()
	expectedError := errors.New("random error")

	t.Run("Logout error", func(t *testing.T) {
		mockAuthStore.EXPECT().Logout(ctx, expectedID).Return(expectedError).Once()
		require.ErrorIs(t, model.ErrUserNotFound, au.Logout(ctx, expectedID))
	})

	t.Run("no error", func(t *testing.T) {
		mockAuthStore.EXPECT().Logout(ctx, expectedID).Return(nil).Once()
		require.Nil(t, au.Logout(ctx, expectedID))
	})
}
