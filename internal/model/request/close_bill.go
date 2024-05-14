package request

import "github.com/google/uuid"

type CloseBill struct {
	BillID uuid.UUID `json:"bill_id" binding:"required"`
	Amount string    `json:"amount" binding:"required"`
}
