package e2e_test

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"testsystem/internal/app"
)

func TestHealth_E2E(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "health.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	router, err := app.NewRouter(db)
	if err != nil {
		t.Fatalf("new router: %v", err)
	}
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("get /health: %v", err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d, body=%q", resp.StatusCode, string(b))
	}
	if string(b) != "ok" {
		t.Fatalf("want body %q, got %q", "ok", string(b))
	}
}
