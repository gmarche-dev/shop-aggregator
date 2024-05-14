package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/model/response"
)

type Initialisation struct {
}

func NewInitialisation() *Initialisation {
	return &Initialisation{}
}

func (i *Initialisation) AppInitialisation(c *gin.Context) {
	ai := response.AppInitialisation{
		StoreTypes: []response.AppInitialisationRow{
			{Name: "shop", Field: model.StoreTypeShop},
			{Name: "web", Field: model.StoreTypeWeb},
		},
		BulkProducts: []response.AppInitialisationBulkProducts{
			{Name: "Meat", Field: model.BulkProductIDMeat.String(), Format: model.SizeFormatSizeTypeWeight},
			{Name: "Vegetable", Field: model.BulkProductIDVegetable.String(), Format: model.SizeFormatSizeTypeWeight},
			{Name: "Fruit", Field: model.BulkProductIDFruits.String(), Format: model.SizeFormatSizeTypeWeight},
		},
		ProductTypes: []response.AppInitialisationRow{
			{Name: "Bulk", Field: model.ProductBulk},
			{Name: "Barcoded", Field: model.ProductBarcoded},
		},
		Formats: map[string]map[string]response.AppInitialisationFormat{
			model.SizeFormatVolume: {
				model.SizeFormatVolumeMl: {
					Field: model.SizeFormatVolumeMl,
					Name:  "milliliter",
					Conversion: map[string]float64{
						model.SizeFormatVolumeL: 1000,
					},
				},
				model.SizeFormatVolumeL: {
					Field: model.SizeFormatVolumeL,
					Name:  "liter",
					Conversion: map[string]float64{
						model.SizeFormatVolumeMl: 0.001,
					},
				},
			},
			model.SizeFormatSizeTypeWeight: {
				model.SizeFormatWeightGr: {
					Field: model.SizeFormatWeightGr,
					Name:  "grams",
					Conversion: map[string]float64{
						model.SizeFormatWeightKg: 1000,
					},
				},
				model.SizeFormatWeightKg: {
					Field: model.SizeFormatWeightKg,
					Name:  "liter",
					Conversion: map[string]float64{
						model.SizeFormatWeightGr: 0.001,
					},
				},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{"message": "app initialisation", "data": ai})
}
