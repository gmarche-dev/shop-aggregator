package request

type CreateBrand struct {
	BrandName string `json:"name" binding:"required"`
}
