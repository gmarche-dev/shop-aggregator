package response

import "github.com/google/uuid"

type StartBill struct {
	BillID  uuid.UUID `json:"bill_id"`
	StoreID uuid.UUID `json:"store_id"`
}
