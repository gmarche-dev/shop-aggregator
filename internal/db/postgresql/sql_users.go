package postgresql

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"shop-aggregator/internal/model"
)

type User struct {
	db *Client
}

func NewUsers(db *Client) *User {
	return &User{
		db: db,
	}
}

const (
	UpsertUserQuery = `
		INSERT INTO users(login, email, password) 
		VALUES ($1,$2,$3)
		ON CONFLICT(login)
		DO UPDATE SET email = EXCLUDED.email
		RETURNING user_id`
	GetUserByLoginQuery = `SELECT user_id, login, email, password FROM users WHERE login = $1`
	GetUserByIDQuery    = `SELECT user_id, login, email, password FROM users WHERE user_id = $1`
	GetUserByEmailQuery = `SELECT user_id, login, email, password FROM users WHERE email = $1`
	UpdatePasswordQuery = `UPDATE users set password = $2  WHERE user_id = $1`
	UpdateEmailQuery    = `UPDATE users set email = $2  WHERE user_id = $1`
)

func (u *User) Upsert(ctx context.Context, m *model.User) error {
	row := u.db.QueryRow(ctx, UpsertUserQuery, m.Login, m.Email, m.HashPassword)
	err := row.Scan(&m.ID)
	return err
}

func (u *User) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := u.db.QueryRow(ctx, GetUserByIDQuery, id)
	var m model.User
	err := row.Scan(&m.ID, &m.Login, &m.Email, &m.HashPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (u *User) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	row := u.db.QueryRow(ctx, GetUserByLoginQuery, login)
	var m model.User
	err := row.Scan(&m.ID, &m.Login, &m.Email, &m.HashPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	row := u.db.QueryRow(ctx, GetUserByEmailQuery, email)
	var m model.User
	err := row.Scan(&m.ID, &m.Login, &m.Email, &m.HashPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (u *User) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	_, err := u.db.Exec(ctx, UpdatePasswordQuery, id, password)
	return err
}

func (u *User) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	_, err := u.db.Exec(ctx, UpdateEmailQuery, id, email)
	return err
}
