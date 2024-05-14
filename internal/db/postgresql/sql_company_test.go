package postgresql

import (
	"context"
	"github.com/stretchr/testify/suite"
	"shop-aggregator/internal/model"
	"testing"
)

type SqlCompanyTestSuite struct {
	DBTestSuite
	Company *Company
}

func (s *SqlCompanyTestSuite) SetupTest() {
	s.Company = NewCompany(s.DB)
}

func (s *SqlCompanyTestSuite) TearDownTest() {
	_, err := s.DB.Exec(s.ctx, "TRUNCATE TABLE company")
	s.Require().NoError(err)
}

func (s *SqlCompanyTestSuite) TestCompany() {
	s.Run("no error", func() {
		expected := []*model.Company{
			{
				CompanyName: "company_1",
			},
			{
				CompanyName: "cmpany_1",
			},
			{
				CompanyName: "copany_1",
			},
			{
				CompanyName: "comany_1",
			},
		}
		for _, brand := range expected {
			err := s.Company.Insert(s.ctx, brand)
			s.NoError(err)
		}

		brands, err := s.Company.SelectCompanies(s.ctx, "c")
		s.NoError(err)
		s.ElementsMatch(expected, brands)

		searchCompanys, err := s.Company.SelectCompanies(s.ctx, "co")
		s.NoError(err)
		s.Len(searchCompanys, 3)

		searchCompanys, err = s.Company.SelectCompanies(s.ctx, "com")
		s.NoError(err)
		s.Len(searchCompanys, 2)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.EqualError(s.Company.Insert(ctx, &model.Company{}), `context canceled`)
		brands, err := s.Company.SelectCompanies(ctx, "my brand")
		s.EqualError(err, `context canceled`)
		s.Nil(brands)
	})
}

func TestCompanyTestSuite(t *testing.T) {
	suite.Run(t, new(SqlCompanyTestSuite))
}
