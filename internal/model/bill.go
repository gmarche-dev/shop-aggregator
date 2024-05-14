package model

import "github.com/google/uuid"

const (
	BillStateCreate    = "create"
	BillStateCompleted = "complete"
	BillStateCanceled  = "cancel"
)

type Bill struct {
	BillID  uuid.UUID
	UserID  uuid.UUID
	StoreID uuid.UUID
	Amount  string
	State   string
}
