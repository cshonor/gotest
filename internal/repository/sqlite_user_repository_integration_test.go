package repository

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSQLiteUserRepository_CreateAndList_Integration(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo, err := NewSQLiteUserRepository(db)
	if err != nil {
		t.Fatalf("new repo: %v", err)
	}

	ctx := context.Background()
	u1, err := repo.Create(ctx, "alice")
	if err != nil {
		t.Fatalf("create user1: %v", err)
	}
	if u1.ID == 0 {
		t.Fatalf("expected id generated")
	}

	_, err = repo.Create(ctx, "bob")
	if err != nil {
		t.Fatalf("create user2: %v", err)
	}

	users, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Name != "alice" || users[1].Name != "bob" {
		t.Fatalf("unexpected users order/data: %+v", users)
	}
}
