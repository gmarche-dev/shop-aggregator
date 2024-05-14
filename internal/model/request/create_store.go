package request

type CreateStore struct {
	Address     string `json:"address"`
	ZipCode     string `json:"zip_code"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Url         string `json:"url"`
	StoreName   string `json:"store_name"`
	StoreType   string `json:"store_type" binding:"required"`
	CompanyName string `json:"company_name" binding:"required"`
}
