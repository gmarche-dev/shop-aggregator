package response

import (
	"github.com/google/uuid"
	"shop-aggregator/internal/model"
)

type Product struct {
	ProductID   uuid.UUID `json:"product_id"`
	EAN         string    `json:"ean"`
	ProductName string    `json:"product_name"`
	BrandID     uuid.UUID `json:"brand_id"`
}

func NewProductFromModel(m *model.Product) *Product {
	return &Product{
		ProductID:   m.ProductID,
		EAN:         m.EAN,
		ProductName: m.ProductName,
		BrandID:     m.BrandID,
	}
}
