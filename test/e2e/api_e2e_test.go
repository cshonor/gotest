package e2e

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"testsystem/internal/app"
	"testsystem/internal/model"
)

func TestAPI_CreateThenListUsers_E2E(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "e2e.db")
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

	body := []byte(`{"name":"charlie"}`)
	resp, err := http.Post(server.URL+"/api/v1/users", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post user: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("want 201, got %d", resp.StatusCode)
	}

	var created model.User
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created user: %v", err)
	}
	if created.Name != "charlie" || created.ID == 0 {
		t.Fatalf("unexpected created user: %+v", created)
	}

	listResp, err := http.Get(server.URL + "/api/v1/users")
	if err != nil {
		t.Fatalf("get users: %v", err)
	}
	defer listResp.Body.Close()
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", listResp.StatusCode)
	}

	var users []model.User
	if err := json.NewDecoder(listResp.Body).Decode(&users); err != nil {
		t.Fatalf("decode users: %v", err)
	}
	if len(users) != 1 || users[0].Name != "charlie" {
		t.Fatalf("unexpected users: %+v", users)
	}
}
