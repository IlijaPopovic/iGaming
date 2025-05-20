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

	playerRepo := repository.NewPlayerRepository(db)
    playerHandler := handlers.NewPlayerHandler(playerRepo)
	rankingHandler := handlers.NewRankingHandler(playerRepo)

	betRepo := repository.NewTournamentBetRepository(db, playerRepo, tournamentRepo)
	betHandler := handlers.NewTournamentBetHandler(betRepo)

	// ______>
	
	router.Get("/tournaments", tournamentHandler.GetTournaments)
	router.Post("/tournaments", tournamentHandler.CreateTournament)

	router.Post("/tournaments/prizes/{id}", tournamentHandler.DistributePrizes)

	router.Get("/players", playerHandler.GetPlayers)
	router.Post("/players", playerHandler.CreatePlayer)

	router.Get("/bets", betHandler.GetBets)
	router.Post("/bets", betHandler.CreateBet)

	router.Get("/rankings", rankingHandler.GetPlayerRankings)

	

	return router
}