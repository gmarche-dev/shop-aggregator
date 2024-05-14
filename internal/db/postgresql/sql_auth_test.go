package postgresql

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"shop-aggregator/internal/model"
	"shop-aggregator/internal/utils"
	"testing"
)

type SqlAuthTestSuite struct {
	DBTestSuite
	Auth *Auth
	User *User
	user model.User
}

func (s *SqlAuthTestSuite) SetupTest() {
	s.user = model.User{
		Login:        "user42",
		Email:        "user42@test.com",
		HashPassword: "iamuser42forever",
	}
	s.Auth = NewAuth(s.DB)
	s.User = NewUsers(s.DB)
}

func (s *SqlAuthTestSuite) TearDownTest() {
	_, err := s.DB.Exec(s.ctx, "TRUNCATE TABLE user")
	_, err = s.DB.Exec(s.ctx, "TRUNCATE TABLE auth")
	s.Require().NoError(err)
}

func (s *SqlAuthTestSuite) getIDAndTokenByToken(t string) (uuid.UUID, string, bool) {
	var id uuid.UUID
	var token string
	var isActive bool
	row := s.DB.QueryRow(s.ctx, "SELECT user_id, token, is_active FROM auth WHERE token = $1", t)
	err := row.Scan(&id, &token, &isActive)
	s.Require().NoError(err)
	return id, token, isActive
}

func (s *SqlAuthTestSuite) TestAuth() {
	s.Run("no error", func() {
		user := s.user
		s.Require().NoError(s.User.Upsert(s.ctx, &user))
		s.NotEqual(uuid.Nil, user.ID)

		// Generate a new user ID and token
		userID := user.ID
		userToken, err := utils.GenerateToken(128)
		s.NoError(err)

		// Insert the user ID and token into the auth table
		s.NoError(s.Auth.Upsert(s.ctx, userID, userToken))

		// Retrieve the inserted user ID, token, and validity from the auth table
		id, token, isActive := s.getIDAndTokenByToken(userToken)

		// Assert that the retrieved values match the inserted values
		s.Equal(userID, id)
		s.Equal(userToken, token)
		s.True(isActive)

		// EnsureValidToken should return the inserted user ID without error
		id, err = s.Auth.EnsureValidToken(s.ctx, userToken)
		s.Equal(userID, id)
		s.NoError(err)

		// Logout the user by setting the is_valid column to false
		s.NoError(s.Auth.Logout(s.ctx, userID))

		// Retrieve the user ID, token, and validity from the auth table after logout
		id, token, isActive = s.getIDAndTokenByToken(userToken)

		// Assert that the validity is now false
		s.False(isActive)

		// valid again a row
		userToken, err = utils.GenerateToken(128)
		s.NoError(err)

		// Insert the user ID and token into the auth table
		s.NoError(s.Auth.Upsert(s.ctx, userID, userToken))

		// Retrieve the inserted user ID, token, and validity from the auth table
		id, token, isActive = s.getIDAndTokenByToken(userToken)
		s.Equal(userID, id)
		s.Equal(userToken, token)
		s.True(isActive)
	})

	s.Run("context cancel error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.EqualError(s.Auth.Upsert(ctx, uuid.New(), "token"), `context canceled`)
		s.EqualError(s.Auth.Logout(ctx, uuid.New()), `context canceled`)
		id, err := s.Auth.EnsureValidToken(ctx, "token")
		s.Equal(uuid.Nil, id)
		s.EqualError(err, `context canceled`)
	})
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(SqlAuthTestSuite))
}
