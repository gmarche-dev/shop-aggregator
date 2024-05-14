package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"shop-aggregator/internal/model"
	"testing"
)

type SqlUserProductTestSuite struct {
	DBTestSuite
	Store       *Store
	Bill        *Bill
	UserProduct *UserProduct
	Product     *Product
	Brand       *Brand
}

func (s *SqlUserProductTestSuite) SetupTest() {
	s.Store = NewStore(s.DB)
	s.Bill = NewBill(s.DB)
	s.UserProduct = NewUserProduct(s.DB)
	s.Product = NewProduct(s.DB)
	s.Brand = NewBrand(s.DB)
}

func (s *SqlUserProductTestSuite) insertNewBill(userID, storeID uuid.UUID, amount string) *model.Bill {
	bill := &model.Bill{
		BillID:  uuid.New(),
		UserID:  userID,
		StoreID: storeID,
		Amount:  amount,
	}
	s.Require().NoError(s.Bill.Insert(s.ctx, bill))
	return bill
}

func (s *SqlUserProductTestSuite) insertNewStore(address, cp, name, city, country string) *model.Store {
	store := &model.Store{
		Address:   address,
		ZipCode:   cp,
		StoreName: name,
		City:      city,
		Country:   country,
		CompanyID: uuid.New(),
	}
	s.Require().NoError(s.Store.Insert(s.ctx, store))
	return store
}

func (s *SqlUserProductTestSuite) insertNewProduct(ean, productName string, brandID uuid.UUID) *model.Product {
	product := &model.Product{
		EAN:         ean,
		ProductName: productName,
		BrandID:     brandID,
	}

	s.Require().NoError(s.Product.Insert(s.ctx, product))
	return product
}

func (s *SqlUserProductTestSuite) insertNewBrand(brandName string) *model.Brand {
	brand := &model.Brand{
		BrandName: brandName,
	}

	s.Require().NoError(s.Brand.Insert(s.ctx, brand))
	return brand
}

func (s *SqlUserProductTestSuite) insertNewUserProduct(userID, productID, billID uuid.UUID, price string, quantity int64) *model.UserProduct {
	up := &model.UserProduct{
		ProductID: productID,
		BillID:    billID,
		Price:     price,
		Quantity:  quantity,
	}

	s.Require().NoError(s.UserProduct.Insert(s.ctx, up, userID))
	return up
}

func (s *SqlUserProductTestSuite) getUserProductByID(userProductID uuid.UUID) *model.UserProduct {
	up := &model.UserProduct{}
	row := s.DB.QueryRow(s.ctx, "SELECT user_product_id,product_id,bill_id,price,quantity FROM user_product where user_product_id = $1", userProductID)
	err := row.Scan(&up.UserProductID, &up.ProductID, &up.BillID, &up.Price, &up.Quantity)
	s.Require().NoError(err)
	return up
}

func (s *SqlUserProductTestSuite) TestInsert() {
	s.Run("no error", func() {
		insertUp := s.insertNewUserProduct(uuid.New(), uuid.New(), uuid.New(), "42.42", 42)
		s.NotEqual(uuid.Nil, insertUp.UserProductID)
		checkInsert := s.getUserProductByID(insertUp.UserProductID)
		s.Equal(insertUp, checkInsert)
	})

	s.Run("context error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.Require().EqualError(s.UserProduct.Insert(ctx, &model.UserProduct{}, uuid.New()), `context canceled`)
	})
}

func (s *SqlUserProductTestSuite) TestSelectProductsByUserID() {
	s.Run("no error", func() {
		userID := uuid.New()
		store := s.insertNewStore("7 rue du labrador", "02140", "intermarché", "vervins", "france")
		brand := s.insertNewBrand("brandName")
		product := s.insertNewProduct("ean13", "productName", brand.BrandID)
		bill := s.insertNewBill(userID, store.StoreID, "185")

		expectedUserProduct := model.UserProduct{
			ProductID:   product.ProductID,
			BillID:      bill.BillID,
			Price:       "42",
			Quantity:    42,
			ProductName: "productName",
			Ean:         product.EAN,
			StoreName:   "intermarché",
			StoreID:     store.StoreID,
			BrandName:   "brandName",
			BrandID:     brand.BrandID,
		}

		s.Require().NoError(s.UserProduct.Insert(s.ctx, &expectedUserProduct, userID))

		userProducts, err := s.UserProduct.SelectProductsByUserID(s.ctx, userID)
		s.Require().NoError(err)
		s.Require().Len(userProducts, 1)
		s.Equal(expectedUserProduct, *userProducts[0])
	})

	s.Run("context error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		us, err := s.UserProduct.SelectProductsByUserID(ctx, uuid.New())
		s.Require().Nil(us)
		s.Require().EqualError(err, `context canceled`)
	})
}

func (s *SqlUserProductTestSuite) TestSelectProductsByUserIDAndStoreID() {
	s.Run("no error", func() {
		userID := uuid.New()
		store := s.insertNewStore("7 rue du labrador", "02140", "intermarché", "vervins", "france")
		brand := s.insertNewBrand("brandName")
		product := s.insertNewProduct("ean13", "productName", brand.BrandID)
		bill := s.insertNewBill(userID, store.StoreID, "185")

		expectedUserProduct := model.UserProduct{
			ProductID:   product.ProductID,
			BillID:      bill.BillID,
			Price:       "42",
			Quantity:    42,
			ProductName: "productName",
			Ean:         product.EAN,
			StoreName:   "intermarché",
			StoreID:     store.StoreID,
			BrandName:   "brandName",
			BrandID:     brand.BrandID,
		}

		s.Require().NoError(s.UserProduct.Insert(s.ctx, &expectedUserProduct, userID))

		userProducts, err := s.UserProduct.SelectProductsByUserIDAndStoreID(s.ctx, userID, store.StoreID)
		s.Require().NoError(err)
		s.Require().Len(userProducts, 1)
		s.Equal(expectedUserProduct, *userProducts[0])
	})

	s.Run("context error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		us, err := s.UserProduct.SelectProductsByUserIDAndStoreID(ctx, uuid.New(), uuid.New())
		s.Require().Nil(us)
		s.Require().EqualError(err, `context canceled`)
	})
}

