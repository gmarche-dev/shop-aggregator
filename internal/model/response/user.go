package response

import "shop-aggregator/internal/model"

type User struct {
	Login string `json:"login"`
	Email string `json:"email"`
}

func NewUserFromModel(m model.User) User {
	return User{
		Login: m.Login,
		Email: m.Email,
	}
}
