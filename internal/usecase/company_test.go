package usecase_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/usecase"
	"testing"
)

func TestCompany_GetCompanyByPartialName(t *testing.T) {
	ctx := context.Background()

	expectedError := errors.New("random error")
	mockCompanyStorer := NewCompanyStorer(t)
	expectedResponse := []*model.Company{
		{},
		{},
		{},
		{},
	}
	c := usecase.NewCompany(mockCompanyStorer)

	t.Run("SelectCompanies error", func(t *testing.T) {
		mockCompanyStorer.EXPECT().SelectCompanies(ctx, "coin").Return(nil, expectedError).Once()
		cs, err := c.SelectByPartialName(ctx, "coin")
		assert.ErrorIs(t, model.ErrSelectCompaniesError, err)
		assert.Nil(t, cs)
	})

	t.Run("no error", func(t *testing.T) {
		mockCompanyStorer.EXPECT().SelectCompanies(ctx, "coin").Return(expectedResponse, nil).Once()
		cs, err := c.SelectByPartialName(ctx, "coin")
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, cs)
	})
}
