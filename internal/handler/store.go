package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/model/request"
	"shop-aggregator/internal/model/response"
)

type StoreUseCase interface {
	CreateStore(ctx context.Context, store *model.Store, companyName string) (*model.Store, error)
	GetStoreByZipCodeOrName(ctx context.Context, storeType, search string) ([]*model.Store, error)
}

type Store struct {
	StoreUseCase StoreUseCase
}

func NewStore(su StoreUseCase) *Store {
	return &Store{
		StoreUseCase: su,
	}
}

func (s *Store) CreateStore(c *gin.Context) {
	var cs request.CreateStore
	if err := c.ShouldBindJSON(&cs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cs.StoreType == model.StoreTypeShop {
		var errStore error
		if cs.Address == "" {
			errStore = fmt.Errorf("%w address needed \n", errStore)
		}
		if cs.ZipCode == "" {
			errStore = fmt.Errorf("%w zip code needed \n", errStore)
		}
		if cs.Country == "" {
			errStore = fmt.Errorf("%w zip country needed \n", errStore)
		}
		if cs.City == "" {
			errStore = fmt.Errorf("%w zip city needed \n", errStore)
		}
		if errStore != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errStore.Error()})
			return
		}
		cs.Url = ""
	}

	if cs.StoreType == model.StoreTypeWeb {
		var errStore error
		if cs.Url == "" {
			errStore = fmt.Errorf("%w url needed \n", errStore)
		}
		if errStore != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errStore.Error()})
			return
		}

		cs.Address = ""
		cs.ZipCode = ""
		cs.Country = ""
		cs.City = ""
		cs.StoreName = cs.Url
	}

	sm := newStoreFromRequest(cs)
	sm, err := s.StoreUseCase.CreateStore(c.Request.Context(), sm, cs.CompanyName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "company created", "data": response.NewStoreFromModel(sm)})
}

func (s *Store) GetStoreByZipCodeOrName(c *gin.Context) {
	storeType := c.Param("store_type")
	search := c.Param("search")
	stores, err := s.StoreUseCase.GetStoreByZipCodeOrName(c.Request.Context(), storeType, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "company created", "data": response.NewStoresFromModels(stores)})
}

func newStoreFromRequest(r request.CreateStore) *model.Store {
	return &model.Store{
		Address:   r.Address,
		ZipCode:   r.ZipCode,
		City:      r.City,
		Country:   r.Country,
		StoreName: r.StoreName,
		Url:       r.Url,
		StoreType: r.StoreType,
	}
}
