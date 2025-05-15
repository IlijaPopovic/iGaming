package handlers

import (
	"encoding/json"
	"igaming/internal/repository"
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
// @Success 200 {array} object
// @Failure 500 {object} object
// @Router /tournaments [get]
func (h *TournamentHandler) GetTournaments(w http.ResponseWriter, r *http.Request) {

    tournaments, err := h.repo.GetAll(r.Context())
    if err != nil {
        http.Error(w, "Failed to get tournaments", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tournaments)
}