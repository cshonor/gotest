package app

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"testsystem/internal/handler"
	"testsystem/internal/repository"
	"testsystem/internal/service"
)

func NewRouter(db *sql.DB) (http.Handler, error) {
	userRepo, err := repository.NewSQLiteUserRepository(db)
	if err != nil {
		return nil, err
	}
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", userHandler.CreateUser)
		r.Get("/users", userHandler.ListUsers)
	})
	return r, nil
}
