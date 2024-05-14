package response

import (
	"github.com/google/uuid"
	"shop-aggregator/internal/model"
)

type Store struct {
	StoreID   uuid.UUID `json:"store_id"`
	Address   string    `json:"address"`
	ZipCode   string    `json:"zip_code"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	StoreName string    `json:"store_name"`
	StoreType string    `json:"store_type"`
	Url       string    `json:"url"`
	CompanyID uuid.UUID `json:"company_id"`
}

func NewStoreFromModel(m *model.Store) *Store {
	if m == nil {
		return &Store{}
	}

	return &Store{
		StoreID:   m.StoreID,
		Address:   m.Address,
		ZipCode:   m.ZipCode,
		City:      m.City,
		Country:   m.Country,
		StoreName: m.StoreName,
		StoreType: m.StoreType,
		Url:       m.Url,
		CompanyID: m.CompanyID,
	}
}

func NewStoresFromModels(mss []*model.Store) []*Store {
	r := []*Store{}
	for _, ms := range mss {
		r = append(r, NewStoreFromModel(ms))
	}

	return r
}
