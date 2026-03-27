package repository

import (
	"context"

	"testsystem/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, name string) (model.User, error)
	List(ctx context.Context) ([]model.User, error)
}
