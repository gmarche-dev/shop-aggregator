package response

type AppInitialisation struct {
	StoreTypes   []AppInitialisationRow                        `json:"store_type"`
	ProductTypes []AppInitialisationRow                        `json:"product_types"`
	BulkProducts []AppInitialisationBulkProducts               `json:"bulk_products"`
	Formats      map[string]map[string]AppInitialisationFormat `json:"formats"`
}

type AppInitialisationRow struct {
	Name  string `json:"name"`
	Field string `json:"field"`
}

type AppInitialisationBulkProducts struct {
	Field  string `json:"field"`
	Name   string `json:"name"`
	Format string `json:"format"`
}

type AppInitialisationFormat struct {
	Field      string             `json:"field"`
	Name       string             `json:"name"`
	Type       string             `json:"type"`
	Conversion map[string]float64 `json:"conversion"`
}
