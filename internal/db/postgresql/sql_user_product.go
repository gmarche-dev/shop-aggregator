package postgresql

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"shop-aggregator/internal/model"
)

type UserProduct struct {
	db *Client
}

func NewUserProduct(db *Client) *UserProduct {
	return &UserProduct{
		db: db,
	}
}

const (
	InsertUserProductQuery = `
		INSERT INTO user_product (product_id, user_id, bill_id, price, quantity, product_type, product_size, size_format)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING user_product_id;`
	SelectProductsByUserIDQuery = `
		SELECT 
		    up.user_product_id, 
		    up.product_id,
		    p.product_name,
		    p.ean,
		    br.brand_id,
		    br.brand_name,
		    up.bill_id, 
		    s.store_id,
		    s.store_name,
		    up.price,
		    up.quantity,
		    up.product_type,
		    up.product_size,
		    up.size_format
		FROM user_product up 
		INNER JOIN product p ON up.product_id = p.product_id
		INNER JOIN brand br ON p.brand_id = br.brand_id
		INNER JOIN bill b ON up.bill_id = b.bill_id
		INNER JOIN store s ON b.store_id = s.store_id
		WHERE up.user_id = $1
		ORDER BY up.created_at`

	SelectProductByIDQuery = `
		SELECT 
		    up.user_product_id, 
		    up.product_id,
		    p.product_name,
		    p.ean,
		    br.brand_id,
		    br.brand_name,
		    up.bill_id, 
		    s.store_id,
		    s.store_name,
		    up.price,
		    up.quantity,
		    up.product_type,
		    up.product_size,
		    up.size_format
		FROM user_product up 
		INNER JOIN product p ON up.product_id = p.product_id
		INNER JOIN brand br ON p.brand_id = br.brand_id
		INNER JOIN bill b ON up.bill_id = b.bill_id
		INNER JOIN store s ON b.store_id = s.store_id
		WHERE up.user_product_id = $1
		ORDER BY up.created_at`

	SelectProductsByUserIDAndStoreIDQuery = `
		SELECT 
		    up.user_product_id, 
		    up.product_id,
		    p.product_name,
		    p.ean,
		    br.brand_id,
		    br.brand_name,
		    up.bill_id, 
		    s.store_id,
		    s.store_name,
		    up.price,
		    up.quantity,
		    up.product_type,
		    up.product_size,
		    up.size_format
		FROM user_product up 
		INNER JOIN product p ON up.product_id = p.product_id
		INNER JOIN brand br ON p.brand_id = br.brand_id
		INNER JOIN bill b ON up.bill_id = b.bill_id
		INNER JOIN store s ON b.store_id = s.store_id
		WHERE up.user_id = $1
		AND s.store_id = $2
		ORDER BY up.created_at`

	SelectProductsByBillIDQuery = `
		SELECT 
		    up.user_product_id, 
		    up.product_id,
		    p.product_name,
		    p.ean,
		    br.brand_id,
		    br.brand_name,
		    up.bill_id, 
		    s.store_id,
		    s.store_name,
		    up.price,
		    up.quantity,
		    up.product_type,
		    up.product_size,
		    up.size_format
		FROM user_product up 
		INNER JOIN product p ON up.product_id = p.product_id
		INNER JOIN brand br ON p.brand_id = br.brand_id
		INNER JOIN bill b ON up.bill_id = b.bill_id
		INNER JOIN store s ON b.store_id = s.store_id
		WHERE up.bill_id = $1
		ORDER BY up.created_at`

	SelectMostRecentUserProductByStoreIDQuery = `
		WITH RecentProducts AS (
			SELECT 
				up.product_id, 
				up.user_product_id,
				up.created_at,
				ROW_NUMBER() OVER (PARTITION BY up.product_id ORDER BY up.created_at DESC) as rn
			FROM 
				user_product up
			JOIN 
				bill b ON up.bill_id = b.bill_id
			WHERE 
				b.store_id = $1
		)
		SELECT 
		    up.user_product_id, 
		    up.product_id,
		    p.product_name,
		    p.ean,
		    br.brand_id,
		    br.brand_name,
		    up.bill_id, 
		    s.store_id,
		    s.store_name,
		    up.price,
		    up.quantity,
		    up.product_type,
		    up.product_size,
		    up.size_format
		FROM 
			RecentProducts rp
		JOIN 
			user_product up ON rp.user_product_id = up.user_product_id
		JOIN 
			bill b ON up.bill_id = b.bill_id
		JOIN 
			store s ON b.store_id = s.store_id
		JOIN 
			product p ON rp.product_id = p.product_id
		JOIN 
			brand br ON p.brand_id = br.brand_id
		WHERE 
			rp.rn = 1
		ORDER BY up.created_at`
	UpdateUserProductQuantityQuery = `UPDATE user_product set quantity = $1, product_type = $2, product_size = $3, size_format = $4 where user_product_id = $5`
	DeleteUserProduct              = `DELETE FROM user_product where user_product_id = $1 RETURNING bill_id`
)

func (up *UserProduct) Insert(ctx context.Context, userProduct *model.UserProduct, userID uuid.UUID) error {
	row := up.db.QueryRow(ctx, InsertUserProductQuery, userProduct.ProductID, userID, userProduct.BillID, userProduct.Price, userProduct.Quantity, userProduct.ProductType, userProduct.ProductSize, userProduct.SizeFormat)
	err := row.Scan(&userProduct.UserProductID)
	return err
}

