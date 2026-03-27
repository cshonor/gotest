package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"testsystem/internal/app"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "file:app.db?cache=shared&mode=rwc"
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("open db error: %v", err)
	}
	defer db.Close()

	router, err := app.NewRouter(db)
	if err != nil {
		log.Fatalf("init app error: %v", err)
	}

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
