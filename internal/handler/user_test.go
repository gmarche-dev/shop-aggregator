package handler_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"shop-aggregator/internal/model/request"
)

func (s *HandlerTestSuite) TestCreateUser() {
	s.Run("no error", func() {
		bodyCreate := request.CreateUser{
			Login:    "login",
			Password: "password",
			Email:    "test@test.com",
		}
		bodyCreateBytes, err := json.Marshal(bodyCreate)
		s.Require().NoError(err)

		wCreate := s.request("POST", "/create-user", bodyCreateBytes)

		s.Equal(200, wCreate.Code)
		s.Equal(`{"message":"User created"}`, wCreate.Body.String())
		user, err := s.HandlerRepositories.Users.GetUserByLogin(s.ctx, "login")
		s.NoError(err)
		s.NotNil(user)
		s.Equal("login", user.Login)
		s.Equal("test@test.com", user.Email)
	})

	s.Run("email already exists, return an error", func() {
		s.createUser("login2", "password", "test2@test.com")
		bodyCreate := request.CreateUser{
			Login:    "login22",
			Password: "password",
			Email:    "test2@test.com",
		}
		bodyCreateBytes, err := json.Marshal(bodyCreate)
		s.Require().NoError(err)

		wCreate := s.request("POST", "/create-user", bodyCreateBytes)

		s.Equal(400, wCreate.Code)
		s.Equal(`{"error":"user exist for email test2@test.com"}`, wCreate.Body.String())
	})

	s.Run("login already exists, return an error", func() {
		s.createUser("login3", "password", "test3@test.com")
		bodyCreate := request.CreateUser{
			Login:    "login3",
			Password: "password",
			Email:    "test33@test.com",
		}
		bodyCreateBytes, err := json.Marshal(bodyCreate)
		s.Require().NoError(err)

		wCreate := s.request("POST", "/create-user", bodyCreateBytes)

		s.Equal(400, wCreate.Code)
		s.Equal(`{"error":"user exist for login login3"}`, wCreate.Body.String())
	})

	s.Run("no login, return an error", func() {
		bodyCreate := request.CreateUser{
			Login:    "",
			Password: "password",
			Email:    "test@test.com",
		}
		bodyCreateBytes, err := json.Marshal(bodyCreate)
		s.Require().NoError(err)

		wCreate := s.request("POST", "/create-user", bodyCreateBytes)

		s.Equal(400, wCreate.Code)
		s.Equal(`{"error":"Key: 'CreateUser.Login' Error:Field validation for 'Login' failed on the 'required' tag"}`, wCreate.Body.String())
	})

	s.Run("no password, return an error", func() {
		bodyCreate := request.CreateUser{
			Login:    "Login",
			Password: "",
			Email:    "test@test.com",
		}
		bodyCreateBytes, err := json.Marshal(bodyCreate)
		s.Require().NoError(err)

		wCreate := s.request("POST", "/create-user", bodyCreateBytes)

		s.Equal(400, wCreate.Code)
		s.Equal(`{"error":"Key: 'CreateUser.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`, wCreate.Body.String())
	})

	s.Run("no email, return an error", func() {
		bodyCreate := request.CreateUser{
			Login:    "login",
			Password: "password",
			Email:    "",
		}
		bodyCreateBytes, err := json.Marshal(bodyCreate)
		s.Require().NoError(err)

		wCreate := s.request("POST", "/create-user", bodyCreateBytes)

		s.Equal(400, wCreate.Code)
		s.Equal(`{"error":"Key: 'CreateUser.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`, wCreate.Body.String())
	})
}

func (s *HandlerTestSuite) TestResetPassword() {
	s.Run("no error", func() {
		// create user and login
		token := s.createUserAndGenerateToken("login", "password", "email")

		// update password
		newPassword := "newpassword"
		bodyUpdatePassword := request.UpdatePassword{
			OldPassword: "password",
			Password:    newPassword,
		}
		bodyUpdatePasswordBytes, err := json.Marshal(bodyUpdatePassword)
		s.Require().NoError(err)
		wUpdatePassword := s.requestWithToken("POST", "/user/reset-password", token, bodyUpdatePasswordBytes)
		s.T().Log(wUpdatePassword.Body.String())
		s.Equal(200, wUpdatePassword.Code)

		// logout
		s.logout(token)

		// login with new password
		newToken := s.login("login", newPassword)
		userID, err := s.HandlerRepositories.Auth.EnsureValidToken(s.ctx, newToken)
		s.NoError(err)
		s.NotEqual(uuid.Nil, userID)
	})

	s.Run("old password is missing", func() {
		token := s.createUserAndGenerateToken("login2", "password2", "email2")

		newPassword := "newpassword"
		bodyUpdatePassword := request.UpdatePassword{
			OldPassword: "",
			Password:    newPassword,
		}
		bodyUpdatePasswordBytes, err := json.Marshal(bodyUpdatePassword)
		s.Require().NoError(err)
		wUpdatePassword := s.requestWithToken("POST", "/user/reset-password", token, bodyUpdatePasswordBytes)
		s.T().Log(wUpdatePassword.Body.String())
		s.Equal(400, wUpdatePassword.Code)
		s.Equal(`{"error":"Key: 'UpdatePassword.OldPassword' Error:Field validation for 'OldPassword' failed on the 'required' tag"}`, wUpdatePassword.Body.String())
	})

	s.Run("new password is missing", func() {
		token := s.createUserAndGenerateToken("login3", "password3", "email3")

		newPassword := ""
		bodyUpdatePassword := request.UpdatePassword{
			OldPassword: "password3",
			Password:    newPassword,
		}
		bodyUpdatePasswordBytes, err := json.Marshal(bodyUpdatePassword)
		s.Require().NoError(err)
		wUpdatePassword := s.requestWithToken("POST", "/user/reset-password", token, bodyUpdatePasswordBytes)
		s.T().Log(wUpdatePassword.Body.String())
		s.Equal(400, wUpdatePassword.Code)
		s.Equal(`{"error":"Key: 'UpdatePassword.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`, wUpdatePassword.Body.String())
	})

	s.Run("old password is wrong", func() {
		token := s.createUserAndGenerateToken("login4", "password4", "email4")

		newPassword := "newPassword4"
		bodyUpdatePassword := request.UpdatePassword{
			OldPassword: "password3",
			Password:    newPassword,
		}
		bodyUpdatePasswordBytes, err := json.Marshal(bodyUpdatePassword)
		s.Require().NoError(err)
		wUpdatePassword := s.requestWithToken("POST", "/user/reset-password", token, bodyUpdatePasswordBytes)
		s.T().Log(wUpdatePassword.Body.String())
		s.Equal(400, wUpdatePassword.Code)
		s.Equal(`{"error":"invalid old password"}`, wUpdatePassword.Body.String())
	})
}
