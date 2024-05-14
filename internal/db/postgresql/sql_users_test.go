package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"shop-aggregator/internal/model"
	"testing"
)

type SqlUserTestSuite struct {
	DBTestSuite
	User *User
	user model.User
}

func (s *SqlUserTestSuite) SetupTest() {
	s.user = model.User{
		Login:        "user42",
		Email:        "user42@test.com",
		HashPassword: "iamuser42forever",
	}
	s.User = NewUsers(s.DB)
}

func (s *SqlUserTestSuite) TearDownTest() {
	_, err := s.DB.Exec(s.ctx, "TRUNCATE TABLE users")
	s.Require().NoError(err)
}

func (s *SqlUserTestSuite) TestUpsert() {
	s.Run("no error", func() {
		// insert user
		user := s.user
		s.Require().NoError(s.User.Upsert(s.ctx, &user))
		s.NotEqual(uuid.Nil, user.ID)

		// check user_id is valid
		checkUser, err := s.User.GetUserByID(s.ctx, user.ID)
		s.Require().NoError(err)
		s.Equal(user, *checkUser)

		// update email
		user.Email = "newemail@test.com"
		s.Require().NoError(s.User.Upsert(s.ctx, &user))
		s.NotEqual(uuid.Nil, user.ID)

		// check email is update
		checkUser, err = s.User.GetUserByID(s.ctx, user.ID)
		s.Require().NoError(err)
		s.Equal(user, *checkUser)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.EqualError(s.User.Upsert(ctx, &model.User{}), `context canceled`)
	})
}

func (s *SqlUserTestSuite) TestGetUserByID() {
	s.Run("no error", func() {
		user := s.user
		u, err := s.User.GetUserByID(s.ctx, uuid.New())
		s.Nil(u)
		s.NoError(err)

		s.Require().NoError(s.User.Upsert(s.ctx, &user))
		s.NotEqual(uuid.Nil, user.ID)

		checkUser, err := s.User.GetUserByID(s.ctx, user.ID)
		s.Require().NoError(err)
		s.Equal(user, *checkUser)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		u, err := s.User.GetUserByID(ctx, uuid.New())
		s.Nil(u)
		s.EqualError(err, `context canceled`)
	})
}

func (s *SqlUserTestSuite) TestGetUserByLogin() {
	s.Run("no error", func() {
		user := s.user
		u, err := s.User.GetUserByLogin(s.ctx, user.Login)
		s.Nil(u)
		s.NoError(err)

		s.Require().NoError(s.User.Upsert(s.ctx, &user))
		s.NotEqual(uuid.Nil, user.ID)

		checkUser, err := s.User.GetUserByLogin(s.ctx, user.Login)
		s.Require().NoError(err)
		s.Equal(user, *checkUser)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		u, err := s.User.GetUserByLogin(ctx, "niet")
		s.Nil(u)
		s.EqualError(err, `context canceled`)
	})
}

func (s *SqlUserTestSuite) TestGetUserByEmail() {
	s.Run("no error", func() {
		user := s.user
		u, err := s.User.GetUserByEmail(s.ctx, user.Email)
		s.Nil(u)
		s.NoError(err)

		s.Require().NoError(s.User.Upsert(s.ctx, &user))
		s.NotEqual(uuid.Nil, user.ID)

		checkUser, err := s.User.GetUserByEmail(s.ctx, user.Email)
		s.Require().NoError(err)
		s.Equal(user, *checkUser)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		u, err := s.User.GetUserByEmail(ctx, "niet")
		s.Nil(u)
		s.EqualError(err, `context canceled`)
	})
}

func (s *SqlUserTestSuite) TestUpdatePassword() {
	s.Run("no error", func() {
		user := s.user
		s.Require().NoError(s.User.Upsert(s.ctx, &user))
		s.NotEqual(uuid.Nil, user.ID)

		user.HashPassword = "thisismynewpassword"
		s.Require().NoError(s.User.UpdatePassword(s.ctx, user.ID, user.HashPassword))

		checkUser, err := s.User.GetUserByEmail(s.ctx, user.Email)
		s.Require().NoError(err)
		s.Equal(user, *checkUser)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.EqualError(s.User.UpdatePassword(ctx, uuid.New(), "niet"), `context canceled`)
	})
}

func (s *SqlUserTestSuite) TestUpdateEmail() {
	s.Run("no error", func() {
		user := s.user
		s.Require().NoError(s.User.Upsert(s.ctx, &user))
		s.NotEqual(uuid.Nil, user.ID)

		user.Email = "new@test.email"
		s.Require().NoError(s.User.UpdateEmail(s.ctx, user.ID, user.Email))

		checkUser, err := s.User.GetUserByEmail(s.ctx, user.Email)
		s.Require().NoError(err)
		s.Equal(user, *checkUser)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.EqualError(s.User.UpdatePassword(ctx, uuid.New(), "niet"), `context canceled`)
	})
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(SqlUserTestSuite))
}
