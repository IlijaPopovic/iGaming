package handlers

import (
	"encoding/json"
	"errors"
	"igaming/internal/handlers/dtos"
	"igaming/internal/models"
	"igaming/internal/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"
)

type TournamentHandler struct {
    repo *repository.TournamentRepository
}

func NewTournamentHandler(repo *repository.TournamentRepository) *TournamentHandler {
    return &TournamentHandler{repo: repo}
}

// GetTournaments godoc
// @Summary Get all tournaments
// @Description Get list of all tournaments
// @Tags tournaments
// @Accept json
// @Produce json
// @Success 200 {array} models.Tournament
// @Failure 500 {object} handlers.ErrorResponse
// @Router /tournaments [get]
func (h *TournamentHandler) GetTournaments(w http.ResponseWriter, r *http.Request) {

    tournaments, err := h.repo.GetAllTournaments(r.Context())
    if err != nil {
        http.Error(w, "Failed to get tournaments", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tournaments)
}

// CreateTournament godoc
// @Summary Create a new tournament
// @Description Creates a new tournament with the provided details
// @Tags tournaments
// @Accept  json
// @Produce  json
// @Param   request body     dtos.CreateTournamentRequest  true  "Tournament Creation Data"
// @Success 201     {object} dtos.TournamentResponse
// @Failure 400     {object} handlers.ErrorResponse
// @Failure 500     {object} handlers.ErrorResponse
// @Router /tournaments [post]
func (h *TournamentHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
    var req dtos.CreateTournamentRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
        return
    }

    if req.Name == "" {
        respondWithError(w, http.StatusBadRequest, "Name is required")
        return
    }
    
    if req.PrizePool <= 0 {
        respondWithError(w, http.StatusBadRequest, "Prize pool must be positive")
        return
    }
    
    if req.EndDate.Before(req.StartDate) {
        respondWithError(w, http.StatusBadRequest, "End date must be after start date")
        return
    }

    tournament := models.Tournament{
        Name:      req.Name,
        PrizePool: req.PrizePool,
        StartDate: req.StartDate,
        EndDate:   req.EndDate,
    }

    if err := h.repo.Create(r.Context(), &tournament); err != nil {
        log.Printf("Failed to create tournament: %v", err) // Add logging
        respondWithError(w, http.StatusInternalServerError, "Failed to create tournament: "+err.Error())
        return
    }
    
    response := dtos.TournamentResponse{
        ID:        tournament.ID,
        Name:      tournament.Name,
        PrizePool: tournament.PrizePool,
        StartDate: tournament.StartDate,
        EndDate:   tournament.EndDate,
        CreatedAt: tournament.CreatedAt,
    }
    
    respondWithJSON(w, http.StatusCreated, response)
}

// Change this, this is not
func respondWithError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(payload)
}

// DistributePrizes godoc
// @Summary Distribute tournament prizes
// @Description Calculate and distribute prizes for a completed tournament
// @Tags tournaments
// @Accept  json
// @Produce  json
// @Param   id path int true "Tournament ID"
// @Success 202 {object} map[string]interface{} "message: Prizes distributed successfully"
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /tournaments/{id}/prizes [post]
func (h *TournamentHandler) DistributePrizes(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    tournamentID, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid tournament ID")
        return
    }

    exists, err := h.repo.Exists(r.Context(), uint(tournamentID))
    if err != nil || !exists {
        respondWithError(w, http.StatusNotFound, "Tournament not found")
        return
    }

    // Execute distribution
    if err := h.repo.DistributePrizes(r.Context(), uint(tournamentID)); err != nil {
        log.Printf("Prize distribution error: %v", err)
        
        var mysqlErr *mysql.MySQLError
        if errors.As(err, &mysqlErr) {
            switch mysqlErr.Number {
            case 1329: // No data found
                respondWithError(w, http.StatusBadRequest, "No eligible bets for tournament")
                return
            case 1365: // Division by zero
                respondWithError(w, http.StatusConflict, "Cannot distribute prizes - invalid participant count")
                return
            }
        }

        respondWithError(w, http.StatusInternalServerError, "Prize distribution failed")
        return
    }

    respondWithJSON(w, http.StatusAccepted, map[string]interface{}{
        "message": "Prizes distributed successfully",
        "tournament_id": tournamentID,
    })
}