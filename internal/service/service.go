package service

import (
	"context"
	"mime/multipart"

	"github.com/Enottt20/astral-test/internal/domain"
	"github.com/Enottt20/astral-test/internal/storage"
	"github.com/redis/go-redis/v9"
)

type Auth interface {
	Register(ctx context.Context, req domain.RegisterRequest) (string, error)
	Authenticate(ctx context.Context, req domain.AuthRequest) (string, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
	Logout(ctx context.Context, token string) error
}

type Users interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
}

type Documents interface {
	Upload(ctx context.Context, token, meta, jsonData string, file multipart.File, header *multipart.FileHeader, isFileLoaded bool) (*domain.Document, error)
	GetAll(ctx context.Context, token, login, key, value string, limit int) ([]*domain.Document, error)
	GetByID(ctx context.Context, token, id string) (*domain.Document, []byte, error)
	Delete(ctx context.Context, token, id string) error
}

type Service struct {
	Auth
	Users
	Documents
}

func NewService(repo *storage.Repository, adminToken string, cache *redis.Client) *Service {
	return &Service{
		Auth:      NewAuthService(repo, adminToken),
		Users:     NewUserService(repo),
		Documents: NewDocumentService(repo, cache),
	}
}
