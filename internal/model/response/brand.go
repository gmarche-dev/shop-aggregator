package response

import (
	"github.com/google/uuid"
	"shop-aggregator/internal/model"
)

type Brand struct {
	BrandID   uuid.UUID `json:"brand_id"`
	BrandName string    `json:"brand_name"`
}

func NewBrandFromModel(m *model.Brand) *Brand {
	return &Brand{
		BrandID:   m.BrandID,
		BrandName: m.BrandName,
	}
}

func NewBrandsFromModels(ms []*model.Brand) []*Brand {
	var brs []*Brand
	for _, m := range ms {
		brs = append(brs, NewBrandFromModel(m))
	}
	return brs
}
