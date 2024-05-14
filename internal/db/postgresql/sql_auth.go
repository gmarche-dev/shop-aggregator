package postgresql

import (
	"context"
	"github.com/google/uuid"
)

type Auth struct {
	db *Client
}

func NewAuth(db *Client) *Auth {
	return &Auth{
		db: db,
	}
}

const (
	LoginQuery = `
		INSERT INTO auth (user_id,token,is_active)
		VALUES($1,$2,true)
		ON CONFLICT(user_id)
		DO UPDATE SET 
		    token = EXCLUDED.token,
		    is_active = true,
		    updated_at = NOW()
	`
	EnsureValidTokenQuery = `
		SELECT u.user_id
		FROM users u
		INNER JOIN auth a ON a.user_id = u.user_id
		WHERE a.token = $1
		AND a.is_active
	`
	LogoutQuery = `UPDATE auth SET is_active = false, updated_at = NOW() WHERE user_id = $1`
)

func (a *Auth) Upsert(ctx context.Context, id uuid.UUID, token string) error {
	_, err := a.db.Exec(ctx, LoginQuery, id, token)
	return err
}

func (a *Auth) EnsureValidToken(ctx context.Context, token string) (uuid.UUID, error) {
	row := a.db.QueryRow(ctx, EnsureValidTokenQuery, token)

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (a *Auth) Logout(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(ctx, LogoutQuery, id)
	return err
}
