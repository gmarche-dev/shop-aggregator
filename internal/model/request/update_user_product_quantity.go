package request

import "github.com/google/uuid"

type UpdateUserProductQuantity struct {
	BillID        uuid.UUID `json:"bill_id"`
	UserProductID uuid.UUID `json:"user_product_id"`
	Quantity      int64     `json:"quantity"`
	ProductType   string    `json:"product_type"`
	ProductSize   string    `json:"product_size"`
	SizeFormat    string    `json:"size_format"`
}
