package handlers

import (
	"encoding/json"
	"igaming/internal/handlers/dtos"
	"igaming/internal/models"
	"igaming/internal/repository"
	"net/http"
)

type TournamentBetHandler struct {
	repo *repository.TournamentBetRepository
}

func NewTournamentBetHandler(repo *repository.TournamentBetRepository) *TournamentBetHandler {
	return &TournamentBetHandler{repo: repo}
}

// CreateBet godoc
// @Summary Place a new bet
// @Description Place a wager on a tournament
// @Tags bets
// @Accept json
// @Produce json
// @Param request body dtos.CreateTournamentBetRequest true "Bet details"
// @Success 201 {object} dtos.TournamentBetResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bets [post]
func (h *TournamentBetHandler) CreateBet(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateTournamentBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	bet := models.TournamentBet{
		PlayerID:     req.PlayerID,
		TournamentID: req.TournamentID,
		BetAmount:    req.BetAmount,
	}

	if err := h.repo.Create(r.Context(), &bet); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "player does not exist" || err.Error() == "tournament does not exist" {
			status = http.StatusNotFound
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, dtos.TournamentBetResponse{
		ID:           bet.ID,
		PlayerID:     bet.PlayerID,
		TournamentID: bet.TournamentID,
		BetAmount:    bet.BetAmount,
		CreatedAt:    bet.CreatedAt,
	})
}

// GetBets godoc
// @Summary Get all bets
// @Description Retrieve list of all placed bets
// @Tags bets
// @Accept json
// @Produce json
// @Success 200 {array} dtos.TournamentBetResponse
// @Failure 500 {object} ErrorResponse
// @Router /bets [get]
func (h *TournamentBetHandler) GetBets(w http.ResponseWriter, r *http.Request) {
	bets, err := h.repo.GetAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve bets: "+err.Error())
		return
	}

	response := make([]dtos.TournamentBetResponse, 0, len(bets))
	for _, bet := range bets {
		response = append(response, dtos.TournamentBetResponse{
			ID:           bet.ID,
			PlayerID:     bet.PlayerID,
			TournamentID: bet.TournamentID,
			BetAmount:    bet.BetAmount,
			CreatedAt:    bet.CreatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}