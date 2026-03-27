package service

import (
	"context"
	"errors"
	"strings"

	"testsystem/internal/model"
	"testsystem/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, name string) (model.User, error) {
	if len(strings.TrimSpace(name)) < 2 {
		return model.User{}, errors.New("name must be at least 2 chars")
	}
	return s.repo.Create(ctx, name)
}

func (s *UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}
