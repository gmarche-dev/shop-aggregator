package response

import (
	"github.com/google/uuid"
	"shop-aggregator/internal/model"
)

type UserProduct struct {
	UserProductID uuid.UUID `json:"user_product_id"`
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Ean           string    `json:"ean"`
	BrandID       uuid.UUID `json:"brand_id"`
	BrandName     string    `json:"brand_name"`
	BillID        uuid.UUID `json:"bill_id"`
	Price         string    `json:"price"`
	Quantity      int64     `json:"quantity"`
	ProductType   string    `json:"product_type"`
	ProductSize   string    `json:"product_size"`
	SizeFormat    string    `json:"size_format"`
}

func NewUserProductFromModel(m *model.UserProduct) *UserProduct {
	up := &UserProduct{
		UserProductID: m.UserProductID,
		ProductID:     m.ProductID,
		ProductName:   m.ProductName,
		Ean:           m.Ean,
		BrandID:       m.BrandID,
		BrandName:     m.BrandName,
		BillID:        m.BillID,
		Price:         m.Price,
		Quantity:      m.Quantity,
		ProductType:   m.ProductType,
		ProductSize:   m.ProductSize,
		SizeFormat:    m.SizeFormat,
	}

	return up
}

func NewUserProductsFromModel(mups []*model.UserProduct) []*UserProduct {
	up := []*UserProduct{}
	if mups == nil {
		return up
	}
	for _, mup := range mups {
		up = append(up, NewUserProductFromModel(mup))
	}
	return up
}
