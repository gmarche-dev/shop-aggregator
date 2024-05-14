package usecase

import (
	"context"
	"github.com/rs/zerolog/log"
	"shop-aggregator/internal/model"
)

type CompanyStorer interface {
	SelectCompanies(ctx context.Context, name string) ([]*model.Company, error)
}

type Company struct {
	CompanyStorer CompanyStorer
}

func NewCompany(cs CompanyStorer) *Company {
	return &Company{
		CompanyStorer: cs,
	}
}

func (c *Company) SelectByPartialName(ctx context.Context, partialName string) ([]*model.Company, error) {
	cs, err := c.CompanyStorer.SelectCompanies(ctx, partialName)
	if err != nil {
		log.Error().Caller().Err(err)
		return nil, model.ErrSelectCompaniesError
	}

	return cs, nil
}
