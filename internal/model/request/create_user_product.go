package request

import (
	"github.com/google/uuid"
)

type CreateUserProduct struct {
	ProductID   uuid.UUID `json:"product_id" binding:"required"`
	BillID      uuid.UUID `json:"bill_id" binding:"required"`
	ProductType string    `json:"product_type" binding:"required"`
	ProductSize string    `json:"product_size"`
	SizeFormat  string    `json:"size_format"`
	Price       string    `json:"price" binding:"required"`
	Quantity    int64     `json:"quantity" binding:"required"`
}
