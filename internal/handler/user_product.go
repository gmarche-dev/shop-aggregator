package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/model/request"
	"shop-aggregator/internal/model/response"
)

type UserProductUseCase interface {
	Create(ctx context.Context, um *model.UserProduct, userID uuid.UUID) (*model.UserProduct, error)
	SelectProductsByBillID(ctx context.Context, billID uuid.UUID) ([]*model.UserProduct, error)
	UpdateQuantity(ctx context.Context, billID, userProductID uuid.UUID, productType, productSize, sizeFormat string, quantity int64) ([]*model.UserProduct, error)
	DeleteUserProduct(ctx context.Context, userProductID uuid.UUID) ([]*model.UserProduct, error)
}

type UserProduct struct {
	UserProductUseCase UserProductUseCase
}

func NewUserProduct(upu UserProductUseCase) *UserProduct {
	return &UserProduct{
		UserProductUseCase: upu,
	}
}

func (up *UserProduct) Create(c *gin.Context) {
	var cp request.CreateUserProduct
	if err := c.ShouldBindJSON(&cp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	pum := newUserProductFromRequest(&cp)
	pum, err := up.UserProductUseCase.Create(c.Request.Context(), pum, uuid.MustParse(id.(string)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user product created", "data": response.NewUserProductFromModel(pum)})
}

func (up *UserProduct) SelectProductsByBillID(c *gin.Context) {
	id, exists := c.Get("bill_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bill not found"})
		return
	}

	pum, err := up.UserProductUseCase.SelectProductsByBillID(c.Request.Context(), uuid.MustParse(id.(string)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user products", "data": response.NewUserProductsFromModel(pum)})
}

func (up *UserProduct) UpdateQuantity(c *gin.Context) {
	var uupq request.UpdateUserProductQuantity
	if err := c.ShouldBindJSON(&uupq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pum, err := up.UserProductUseCase.UpdateQuantity(c.Request.Context(), uupq.BillID, uupq.UserProductID, uupq.ProductType, uupq.ProductSize, uupq.SizeFormat, uupq.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user products", "data": response.NewUserProductsFromModel(pum)})
}

func (up *UserProduct) Delete(c *gin.Context) {
	userProductID := c.Param("user_product_id")
	pum, err := up.UserProductUseCase.DeleteUserProduct(c.Request.Context(), uuid.MustParse(userProductID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user products", "data": response.NewUserProductsFromModel(pum)})
}

func newUserProductFromRequest(r *request.CreateUserProduct) *model.UserProduct {
	return &model.UserProduct{
		ProductID:   r.ProductID,
		BillID:      r.BillID,
		Price:       r.Price,
		Quantity:    r.Quantity,
		ProductSize: r.ProductSize,
		ProductType: r.ProductType,
		SizeFormat:  r.SizeFormat,
	}
}
