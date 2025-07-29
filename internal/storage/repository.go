package storage

import (
	"context"
	"time"

	"github.com/Enottt20/astral-test/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Documents interface {
	Create(ctx context.Context, token string, doc *domain.Document, jsonData string, fileData []byte) error
	GetAll(ctx context.Context, token, login, key, value string, limit int) ([]*domain.Document, error)
	GetByID(ctx context.Context, token, id string) (*domain.Document, []byte, error)
	GetFileData(ctx context.Context, token, id string) ([]byte, error)
	Delete(ctx context.Context, token, id string) error
}

type Users interface {
	Create(ctx context.Context, req domain.RegisterRequest) (string, error)
	GetByCredentials(ctx context.Context, login, password string) (*domain.User, error)
	CreateSession(ctx context.Context, userID int, token string, expiresAt time.Time) error
	DeleteSession(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (bool, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
}

type Repository struct {
	Documents
	Users
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Documents: NewDocumentsPostgres(db),
		Users:     NewUsersPostgres(db),
	}
}
