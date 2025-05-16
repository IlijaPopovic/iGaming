package handlers

import (
	"encoding/json"
	"igaming/internal/handlers/dtos"
	"igaming/internal/models"
	"igaming/internal/repository"
	"log"
	"net/http"
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
    
    // Handle JSON decoding errors
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
        return
    }

    // Validate request data
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

    // Convert DTO to model
    tournament := models.Tournament{
        Name:      req.Name,
        PrizePool: req.PrizePool,
        StartDate: req.StartDate,
        EndDate:   req.EndDate,
    }

    // Save to repository
    if err := h.repo.Create(r.Context(), &tournament); err != nil {
        log.Printf("Failed to create tournament: %v", err) // Add logging
        respondWithError(w, http.StatusInternalServerError, "Failed to create tournament: "+err.Error())
        return
    }
    
    // Convert back to response DTO
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