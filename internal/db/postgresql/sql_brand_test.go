package postgresql

import (
	"context"
	"github.com/stretchr/testify/suite"
	"shop-aggregator/internal/model"
	"testing"
)

type SqlBrandTestSuite struct {
	DBTestSuite
	Brand *Brand
}

func (s *SqlBrandTestSuite) SetupTest() {
	s.Brand = NewBrand(s.DB)
}

func (s *SqlBrandTestSuite) TearDownTest() {
	_, err := s.DB.Exec(s.ctx, "TRUNCATE TABLE brand")
	s.Require().NoError(err)
}

func (s *SqlBrandTestSuite) TestBrand() {
	s.Run("no error", func() {
		expected := []*model.Brand{
			{
				BrandName: "company_1",
			},
			{
				BrandName: "cmpany_1",
			},
			{
				BrandName: "copany_1",
			},
			{
				BrandName: "comany_1",
			},
		}
		for _, brand := range expected {
			err := s.Brand.Insert(s.ctx, brand)
			s.NoError(err)
		}

		brands, err := s.Brand.SelectBrands(s.ctx, "c")
		s.NoError(err)
		s.ElementsMatch(expected, brands)

		searchBrands, err := s.Brand.SelectBrands(s.ctx, "co")
		s.NoError(err)
		s.Len(searchBrands, 3)

		searchBrands, err = s.Brand.SelectBrands(s.ctx, "com")
		s.NoError(err)
		s.Len(searchBrands, 2)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.EqualError(s.Brand.Insert(ctx, &model.Brand{}), `context canceled`)
		brands, err := s.Brand.SelectBrands(ctx, "my brand")
		s.EqualError(err, `context canceled`)
		s.Nil(brands)
	})
}

func TestBrandTestSuite(t *testing.T) {
	suite.Run(t, new(SqlBrandTestSuite))
}
