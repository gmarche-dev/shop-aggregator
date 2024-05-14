package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"shop-aggregator/internal/model"
)

type Brand struct {
	db *Client
}

func NewBrand(db *Client) *Brand {
	return &Brand{
		db: db,
	}
}

const (
	InsertBrandQuery = `
		INSERT INTO brand (brand_name)
		VALUES ($1)
		RETURNING brand_id`
	SelectBrandsQuery      = `SELECT brand_id, brand_name FROM brand WHERE brand_name LIKE CONCAT(CAST($1 AS text), '%')`
	SelectBrandByNameQuery = `SELECT brand_id, brand_name FROM brand WHERE brand_name = $1`
)

func (b *Brand) Insert(ctx context.Context, brand *model.Brand) error {
	row := b.db.QueryRow(ctx, InsertBrandQuery, brand.BrandName)
	err := row.Scan(&brand.BrandID)
	return err
}

func (b *Brand) SelectBrands(ctx context.Context, name string) ([]*model.Brand, error) {
	rows, err := b.db.Query(ctx, SelectBrandsQuery, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	brands := []*model.Brand{}
	for rows.Next() {
		brand := &model.Brand{}
		err := rows.Scan(&brand.BrandID, &brand.BrandName)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil
}

func (b *Brand) SelectBrandByName(ctx context.Context, name string) (*model.Brand, error) {
	row := b.db.QueryRow(ctx, SelectBrandByNameQuery, name)
	var brand model.Brand
	if err := row.Scan(&brand.BrandID, &brand.BrandName); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, nil
		}
		return nil, err
	}
	return &brand, nil
}
