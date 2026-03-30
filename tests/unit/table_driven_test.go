package unit_test

import (
	"context"
	"testing"

	"testsystem/internal/model"
	"testsystem/internal/service"
)

func TestCreateUser_TableDriven(t *testing.T) {
	repo := &fakeUserRepo{
		createFn: func(_ context.Context, name string) (model.User, error) {
			return model.User{ID: 1, Name: name}, nil
		},
		listFn: func(_ context.Context) ([]model.User, error) {
			return nil, nil
		},
	}
	svc := service.NewUserService(repo)

	cases := []struct {
		name     string
		input    string
		wantErr  bool
		wantName string
	}{
		{name: "too-short", input: "a", wantErr: true},
		{name: "ok", input: "tom", wantErr: false, wantName: "tom"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := svc.CreateUser(context.Background(), tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if u.Name != tc.wantName {
				t.Fatalf("want name %q, got %q", tc.wantName, u.Name)
			}
		})
	}
}
