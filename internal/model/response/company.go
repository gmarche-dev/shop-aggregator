package response

import (
	"github.com/google/uuid"
	"shop-aggregator/internal/model"
)

type Company struct {
	CompanyID   uuid.UUID `json:"company_id"`
	CompanyName string    `json:"company_name"`
}

func NewCompanyFromModel(m *model.Company) *Company {
	if m == nil {
		return &Company{}
	}
	return &Company{
		CompanyID:   m.CompanyID,
		CompanyName: m.CompanyName,
	}
}

func NewCompaniesFromModels(ms []*model.Company) []*Company {
	brs := []*Company{}
	for _, m := range ms {
		brs = append(brs, NewCompanyFromModel(m))
	}
	return brs
}
