package handler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/model/request"
	"shop-aggregator/internal/model/response"
)

type ProductUseCase interface {
	Create(ctx context.Context, m *model.Product, brandName string) (*model.Product, error)
	GetProductByEAN(ctx context.Context, ean string) (*model.Product, error)
}

type Product struct {
	ProductUseCase ProductUseCase
}

func NewProduct(pu ProductUseCase) *Product {
	return &Product{
		ProductUseCase: pu,
	}
}

func (p *Product) Create(c *gin.Context) {
	var cp request.CreateProduct
	if err := c.ShouldBindJSON(&cp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pm := newProductFromRequest(&cp)
	pm, err := p.ProductUseCase.Create(c.Request.Context(), pm, cp.BrandName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "company created", "data": response.NewProductFromModel(pm)})
}

func (p *Product) GetProductByEAN(c *gin.Context) {
	ean := c.Param("ean")
	product, err := p.ProductUseCase.GetProductByEAN(c.Request.Context(), ean)
	if err != nil {
		if errors.Is(err, model.ErrNotExistsError) {
			c.JSON(http.StatusNoContent, gin.H{"message": "product not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "company created", "data": response.NewProductFromModel(product)})
}

func newProductFromRequest(r *request.CreateProduct) *model.Product {
	return &model.Product{
		EAN:         r.EAN,
		ProductName: r.ProductName,
	}
}
