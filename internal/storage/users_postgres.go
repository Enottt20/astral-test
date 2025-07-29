package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/Enottt20/astral-test/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{db: db}
}

func (r *UsersPostgres) Create(ctx context.Context, req domain.RegisterRequest) (string, error) {
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, req.Login, req.Pswd)
	if err != nil {
		return "", err
	}
	return req.Login, err
}

func (r *UsersPostgres) GetByCredentials(ctx context.Context, login, password string) (*domain.User, error) {
	var user domain.User
	query := `SELECT user_id, login, password_hash 
	FROM users 
	WHERE login = $1 AND password_hash = $2`
	err := r.db.GetContext(ctx, &user, query, login, password)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &user, err
}

func (r *UsersPostgres) CreateSession(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO sessions (user_id, token, expires_at)
	VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

func (r *UsersPostgres) DeleteSession(ctx context.Context, token string) error {
	query := `DELETE FROM sessions 
	WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *UsersPostgres) ValidateToken(ctx context.Context, token string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(
	SELECT 1 FROM sessions 
	WHERE token = $1 AND expires_at > NOW())`
	err := r.db.GetContext(ctx, &exists, query, token)
	return exists, err
}

func (r *UsersPostgres) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT user_id, login, password_hash
	FROM users
	WHERE user_id = $1`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &user, nil
}
