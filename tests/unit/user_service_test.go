package unit_test

import (
	"context"
	"errors"
	"testing"

	"testsystem/internal/model"
	"testsystem/internal/service"
)

type fakeUserRepo struct {
	createFn func(ctx context.Context, name string) (model.User, error)
	listFn   func(ctx context.Context) ([]model.User, error)
}

func (f *fakeUserRepo) Create(ctx context.Context, name string) (model.User, error) {
	return f.createFn(ctx, name)
}

func (f *fakeUserRepo) List(ctx context.Context) ([]model.User, error) {
	return f.listFn(ctx)
}

func TestCreateUser_ValidateAndCallRepo(t *testing.T) {
	repo := &fakeUserRepo{
		createFn: func(_ context.Context, name string) (model.User, error) {
			return model.User{ID: 1, Name: name}, nil
		},
		listFn: func(_ context.Context) ([]model.User, error) {
			return nil, nil
		},
	}
	svc := service.NewUserService(repo)

	if _, err := svc.CreateUser(context.Background(), "a"); err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	got, err := svc.CreateUser(context.Background(), "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "alice" || got.ID != 1 {
		t.Fatalf("unexpected user: %+v", got)
	}
}

func TestListUsers_RepoError(t *testing.T) {
	wantErr := errors.New("repo down")
	repo := &fakeUserRepo{
		createFn: func(_ context.Context, _ string) (model.User, error) {
			return model.User{}, nil
		},
		listFn: func(_ context.Context) ([]model.User, error) {
			return nil, wantErr
		},
	}
	svc := service.NewUserService(repo)
	_, err := svc.ListUsers(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("want %v, got %v", wantErr, err)
	}
}
