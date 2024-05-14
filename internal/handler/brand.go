package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/model/request"
	"shop-aggregator/internal/model/response"
)

type BrandUseCase interface {
	Create(ctx context.Context, brandName string) (*model.Brand, error)
	SelectByPartialName(ctx context.Context, name string) ([]*model.Brand, error)
}

type Brand struct {
	BrandUseCase BrandUseCase
}

func NewBrand(bu BrandUseCase) *Brand {
	return &Brand{
		BrandUseCase: bu,
	}
}

func (b *Brand) Create(c *gin.Context) {
	var br request.CreateBrand
	if err := c.ShouldBindJSON(&br); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bm, err := b.BrandUseCase.Create(c.Request.Context(), br.BrandName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "brand created", "data": response.NewBrandFromModel(bm)})
}

func (b *Brand) GetByPartialName(c *gin.Context) {
	name := c.Param("name")
	bm, err := b.BrandUseCase.SelectByPartialName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.NewBrandsFromModels(bm)})
}
