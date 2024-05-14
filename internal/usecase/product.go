package usecase

import (
	"context"
	"github.com/rs/zerolog/log"
	"shop-aggregator/internal/model"
)

type ProductBrandStorer interface {
	Insert(ctx context.Context, company *model.Brand) error
	SelectBrandByName(ctx context.Context, name string) (*model.Brand, error)
}

type ProductStorer interface {
	Insert(ctx context.Context, product *model.Product) error
	GetProductByEAN(ctx context.Context, ean string) (*model.Product, error)
}

type Product struct {
	ProductStorer      ProductStorer
	ProductBrandStorer ProductBrandStorer
}

func NewProduct(ps ProductStorer, pbs ProductBrandStorer) *Product {
	return &Product{
		ProductStorer:      ps,
		ProductBrandStorer: pbs,
	}
}

func (p *Product) createOrGetBrand(ctx context.Context, brandName string) (*model.Brand, error) {
	exists, err := p.ProductBrandStorer.SelectBrandByName(ctx, brandName)
	if err != nil {
		log.Error().Caller().Err(err)
		return nil, model.ErrBrandError
	}
	if exists != nil {
		return exists, nil
	}
	brand := model.Brand{
		BrandName: brandName,
	}
	if err := p.ProductBrandStorer.Insert(ctx, &brand); err != nil {
		log.Error().Caller().Err(err)
		return nil, model.ErrBrandError
	}

	return &brand, nil
}

func (p *Product) Create(ctx context.Context, m *model.Product, brandName string) (*model.Product, error) {
	brand, err := p.createOrGetBrand(ctx, brandName)
	if err != nil {
		return nil, err
	}
	m.BrandID = brand.BrandID
	exists, err := p.ProductStorer.GetProductByEAN(ctx, m.EAN)
	if err != nil {
		log.Error().Caller().Err(err)
		return nil, model.ErrProductError
	}
	if exists != nil {
		return exists, nil
	}
	if err = p.ProductStorer.Insert(ctx, m); err != nil {
		log.Error().Caller().Err(err)
		return nil, model.ErrProductError
	}

	return m, nil
}

func (p *Product) GetProductByEAN(ctx context.Context, ean string) (*model.Product, error) {
	product, err := p.ProductStorer.GetProductByEAN(ctx, ean)
	if err != nil {
		log.Error().Caller().Err(err)
		return nil, model.ErrProductError
	}
	if product == nil {
		return nil, model.ErrNotExistsError
	}

	return product, nil
}
