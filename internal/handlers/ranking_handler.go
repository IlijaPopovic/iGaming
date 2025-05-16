package handlers

import (
	"igaming/internal/repository"
	"net/http"
)

type RankingHandler struct {
    repo *repository.PlayerRepository
}

func NewRankingHandler(repo *repository.PlayerRepository) *RankingHandler {
    return &RankingHandler{repo: repo}
}

// GetPlayerRankings godoc
// @Summary Get player rankings
// @Description Get ranked list of players by account balance
// @Tags rankings
// @Accept  json
// @Produce  json
// @Success 200 {array} models.PlayerRanking
// @Failure 500 {object} ErrorResponse
// @Router /rankings [get]
func (h *RankingHandler) GetPlayerRankings(w http.ResponseWriter, r *http.Request) {
    rankings, err := h.repo.GetRankings(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to get rankings")
        return
    }

    respondWithJSON(w, http.StatusOK, rankings)
}