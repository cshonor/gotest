package full_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"testsystem/internal/app"
)

// “全量测试”在 Go 里通常就是 `go test ./...`。
// 这个目录放一个最小 smoke test，帮助新手理解：如何一条命令跑起整个项目的关键链路。
func TestSmoke_RouterAndHealth(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "smoke.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	router, err := app.NewRouter(db)
	if err != nil {
		t.Fatalf("new router: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
}
