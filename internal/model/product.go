package model

import (
	"github.com/google/uuid"
)

var (
	BulkProductIDMeat      = uuid.MustParse("2e30955b-0f88-43df-8924-1ec21afed0aa")
	BulkProductIDVegetable = uuid.MustParse("eb5be0d0-b3f6-4f2f-b582-9a7dd566b549")
	BulkProductIDFruits    = uuid.MustParse("3c94d3d7-7bce-40d6-8f11-c8fdeb328d41")
)

type Product struct {
	ProductID   uuid.UUID
	EAN         string
	ProductName string
	ProductType string
	BrandID     uuid.UUID
}
