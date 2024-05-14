package request

type CreateCompany struct {
	CompanyName string `json:"name" binding:"required"`
}
