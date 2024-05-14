package model

import "github.com/google/uuid"

type Company struct {
	CompanyID   uuid.UUID
	CompanyName string
}
