package integration_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"testsystem/internal/repository"
)

func TestSQLite_MigrationAndEmptyList(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "migration.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo, err := repository.NewSQLiteUserRepository(db)
	if err != nil {
		t.Fatalf("new repo (runs migration): %v", err)
	}

	users, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected empty users, got %d", len(users))
	}
}
