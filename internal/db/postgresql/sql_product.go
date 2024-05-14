package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"shop-aggregator/internal/model"
)

type Product struct {
	db *Client
}

func NewProduct(db *Client) *Product {
	return &Product{
		db: db,
	}
}

const (
	InsertProductQuery = `
		INSERT INTO product (ean, product_name, brand_id)
		VALUES ($1, $2, $3)
		RETURNING product_id`
	GetProductByEANQuery = `
		SELECT p.product_id, p.ean, p.product_name, p.brand_id
		FROM product p 
		WHERE p.ean = $1`
)

func (p *Product) Insert(ctx context.Context, product *model.Product) error {
	row := p.db.QueryRow(ctx, InsertProductQuery, product.EAN, product.ProductName, product.BrandID)
	err := row.Scan(&product.ProductID)
	return err
}

func (p *Product) GetProductByEAN(ctx context.Context, ean string) (*model.Product, error) {
	row := p.db.QueryRow(ctx, GetProductByEANQuery, ean)
	product := &model.Product{}
	err := row.Scan(&product.ProductID, &product.EAN, &product.ProductName, &product.BrandID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return product, nil
}
