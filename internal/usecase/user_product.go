package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"shop-aggregator/internal/model"
)

type UserProductStorer interface {
	Insert(ctx context.Context, userProduct *model.UserProduct, userID uuid.UUID) error
	SelectProductsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.UserProduct, error)
	SelectProductsByUserIDAndStoreID(ctx context.Context, userID, storeID uuid.UUID) ([]*model.UserProduct, error)
	SelectProductsByBillID(ctx context.Context, billID uuid.UUID) ([]*model.UserProduct, error)
	SelectMostRecentUserProductByStoreID(ctx context.Context, storeID uuid.UUID) ([]*model.UserProduct, error)
	SelectProductByID(ctx context.Context, id uuid.UUID) (*model.UserProduct, error)
	UpdateQuantity(ctx context.Context, quantity int64, productType, productSize, sizeFormat string, userProductID uuid.UUID) error
	DeleteUserProduct(ctx context.Context, userProductID uuid.UUID) (uuid.UUID, error)
}

type UserProduct struct {
	UserProductStorer UserProductStorer
}

func NewUserProduct(ups UserProductStorer) *UserProduct {
	return &UserProduct{
		UserProductStorer: ups,
	}
}

func (up *UserProduct) Create(ctx context.Context, um *model.UserProduct, userID uuid.UUID) (*model.UserProduct, error) {
	if err := up.UserProductStorer.Insert(ctx, um, userID); err != nil {
		log.Error().Caller().Err(err).Msg("Create.Insert")
		return nil, model.ErrUserProductError
	}

	returnUp, err := up.UserProductStorer.SelectProductByID(ctx, um.UserProductID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Create.SelectProductByID")
		return nil, model.ErrUserProductError
	}

	if returnUp == nil {
		log.Error().Caller().Err(fmt.Errorf("returnUp is nil for id %s", um.UserProductID.String())).Msg("Create.SelectProductByID")
		return nil, model.ErrUserProductError
	}

	return returnUp, nil
}

func (up *UserProduct) SelectProductsByBillID(ctx context.Context, billID uuid.UUID) ([]*model.UserProduct, error) {
	ups, err := up.UserProductStorer.SelectProductsByBillID(ctx, billID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("SelectProductsByBillID.SelectProductsByBillID")
		return nil, model.ErrUserProductError
	}

	return ups, nil
}

func (up *UserProduct) UpdateQuantity(ctx context.Context, billID, userProductID uuid.UUID, productType, productSize, sizeFormat string, quantity int64) ([]*model.UserProduct, error) {
	if err := up.UserProductStorer.UpdateQuantity(ctx, quantity, productType, productSize, sizeFormat, userProductID); err != nil {
		log.Error().Caller().Err(err).Msg("UpdateQuantity.UpdateQuantity")
		return nil, model.ErrUserProductError
	}
	ups, err := up.UserProductStorer.SelectProductsByBillID(ctx, billID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("UpdateQuantity.SelectProductsByBillID")
		return nil, model.ErrUserProductError
	}

	return ups, nil
}

func (up *UserProduct) DeleteUserProduct(ctx context.Context, userProductID uuid.UUID) ([]*model.UserProduct, error) {
	billID, err := up.UserProductStorer.DeleteUserProduct(ctx, userProductID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("DeleteUserProduct.DeleteUserProduct")
		return nil, model.ErrUserProductError
	}
	ups, err := up.UserProductStorer.SelectProductsByBillID(ctx, billID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("DeleteUserProduct.SelectProductsByBillID")
		return nil, model.ErrUserProductError
	}

	return ups, nil
}
