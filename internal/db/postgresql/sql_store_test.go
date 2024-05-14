package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"shop-aggregator/internal/model"
	"testing"
)

type SqlStoreTestSuite struct {
	DBTestSuite
	Store *Store
}

func (s *SqlStoreTestSuite) SetupTest() {
	s.Store = NewStore(s.DB)
}

func (s *SqlStoreTestSuite) TearDownTest() {
	_, err := s.DB.Exec(s.ctx, "TRUNCATE TABLE store")
	s.Require().NoError(err)
}

func (s *SqlStoreTestSuite) TestStore() {
	s.Run("no error", func() {
		expected := []*model.Store{
			{
				StoreName: "store_1",
				Address:   "address_1",
				ZipCode:   "78600",
				City:      "city_1",
				Country:   "country_1",
				CompanyID: uuid.New(),
			},
			{
				StoreName: "store_2",
				Address:   "address_2",
				ZipCode:   "78500",
				City:      "city_2",
				Country:   "country_2",
				CompanyID: uuid.New(),
			},
			{
				StoreName: "store_3",
				Address:   "address_3",
				ZipCode:   "78610",
				City:      "city_3",
				Country:   "country_3",
				CompanyID: uuid.New(),
			},
			{
				StoreName: "store_4",
				Address:   "address_4",
				ZipCode:   "78601",
				City:      "city_4",
				Country:   "country_4",
				CompanyID: uuid.New(),
			},
		}
		for _, brand := range expected {
			err := s.Store.Insert(s.ctx, brand)
			s.NoError(err)
		}

		brands, err := s.Store.SelectStoresByZipCodeOrName(s.ctx, model.StoreTypeShop, "7")
		s.NoError(err)
		s.ElementsMatch(expected, brands)

		searchStores, err := s.Store.SelectStoresByZipCodeOrName(s.ctx, model.StoreTypeShop, "786")
		s.NoError(err)
		s.Len(searchStores, 3)

		searchStores, err = s.Store.SelectStoresByZipCodeOrName(s.ctx, model.StoreTypeShop, "7860")
		s.NoError(err)
		s.Len(searchStores, 2)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.EqualError(s.Store.Insert(ctx, &model.Store{}), `context canceled`)
		brands, err := s.Store.SelectStoresByZipCodeOrName(s.ctx, model.StoreTypeShop, "oh no")
		s.EqualError(err, `context canceled`)
		s.Nil(brands)
	})
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(SqlStoreTestSuite))
}
