package model

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Login        string
	Password     string
	HashPassword string
	Email        string
}

type UpdatePassword struct {
	ID          uuid.UUID
	Password    string
	OldPassword string
}

type UpdateEmail struct {
	ID    uuid.UUID
	Email string
}
