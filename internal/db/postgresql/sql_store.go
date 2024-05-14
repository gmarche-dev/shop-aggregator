package postgresql

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"shop-aggregator/internal/model"
)

type Store struct {
	db *Client
}

func NewStore(db *Client) *Store {
	return &Store{
		db: db,
	}
}

const (
	InsertStoreQuery = `
		INSERT INTO store (address, zip_code, city, country, store_name, store_type, url, company_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING store_id`
	SelectStoresByZipCodeQuery = `SELECT store_id,address, zip_code, city, country, store_name, store_type, url, company_id FROM store where zip_code LIKE CONCAT(CAST($1 AS text), '%')`
	SelectStoresByNameQuery    = `SELECT store_id,address, zip_code, city, country, store_name, store_type, url, company_id FROM store where store_type = $1 AND store_name LIKE CONCAT('%', CAST($2 AS text), '%')`
	SelectStoresByIDQuery      = `SELECT store_id,address, zip_code, city, country, store_name, store_type, url, company_id FROM store where store_id = $1`
)

func (s *Store) Insert(ctx context.Context, store *model.Store) error {
	row := s.db.QueryRow(ctx, InsertStoreQuery, store.Address, store.ZipCode, store.City, store.Country, store.StoreName, store.StoreType, store.Url, store.CompanyID)
	err := row.Scan(&store.StoreID)
	return err
}

func (s *Store) SelectStoresByZipCodeOrName(ctx context.Context, storeType, search string) ([]*model.Store, error) {
	var rows pgx.Rows
	var err error
	if storeType == model.StoreTypeWeb {
		rows, err = s.db.Query(ctx, SelectStoresByNameQuery, storeType, search)
	} else {
		rows, err = s.db.Query(ctx, SelectStoresByZipCodeQuery, search)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := []*model.Store{}
	for rows.Next() {
		store := &model.Store{}
		err := rows.Scan(&store.StoreID, &store.Address, &store.ZipCode, &store.City, &store.Country, &store.StoreName, &store.StoreType, &store.Url, &store.CompanyID)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stores, nil
}

func (s *Store) SelectStoreByID(ctx context.Context, storeID uuid.UUID) (*model.Store, error) {
	rows := s.db.QueryRow(ctx, SelectStoresByIDQuery, storeID)
	store := &model.Store{}
	if err := rows.Scan(&store.StoreID, &store.Address, &store.ZipCode, &store.City, &store.Country, &store.StoreName, &store.StoreType, &store.Url, &store.CompanyID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return store, nil
}
