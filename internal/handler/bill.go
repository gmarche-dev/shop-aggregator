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

type BillUseCase interface {
	CloseBill(ctx context.Context, userID, billID uuid.UUID, amount string) error
	StartBill(ctx context.Context, userID, storeID uuid.UUID) (*response.Bill, error)
	GetBillsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Bill, error)
	GetLastBill(ctx context.Context, userID uuid.UUID) (*response.Bill, error)
	CancelBill(ctx context.Context, userID, billID uuid.UUID) error
}

type Bill struct {
	BillUseCase BillUseCase
}

func NewBill(bu BillUseCase) *Bill {
	return &Bill{
		BillUseCase: bu,
	}
}

func (b *Bill) Start(c *gin.Context) {
	var sb request.StartBill
	if err := c.ShouldBindJSON(&sb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	bill, err := b.BillUseCase.StartBill(c.Request.Context(), uuid.MustParse(id.(string)), sb.StoreID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bill started", "data": bill})
}

func (b *Bill) Close(c *gin.Context) {
	var sb request.CloseBill
	if err := c.ShouldBindJSON(&sb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := b.BillUseCase.CloseBill(c.Request.Context(), uuid.MustParse(id.(string)), sb.BillID, sb.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bill close"})
}

func (b *Bill) Cancel(c *gin.Context) {
	var sb request.CancelBill
	if err := c.ShouldBindJSON(&sb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := b.BillUseCase.CancelBill(c.Request.Context(), uuid.MustParse(id.(string)), sb.BillID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bill canceled"})
}

func (b *Bill) GetBillsByUserID(c *gin.Context) {
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	bills, err := b.BillUseCase.GetBillsByUserID(c.Request.Context(), uuid.MustParse(id.(string)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	billResponse := response.NewBillsFromModels(bills)

	c.JSON(http.StatusOK, gin.H{"message": "get bills", "data": billResponse})
}

func (b *Bill) GetLastBill(c *gin.Context) {
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	bill, err := b.BillUseCase.GetLastBill(c.Request.Context(), uuid.MustParse(id.(string)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if bill == nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "get last bill", "data": bill})
}
