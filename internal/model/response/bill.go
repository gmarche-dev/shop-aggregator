package response

import (
	"github.com/google/uuid"
	"shop-aggregator/internal/model"
)

type Bill struct {
	BillID   uuid.UUID      `json:"bill_id"`
	Amount   string         `json:"amount"`
	State    string         `json:"state"`
	Store    *BillStore     `json:"store"`
	Products []*UserProduct `json:"products"`
}

type BillStore struct {
	*Store
	CompanyName string `json:"company_name"`
}

func NewBillStore(s *model.Store, c *model.Company) *BillStore {
	if s == nil {
		return &BillStore{}
	}
	bs := &BillStore{
		Store: NewStoreFromModel(s),
	}
	if c != nil {
		bs.CompanyName = c.CompanyName
	}

	return bs
}

func NewBillFromModel(m *model.Bill, s *model.Store, c *model.Company, ps []*model.UserProduct) *Bill {
	b := &Bill{
		BillID:   m.BillID,
		Amount:   m.Amount,
		State:    m.State,
		Store:    NewBillStore(s, c),
		Products: NewUserProductsFromModel(ps),
	}

	return b
}

func NewBillsFromModels(m []*model.Bill) []*Bill {
	var b []*Bill
	for _, n := range m {
		b = append(b, NewBillFromModel(n, nil, nil, nil))
	}

	return b
}
