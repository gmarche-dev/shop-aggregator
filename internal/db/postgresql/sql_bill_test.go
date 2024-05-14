package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"shop-aggregator/internal/model"
	"testing"
)

type SqlBillTestSuite struct {
	DBTestSuite
	Bill *Bill
}

func (s *SqlBillTestSuite) SetupTest() {
	s.Bill = NewBill(s.DB)
}

func (s *SqlBillTestSuite) TearDownTest() {
	_, err := s.DB.Exec(s.ctx, "TRUNCATE TABLE bill")
	s.Require().NoError(err)
}

func (s *SqlBillTestSuite) TestBill() {
	s.Run("insert and get, no error", func() {
		userID := uuid.New()
		bills := []*model.Bill{
			{
				UserID:  userID,
				StoreID: uuid.New(),
				Amount:  "2",
				State:   model.BillStateCreate,
			},
			{
				UserID:  userID,
				StoreID: uuid.New(),
				Amount:  "5",
				State:   model.BillStateCreate,
			},
			{
				UserID:  userID,
				StoreID: uuid.New(),
				Amount:  "245",
				State:   model.BillStateCreate,
			},
			{
				UserID:  userID,
				StoreID: uuid.New(),
				Amount:  "72",
				State:   model.BillStateCreate,
			},
			{
				UserID:  userID,
				StoreID: uuid.New(),
				Amount:  "245",
				State:   model.BillStateCreate,
			},
		}

		for _, bill := range bills {
			s.Require().NoError(s.Bill.Insert(s.ctx, bill))
		}

		checkBills, err := s.Bill.GetBillsByUserID(s.ctx, userID)
		s.Require().NoError(err)
		s.ElementsMatch(bills, checkBills)

		// insert bills with a new userID
		newUserID := uuid.New()
		newBills := []*model.Bill{
			{
				UserID:  newUserID,
				StoreID: uuid.New(),
				Amount:  "2",
				State:   model.BillStateCreate,
			},
			{
				UserID:  newUserID,
				StoreID: uuid.New(),
				Amount:  "5",
				State:   model.BillStateCreate,
			},
			{
				UserID:  newUserID,
				StoreID: uuid.New(),
				Amount:  "245",
				State:   model.BillStateCreate,
			},
		}

		for _, bill := range newBills {
			s.Require().NoError(s.Bill.Insert(s.ctx, bill))
		}

		checkNewBills, err := s.Bill.GetBillsByUserID(s.ctx, newUserID)
		s.Require().NoError(err)
		s.ElementsMatch(newBills, checkNewBills)
	})

	s.Run("insert and update, no error", func() {
		bill := model.Bill{
			UserID:  uuid.New(),
			StoreID: uuid.New(),
			Amount:  "1234",
			State:   model.BillStateCreate,
		}

		// insert
		s.NoError(s.Bill.Insert(s.ctx, &bill))

		// check well insert
		checkBills, err := s.Bill.GetBillsByUserID(s.ctx, bill.UserID)
		s.NoError(err)
		s.Len(checkBills, 1)
		s.Equal(bill, *checkBills[0])

		// update
		bill.Amount = "987"
		s.NoError(s.Bill.Update(s.ctx, &bill))

		// check update
		checkBills, err = s.Bill.GetBillsByUserID(s.ctx, bill.UserID)
		s.NoError(err)
		s.Len(checkBills, 1)
		s.Equal(bill, *checkBills[0])
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.EqualError(s.Bill.Insert(ctx, &model.Bill{}), `context canceled`)
		s.EqualError(s.Bill.Update(ctx, &model.Bill{}), `context canceled`)
		b, err := s.Bill.GetBillsByUserID(ctx, uuid.New())
		s.Nil(b)
		s.EqualError(err, `context canceled`)
	})
}

func TestBillTestSuite(t *testing.T) {
	suite.Run(t, new(SqlBillTestSuite))
}
