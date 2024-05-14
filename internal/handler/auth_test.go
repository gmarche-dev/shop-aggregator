package handler_test

import (
	"encoding/json"
	"shop-aggregator/internal/model/request"
	"shop-aggregator/internal/model/response"
)

func (s *HandlerTestSuite) TestLogin() {
	s.Run("no error", func() {
		// create user
		s.createUser("login", "password", "test@test.com")

		user, err := s.HandlerRepositories.Users.GetUserByLogin(s.ctx, "login")
		s.NoError(err)
		s.NotNil(user)

		// login
		bodyLogin := request.Login{
			Login:    "login",
			Password: "password",
		}
		bodyLoginBytes, err := json.Marshal(bodyLogin)
		s.Require().NoError(err)
		wLogin := s.request("POST", "/login", bodyLoginBytes)
		s.Equal(200, wLogin.Code)

		var login response.Login
		s.NoError(json.Unmarshal(wLogin.Body.Bytes(), &login))

		row := s.DB.QueryRow(s.ctx, "SELECT token,is_active from auth WHERE user_id = $1", user.ID)
		var expectedToken string
		var isActive bool
		s.NoError(row.Scan(&expectedToken, &isActive))
		s.Equal(expectedToken, login.Token)
		s.True(isActive)
	})

	s.Run("invalid password", func() {
		// create user
		s.createUser("login2", "password2", "test2@test.com")

		// login
		bodyLogin := request.Login{
			Login:    "login2",
			Password: "nananère",
		}
		bodyLoginBytes, err := json.Marshal(bodyLogin)
		s.Require().NoError(err)
		wLogin := s.request("POST", "/login", bodyLoginBytes)
		s.Equal(400, wLogin.Code)
		s.Equal(`{"error":"invalid password"}`, wLogin.Body.String())
	})

	s.Run("login not exists", func() {
		// login
		bodyLogin := request.Login{
			Login:    "login3",
			Password: "nananère",
		}
		bodyLoginBytes, err := json.Marshal(bodyLogin)
		s.Require().NoError(err)
		wLogin := s.request("POST", "/login", bodyLoginBytes)
		s.Equal(400, wLogin.Code)
		s.Equal(`{"error":"user not found"}`, wLogin.Body.String())
	})
}

func (s *HandlerTestSuite) TestLogOut() {
	s.Run("no error", func() {
		// create user
		token := s.createUserAndGenerateToken("login", "password", "test@test.com")
		user, err := s.HandlerRepositories.Users.GetUserByLogin(s.ctx, "login")
		s.NoError(err)
		s.NotNil(user)

		// logout
		wLogout := s.requestWithToken("POST", "/user/logout", token, nil)

		s.Equal(200, wLogout.Code)
		var isActive bool
		row := s.DB.QueryRow(s.ctx, "SELECT is_active from auth WHERE user_id = $1", user.ID)
		s.NoError(row.Scan(&isActive))
		s.False(isActive)
	})

	s.Run("invalid token", func() {
		// logout
		wLogout := s.requestWithToken("POST", "/user/logout", "notavalidtoken", nil)

		s.Equal(401, wLogout.Code)
		s.Equal(`{"error":"invalid authorization token"}`, wLogout.Body.String())
	})
}
