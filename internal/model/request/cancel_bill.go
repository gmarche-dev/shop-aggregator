package request

import "github.com/google/uuid"

type CancelBill struct {
	BillID uuid.UUID `json:"bill_id" binding:"required"`
}
