package server

import (
	"database/sql"
	"igaming/internal/handlers"
	"igaming/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "igaming/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(db *sql.DB) http.Handler {
	router := chi.NewRouter()

	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

    tournamentRepo := repository.NewTournamentRepository(db)
    tournamentHandler := handlers.NewTournamentHandler(tournamentRepo)
	
	router.Get("/tournaments", tournamentHandler.GetTournaments)
	
	return router
}