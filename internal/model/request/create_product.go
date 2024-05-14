package request

type CreateProduct struct {
	EAN         string `json:"ean" binding:"required"`
	ProductName string `json:"product_name" binding:"required"`
	BrandName   string `json:"brand_name" binding:"required"`
}
