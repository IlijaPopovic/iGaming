package handlers

import (
	"encoding/json"
	"igaming/internal/handlers/dtos"
	"igaming/internal/models"
	"igaming/internal/repository"
	"net/http"
)

type PlayerHandler struct {
	repo *repository.PlayerRepository
}

func NewPlayerHandler(repo *repository.PlayerRepository) *PlayerHandler {
	return &PlayerHandler{repo: repo}
}

// CreatePlayer godoc
// @Summary Create a new player
// @Description Register a new player account
// @Tags players
// @Accept json
// @Produce json
// @Param request body dtos.CreatePlayerRequest true "Player registration data"
// @Success 201 {object} dtos.PlayerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /players [post]
func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	player := models.Player{
		Name:           req.Name,
		Email:          req.Email,
		PasswordHash:   req.Password, // >>>>> Should be hashed(LATER)
		AccountBalance: req.AccountBalance,
	}

	if err := h.repo.Create(r.Context(), &player); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create player: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, dtos.PlayerResponse{
		ID:            player.ID,
		Name:          player.Name,
		Email:         player.Email,
		AccountBalance: player.AccountBalance,
		CreatedAt:     player.CreatedAt,
		UpdatedAt:     player.UpdatedAt,
	})
}

// GetPlayers godoc
// @Summary Get all players
// @Description Retrieve list of all registered players
// @Tags players
// @Accept json
// @Produce json
// @Success 200 {array} dtos.PlayerResponse
// @Failure 500 {object} ErrorResponse
// @Router /players [get]
func (h *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := h.repo.GetAllPlayers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve players: "+err.Error())
		return
	}

	response := make([]dtos.PlayerResponse, 0, len(players))
	for _, p := range players {
		response = append(response, dtos.PlayerResponse{
			ID:            p.ID,
			Name:          p.Name,
			Email:         p.Email,
			AccountBalance: p.AccountBalance,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}