package benchmark_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"testsystem/internal/app"
	"testsystem/internal/repository"
	"testsystem/internal/service"
)

func BenchmarkUserService_CreateUser(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "svc_bench.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		b.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo, err := repository.NewSQLiteUserRepository(db)
	if err != nil {
		b.Fatalf("new repo: %v", err)
	}
	svc := service.NewUserService(repo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := svc.CreateUser(context.Background(), "bench-user"); err != nil {
			b.Fatalf("create user: %v", err)
		}
	}
}

func BenchmarkRepository_ListUsers(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "repo_bench.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		b.Fatalf("open db: %v", err)
	}
	defer db.Close()
	repo, err := repository.NewSQLiteUserRepository(db)
	if err != nil {
		b.Fatalf("new repo: %v", err)
	}
	for i := 0; i < 1000; i++ {
		if _, err := repo.Create(context.Background(), "seed-user"); err != nil {
			b.Fatalf("seed user: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.List(context.Background()); err != nil {
			b.Fatalf("list users: %v", err)
		}
	}
}

func BenchmarkHTTP_CreateUser(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "http_bench.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		b.Fatalf("open db: %v", err)
	}
	defer db.Close()
	router, err := app.NewRouter(db)
	if err != nil {
		b.Fatalf("new router: %v", err)
	}

	reqBody := `{"name":"perf-user"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			b.Fatalf("unexpected status: %d", w.Code)
		}
	}
}
