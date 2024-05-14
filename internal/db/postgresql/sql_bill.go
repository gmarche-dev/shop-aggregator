package postgresql

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"shop-aggregator/internal/model"
)

type Bill struct {
	db *Client
}

func NewBill(db *Client) *Bill {
	return &Bill{
		db: db,
	}
}

const (
	InsertBillQuery = `
		INSERT INTO bill (user_id, store_id, amount, bill_state)
		VALUES ($1, $2, $3, $4)
		RETURNING bill_id`
	UpdateBillQuery         = `UPDATE bill SET amount = $1, bill_state = $2 where bill_id = $3`
	GetBillByUserIDQuery    = `SELECT bill_id, user_id, store_id, amount, bill_state FROM bill WHERE user_id = $1`
	ExistsUnclosedBillQuery = `SELECT bill_id, user_id, store_id, amount, bill_state FROM bill WHERE user_id = $1 AND bill_state = $2`
)

func (b *Bill) Insert(ctx context.Context, bill *model.Bill) error {
	row := b.db.QueryRow(ctx, InsertBillQuery, bill.UserID, bill.StoreID, bill.Amount, model.BillStateCreate)
	err := row.Scan(&bill.BillID)
	return err
}

func (b *Bill) Update(ctx context.Context, bill *model.Bill) error {
	_, err := b.db.Exec(ctx, UpdateBillQuery, bill.Amount, bill.State, bill.BillID)
	return err
}

func (b *Bill) GetBillsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Bill, error) {
	rows, err := b.db.Query(ctx, GetBillByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*model.Bill
	for rows.Next() {
		bill := &model.Bill{}
		err = rows.Scan(&bill.BillID, &bill.UserID, &bill.StoreID, &bill.Amount, &bill.State)
		if err != nil {
			return nil, err
		}
		bills = append(bills, bill)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bills, nil
}

func (b *Bill) ExistsUnclosedBill(ctx context.Context, userID uuid.UUID) (*model.Bill, error) {
	row := b.db.QueryRow(ctx, ExistsUnclosedBillQuery, userID, model.BillStateCreate)
	bill := &model.Bill{}

	if err := row.Scan(&bill.BillID, &bill.UserID, &bill.StoreID, &bill.Amount, &bill.State); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return bill, nil
}