func (up *UserProduct) SelectProductsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.UserProduct, error) {
	rows, err := up.db.Query(ctx, SelectProductsByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userProducts []*model.UserProduct
	for rows.Next() {
		userProduct := model.UserProduct{}
		err := rows.Scan(
			&userProduct.UserProductID,
			&userProduct.ProductID,
			&userProduct.ProductName,
			&userProduct.Ean,
			&userProduct.BrandID,
			&userProduct.BrandName,
			&userProduct.BillID,
			&userProduct.StoreID,
			&userProduct.StoreName,
			&userProduct.Price,
			&userProduct.Quantity,
			&userProduct.ProductType,
			&userProduct.ProductSize,
			&userProduct.SizeFormat,
		)
		if err != nil {
			return nil, err
		}
		userProducts = append(userProducts, &userProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userProducts, nil
}

func (up *UserProduct) SelectProductsByUserIDAndStoreID(ctx context.Context, userID, storeID uuid.UUID) ([]*model.UserProduct, error) {
	rows, err := up.db.Query(ctx, SelectProductsByUserIDAndStoreIDQuery, userID, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userProducts []*model.UserProduct
	for rows.Next() {
		userProduct := model.UserProduct{}
		err := rows.Scan(
			&userProduct.UserProductID,
			&userProduct.ProductID,
			&userProduct.ProductName,
			&userProduct.Ean,
			&userProduct.BrandID,
			&userProduct.BrandName,
			&userProduct.BillID,
			&userProduct.StoreID,
			&userProduct.StoreName,
			&userProduct.Price,
			&userProduct.Quantity,
			&userProduct.ProductType,
			&userProduct.ProductSize,
			&userProduct.SizeFormat,
		)
		if err != nil {
			return nil, err
		}
		userProducts = append(userProducts, &userProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userProducts, nil
}

func (up *UserProduct) SelectProductsByBillID(ctx context.Context, billID uuid.UUID) ([]*model.UserProduct, error) {
	rows, err := up.db.Query(ctx, SelectProductsByBillIDQuery, billID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userProducts []*model.UserProduct
	for rows.Next() {
		userProduct := model.UserProduct{}
		err := rows.Scan(
			&userProduct.UserProductID,
			&userProduct.ProductID,
			&userProduct.ProductName,
			&userProduct.Ean,
			&userProduct.BrandID,
			&userProduct.BrandName,
			&userProduct.BillID,
			&userProduct.StoreID,
			&userProduct.StoreName,
			&userProduct.Price,
			&userProduct.Quantity,
			&userProduct.ProductType,
			&userProduct.ProductSize,
			&userProduct.SizeFormat,
		)
		if err != nil {
			return nil, err
		}
		userProducts = append(userProducts, &userProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userProducts, nil
}

func (up *UserProduct) SelectMostRecentUserProductByStoreID(ctx context.Context, storeID uuid.UUID) ([]*model.UserProduct, error) {
	rows, err := up.db.Query(ctx, SelectMostRecentUserProductByStoreIDQuery, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userProducts []*model.UserProduct
	for rows.Next() {
		userProduct := model.UserProduct{}
		err := rows.Scan(
			&userProduct.UserProductID,
			&userProduct.ProductID,
			&userProduct.ProductName,
			&userProduct.Ean,
			&userProduct.BrandID,
			&userProduct.BrandName,
			&userProduct.BillID,
			&userProduct.StoreID,
			&userProduct.StoreName,
			&userProduct.Price,
			&userProduct.Quantity,
			&userProduct.ProductType,
			&userProduct.ProductSize,
			&userProduct.SizeFormat,
		)
		if err != nil {
			return nil, err
		}
		userProducts = append(userProducts, &userProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userProducts, nil
}

func (up *UserProduct) SelectProductByID(ctx context.Context, id uuid.UUID) (*model.UserProduct, error) {
	row := up.db.QueryRow(ctx, SelectProductByIDQuery, id)
	userProduct := model.UserProduct{}
	err := row.Scan(
		&userProduct.UserProductID,
		&userProduct.ProductID,
		&userProduct.ProductName,
		&userProduct.Ean,
		&userProduct.BrandID,
		&userProduct.BrandName,
		&userProduct.BillID,
		&userProduct.StoreID,
		&userProduct.StoreName,
		&userProduct.Price,
		&userProduct.Quantity,
		&userProduct.ProductType,
		&userProduct.ProductSize,
		&userProduct.SizeFormat,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &userProduct, nil
}

func (up *UserProduct) UpdateQuantity(ctx context.Context, quantity int64, productType, productSize, sizeFormat string, userProductID uuid.UUID) error {
	_, err := up.db.Exec(ctx, UpdateUserProductQuantityQuery, quantity, productType, productSize, sizeFormat, userProductID)
	return err
}
func (up *UserProduct) DeleteUserProduct(ctx context.Context, userProductID uuid.UUID) (uuid.UUID, error) {
	row := up.db.QueryRow(ctx, DeleteUserProduct, userProductID)
	var billID uuid.UUID
	if err := row.Scan(&billID); err != nil {
		return uuid.Nil, err
	}
	return billID, nil
}
