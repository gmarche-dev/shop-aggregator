package usecase

import (
	"context"
	"github.com/rs/zerolog/log"
	"shop-aggregator/internal/model"
)

type StoreStorer interface {
	Insert(ctx context.Context, store *model.Store) error
	SelectStoresByZipCodeOrName(ctx context.Context, storeType, search string) ([]*model.Store, error)
}

type CompanyStoreStorer interface {
	Insert(ctx context.Context, company *model.Company) error
	SelectCompanyByName(ctx context.Context, name string) (*model.Company, error)
}

type Store struct {
	StoreStorer   StoreStorer
	CompanyStorer CompanyStoreStorer
}

func NewStore(ss StoreStorer, cs CompanyStoreStorer) *Store {
	return &Store{
		StoreStorer:   ss,
		CompanyStorer: cs,
	}
}

func (s *Store) createOrGetCompany(ctx context.Context, companyName string) (*model.Company, error) {
	exists, err := s.CompanyStorer.SelectCompanyByName(ctx, companyName)
	if err != nil {
		log.Error().Caller().Err(err).Msg("createOrGetCompany.SelectCompanyByName")
		return nil, model.ErrCompanyError
	}
	if exists != nil {
		return exists, nil
	}
	company := model.Company{
		CompanyName: companyName,
	}
	if err := s.CompanyStorer.Insert(ctx, &company); err != nil {
		log.Error().Caller().Err(err).Msg("createOrGetCompany.Insert")
		return nil, model.ErrInsertCompanyError
	}

	return &company, nil
}

func (s *Store) CreateStore(ctx context.Context, store *model.Store, companyName string) (*model.Store, error) {
	company, err := s.createOrGetCompany(ctx, companyName)
	if err != nil {
		return nil, err
	}
	store.CompanyID = company.CompanyID
	var search string
	if store.StoreType == model.StoreTypeShop {
		search = store.ZipCode
	} else {
		search = store.StoreName
	}
	stores, err := s.StoreStorer.SelectStoresByZipCodeOrName(ctx, store.StoreType, search)
	if err != nil {
		log.Error().Caller().Err(err).Msg("CreateStore.SelectStoresByZipCodeOrName")
		return nil, model.ErrStoreError
	}

	if checkStore := sameStore(stores, store); checkStore != nil {
		return checkStore, nil
	}

	if err = s.StoreStorer.Insert(ctx, store); err != nil {
		log.Error().Caller().Err(err).Msg("CreateStore.Insert")
		return nil, model.ErrStoreError
	}

	return store, nil
}

func (s *Store) GetStoreByZipCodeOrName(ctx context.Context, storeType, search string) ([]*model.Store, error) {
	stores, err := s.StoreStorer.SelectStoresByZipCodeOrName(ctx, storeType, search)
	if err != nil {
		log.Error().Caller().Err(err).Msg("GetStoreByZipCodeOrName.SelectStoresByZipCodeOrName")
		return nil, model.ErrStoreError
	}

	return stores, nil
}

func sameStore(olds []*model.Store, new *model.Store) *model.Store {
	for _, old := range olds {
		if old.CompanyID == new.CompanyID {
			if old.City == new.City {
				if old.Address == new.Address {
					return old
				}
			}
		}
	}
	return nil
}