func (s *SqlUserProductTestSuite) TestSelectProductsByBillID() {
	s.Run("no error", func() {
		userID := uuid.New()
		store := s.insertNewStore("7 rue du labrador", "02140", "intermarché", "vervins", "france")
		brand := s.insertNewBrand("brandName")
		product := s.insertNewProduct("ean13", "productName", brand.BrandID)
		bill := s.insertNewBill(userID, store.StoreID, "185")

		expectedUserProduct := model.UserProduct{
			ProductID:   product.ProductID,
			BillID:      bill.BillID,
			Price:       "42",
			Quantity:    42,
			ProductName: "productName",
			Ean:         product.EAN,
			StoreName:   "intermarché",
			StoreID:     store.StoreID,
			BrandName:   "brandName",
			BrandID:     brand.BrandID,
		}

		s.Require().NoError(s.UserProduct.Insert(s.ctx, &expectedUserProduct, userID))

		userProducts, err := s.UserProduct.SelectProductsByBillID(s.ctx, bill.BillID)
		s.Require().NoError(err)
		s.Require().Len(userProducts, 1)
		s.Equal(expectedUserProduct, *userProducts[0])
	})

	s.Run("context error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		us, err := s.UserProduct.SelectProductsByBillID(ctx, uuid.New())
		s.Require().Nil(us)
		s.Require().EqualError(err, `context canceled`)
	})
}

func (s *SqlUserProductTestSuite) TestSelectMostRecentUserProductByStoreID() {
	s.Run("no error", func() {
		userID := uuid.New()
		store := s.insertNewStore("7 rue du labrador", "02140", "intermarché", "vervins", "france")
		brand := s.insertNewBrand("brandName")
		product := s.insertNewProduct("ean13", "productName", brand.BrandID)
		secondProduct := s.insertNewProduct("ean1312", "productName", brand.BrandID)
		bill := s.insertNewBill(userID, store.StoreID, "185")
		secondBill := s.insertNewBill(userID, store.StoreID, "345")

		expectedUserProduct := model.UserProduct{
			ProductID:   product.ProductID,
			BillID:      bill.BillID,
			Price:       "42",
			Quantity:    42,
			ProductName: "productName",
			Ean:         product.EAN,
			StoreName:   "intermarché",
			StoreID:     store.StoreID,
			BrandName:   "brandName",
			BrandID:     brand.BrandID,
		}

		s.Require().NoError(s.UserProduct.Insert(s.ctx, &expectedUserProduct, userID))

		expectedSecondUserProduct := model.UserProduct{
			ProductID:   product.ProductID,
			BillID:      bill.BillID,
			Price:       "42",
			Quantity:    42,
			ProductName: "productName",
			Ean:         product.EAN,
			StoreName:   "intermarché",
			StoreID:     store.StoreID,
			BrandName:   "brandName",
			BrandID:     brand.BrandID,
		}

		s.Require().NoError(s.UserProduct.Insert(s.ctx, &expectedSecondUserProduct, userID))

		expectedThirdUserProduct := model.UserProduct{
			ProductID:   secondProduct.ProductID,
			BillID:      secondBill.BillID,
			Price:       "42",
			Quantity:    42,
			ProductName: "productName",
			Ean:         secondProduct.EAN,
			StoreName:   "intermarché",
			StoreID:     store.StoreID,
			BrandName:   "brandName",
			BrandID:     brand.BrandID,
		}

		s.Require().NoError(s.UserProduct.Insert(s.ctx, &expectedThirdUserProduct, userID))

		expectedFourthUserProduct := model.UserProduct{
			ProductID:   secondProduct.ProductID,
			BillID:      secondBill.BillID,
			Price:       "42",
			Quantity:    42,
			ProductName: "productName",
			Ean:         secondProduct.EAN,
			StoreName:   "intermarché",
			StoreID:     store.StoreID,
			BrandName:   "brandName",
			BrandID:     brand.BrandID,
		}

		s.Require().NoError(s.UserProduct.Insert(s.ctx, &expectedFourthUserProduct, userID))

		userProducts, err := s.UserProduct.SelectMostRecentUserProductByStoreID(s.ctx, store.StoreID)
		s.Require().NoError(err)
		s.Require().Len(userProducts, 2)
		expectedResult := []*model.UserProduct{
			&expectedSecondUserProduct,
			&expectedFourthUserProduct,
		}
		s.ElementsMatch(expectedResult, userProducts)
	})

	s.Run("context error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		us, err := s.UserProduct.SelectMostRecentUserProductByStoreID(ctx, uuid.New())
		s.Require().Nil(us)
		s.Require().EqualError(err, `context canceled`)
	})
}

func (s *SqlUserProductTestSuite) TearDownTest() {
	_, err := s.DB.Exec(s.ctx, "TRUNCATE TABLE store")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE bill")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE product")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE user_product")
	s.Require().NoError(err)
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE brand")
	s.Require().NoError(err)
}

func TestUserProductTestSuite(t *testing.T) {
	suite.Run(t, new(SqlUserProductTestSuite))
}
