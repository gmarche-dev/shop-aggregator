package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/model/response"
)

type CompanyUseCase interface {
	SelectByPartialName(ctx context.Context, partialName string) ([]*model.Company, error)
}

type Company struct {
	CompanyUseCase CompanyUseCase
}

func NewCompany(cu CompanyUseCase) *Company {
	return &Company{
		CompanyUseCase: cu,
	}
}

func (co *Company) GetByPartialName(c *gin.Context) {
	name := c.Param("name")
	cm, err := co.CompanyUseCase.SelectByPartialName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.NewCompaniesFromModels(cm)})
}
