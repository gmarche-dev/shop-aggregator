package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/model/response"
)

type BillStorer interface {
	Insert(ctx context.Context, bill *model.Bill) error
	Update(ctx context.Context, bill *model.Bill) error
	GetBillsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Bill, error)
	ExistsUnclosedBill(ctx context.Context, userID uuid.UUID) (*model.Bill, error)
}

type BillStoreStorer interface {
	SelectStoreByID(ctx context.Context, storeID uuid.UUID) (*model.Store, error)
}

type BillCompanyStorer interface {
	SelectCompanyByID(ctx context.Context, companyID uuid.UUID) (*model.Company, error)
}

type BillUserProductsStorer interface {
	SelectProductsByBillID(ctx context.Context, billID uuid.UUID) ([]*model.UserProduct, error)
}

type Bill struct {
	BillStorer             BillStorer
	BillStoreStorer        BillStoreStorer
	BillCompanyStorer      BillCompanyStorer
	BillUserProductsStorer BillUserProductsStorer
}

func NewBill(
	bs BillStorer,
	bss BillStoreStorer,
	bcs BillCompanyStorer,
	bups BillUserProductsStorer,
) *Bill {
	return &Bill{
		BillStorer:             bs,
		BillStoreStorer:        bss,
		BillCompanyStorer:      bcs,
		BillUserProductsStorer: bups,
	}
}

func (b *Bill) StartBill(ctx context.Context, userID, storeID uuid.UUID) (*response.Bill, error) {
	bill, err := b.BillStorer.ExistsUnclosedBill(ctx, userID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("StartBill.ExistsUnclosedBill")
		return nil, model.ErrBillError
	}
	if bill == nil {
		bill = &model.Bill{
			UserID:  userID,
			StoreID: storeID,
			Amount:  "0.0",
			State:   model.BillStateCreate,
		}
		if err = b.BillStorer.Insert(ctx, bill); err != nil {
			log.Error().Caller().Err(err).Msg("StartBill.Insert")
			return nil, model.ErrBillError
		}
	}

	return b.prepareBillResponse(ctx, bill)
}

func (b *Bill) CloseBill(ctx context.Context, userID, billID uuid.UUID, amount string) error {
	bill, err := b.BillStorer.ExistsUnclosedBill(ctx, userID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("CloseBill.ExistsUnclosedBill")
		return model.ErrBillError
	}
	if bill == nil {
		return model.ErrBillError
	}
	if bill.BillID != billID {
		return model.ErrBillError
	}

	bill.State = model.BillStateCompleted
	bill.Amount = amount
	err = b.BillStorer.Update(ctx, bill)
	if err != nil {
		log.Error().Caller().Err(err).Msg("CloseBill.Update")
		return model.ErrBillError
	}

	return err
}

func (b *Bill) CancelBill(ctx context.Context, userID, billID uuid.UUID) error {
	bill, err := b.BillStorer.ExistsUnclosedBill(ctx, userID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("CancelBill.ExistsUnclosedBill")
		return model.ErrBillError
	}
	if bill == nil {
		return nil
	}
	if bill.BillID != billID {
		return model.ErrBillError
	}

	bill.State = model.BillStateCanceled
	bill.Amount = "0"
	err = b.BillStorer.Update(ctx, bill)
	if err != nil {
		log.Error().Caller().Err(err).Msg("CancelBill.Update")
		return model.ErrBillError
	}

	return err
}

func (b *Bill) GetBillsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Bill, error) {
	bills, err := b.BillStorer.GetBillsByUserID(ctx, userID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("GetBillsByUserID.GetBillsByUserID")
		return nil, model.ErrBillError
	}
	return bills, nil
}

func (b *Bill) GetLastBill(ctx context.Context, userID uuid.UUID) (*response.Bill, error) {
	bill, err := b.BillStorer.ExistsUnclosedBill(ctx, userID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("GetLastBill.ExistsUnclosedBill")
		return nil, model.ErrBillError
	}

	return b.prepareBillResponse(ctx, bill)
}

func (b *Bill) prepareBillResponse(ctx context.Context, bill *model.Bill) (*response.Bill, error) {
	if bill == nil {
		return nil, nil
	}
	products, err := b.BillUserProductsStorer.SelectProductsByBillID(ctx, bill.BillID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("prepareBillResponse.SelectProductsByBillID")
		return nil, err
	}

	store, err := b.BillStoreStorer.SelectStoreByID(ctx, bill.StoreID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("prepareBillResponse.SelectStoreByID")
		return nil, model.ErrStoreError
	}

	company, err := b.BillCompanyStorer.SelectCompanyByID(ctx, store.CompanyID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("prepareBillResponse.SelectCompanyByID")
		return nil, model.ErrStoreError
	}

	return response.NewBillFromModel(bill, store, company, products), nil
}
