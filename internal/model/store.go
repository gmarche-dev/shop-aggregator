package model

import (
	"github.com/google/uuid"
)

const (
	StoreTypeWeb  = "web"
	StoreTypeShop = "shop"
)

type Store struct {
	StoreID   uuid.UUID
	Address   string
	ZipCode   string
	City      string
	Country   string
	StoreName string
	StoreType string
	Url       string
	CompanyID uuid.UUID
}
