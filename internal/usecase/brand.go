package usecase

import (
	"context"
	"github.com/rs/zerolog/log"
	"shop-aggregator/internal/model"
)

type BrandStorer interface {
	Insert(ctx context.Context, brand *model.Brand) error
	SelectBrands(ctx context.Context, name string) ([]*model.Brand, error)
	SelectBrandByName(ctx context.Context, name string) (*model.Brand, error)
}

type Brand struct {
	BrandStorer BrandStorer
}

func NewBrand(bs BrandStorer) *Brand {
	return &Brand{
		BrandStorer: bs,
	}
}

func (b *Brand) Create(ctx context.Context, brandName string) (*model.Brand, error) {
	existingBrand, err := b.BrandStorer.SelectBrandByName(ctx, brandName)
	if err != nil {
		log.Error().Caller().Err(err)
		return nil, err
	}
	if existingBrand != nil {
		return nil, model.ErrBrandExists
	}
	brand := model.Brand{
		BrandName: brandName,
	}
	if err = b.BrandStorer.Insert(ctx, &brand); err != nil {
		return nil, err
	}

	return &brand, nil
}

func (b *Brand) SelectByPartialName(ctx context.Context, name string) ([]*model.Brand, error) {
	brands, err := b.BrandStorer.SelectBrands(ctx, name)
	if err != nil {
		log.Error().Caller().Err(err)
		return nil, model.ErrBrandError
	}

	return brands, err
}
