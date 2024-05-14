package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"shop-aggregator/internal/model"
	"testing"
)

type SqlProductTestSuite struct {
	DBTestSuite
	Product *Product
	product model.Product
}

func (s *SqlProductTestSuite) SetupTest() {
	s.Product = NewProduct(s.DB)
	s.product = model.Product{
		ProductID:   uuid.New(),
		EAN:         "eam",
		ProductName: "product-name",
		BrandID:     uuid.New(),
	}
}

func (s *SqlProductTestSuite) TearDownTest() {
	_, err := s.DB.Exec(s.ctx, "TRUNCATE TABLE product")
	s.Require().NoError(err)
}

func (s *SqlProductTestSuite) TestProduct() {
	s.Run("insert and get, no error", func() {
		product := s.product
		checkProduct, err := s.Product.GetProductByEAN(s.ctx, product.EAN)
		s.NoError(err)
		s.Nil(checkProduct)

		// insert product
		s.NoError(s.Product.Insert(s.ctx, &product))

		// check well insert
		checkProduct, err = s.Product.GetProductByEAN(s.ctx, product.EAN)
		s.NoError(err)
		s.Equal(product, *checkProduct)
	})

	s.Run("insert and update, no error", func() {

	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.EqualError(s.Product.Insert(ctx, &model.Product{}), `context canceled`)
		p, err := s.Product.GetProductByEAN(ctx, "ean")
		s.Nil(p)
		s.EqualError(err, `context canceled`)
	})
}

func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(SqlProductTestSuite))
}
