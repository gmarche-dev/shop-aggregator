package model

import (
	"github.com/google/uuid"
)

const (
	ProductBulk     = "bulk_product"
	ProductBarcoded = "barcoded_product"
)

const (
	SizeFormatSizeTypeWeight = "weight"
	SizeFormatVolume         = "volume"
)

const (
	SizeFormatWeightGr = "gr"
	SizeFormatWeightKg = "kg"
)

const (
	SizeFormatVolumeMl = "ml"
	SizeFormatVolumeL  = "l"
)

type UserProduct struct {
	UserProductID uuid.UUID
	ProductID     uuid.UUID
	UserID        uuid.UUID
	ProductName   string
	Ean           string
	BrandID       uuid.UUID
	BrandName     string
	StoreID       uuid.UUID
	StoreName     string
	BillID        uuid.UUID
	Price         string
	ProductType   string
	ProductSize   string
	SizeFormat    string
	Quantity      int64
}
