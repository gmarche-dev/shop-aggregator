package request

import "github.com/google/uuid"

type StartBill struct {
	StoreID uuid.UUID `json:"store_id" binding:"required"`
}
