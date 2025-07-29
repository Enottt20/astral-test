package service

import (
	"context"

	"github.com/Enottt20/astral-test/internal/domain"
	"github.com/Enottt20/astral-test/internal/storage"
)

type UserService struct {
	repo *storage.Repository
}

func NewUserService(repo *storage.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return s.repo.Users.GetByID(ctx, id)
}
