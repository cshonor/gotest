package benchmark_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"testsystem/internal/app"
)

func BenchmarkHTTP_Endpoints(b *testing.B) {
	dbPath := filepath.Join(b.TempDir(), "http_endpoints_bench.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		b.Fatalf("open db: %v", err)
	}
	defer db.Close()

	router, err := app.NewRouter(db)
	if err != nil {
		b.Fatalf("new router: %v", err)
	}

	b.Run("GET_/health", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				b.Fatalf("unexpected status: %d", w.Code)
			}
		}
	})
}
