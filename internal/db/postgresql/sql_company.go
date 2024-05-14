package postgresql

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"shop-aggregator/internal/model"
)

type Company struct {
	db *Client
}

func NewCompany(db *Client) *Company {
	return &Company{
		db: db,
	}
}

const (
	InsertCompanyQuery = `
		INSERT INTO company (company_name)
		VALUES ($1)
		RETURNING company_id`
	SelectCompaniesQuery     = `SELECT company_id, company_name FROM company WHERE company_name LIKE CONCAT(CAST($1 AS text), '%')`
	SelectCompanyByNameQuery = `SELECT company_id, company_name FROM company WHERE company_name = $1`
	SelectCompanyByIDQuery   = `SELECT company_id, company_name FROM company WHERE company_id = $1`
)

func (c *Company) Insert(ctx context.Context, company *model.Company) error {
	row := c.db.QueryRow(ctx, InsertCompanyQuery, company.CompanyName)
	err := row.Scan(&company.CompanyID)
	return err
}

func (c *Company) SelectCompanyByName(ctx context.Context, name string) (*model.Company, error) {
	rows := c.db.QueryRow(ctx, SelectCompanyByNameQuery, name)
	var company model.Company
	if err := rows.Scan(&company.CompanyID, &company.CompanyName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &company, nil
}

func (c *Company) SelectCompanyByID(ctx context.Context, companyID uuid.UUID) (*model.Company, error) {
	rows := c.db.QueryRow(ctx, SelectCompanyByIDQuery, companyID)
	var company model.Company
	if err := rows.Scan(&company.CompanyID, &company.CompanyName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &company, nil
}

func (c *Company) SelectCompanies(ctx context.Context, name string) ([]*model.Company, error) {
	rows, err := c.db.Query(ctx, SelectCompaniesQuery, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []*model.Company{}
	for rows.Next() {
		company := &model.Company{}
		err := rows.Scan(&company.CompanyID, &company.CompanyName)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}
