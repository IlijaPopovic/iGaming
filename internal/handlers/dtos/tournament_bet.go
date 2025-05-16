package dtos

import "time"

// CreateTournamentBetRequest represents a bet placement
type CreateTournamentBetRequest struct {
	// ID of the player making the bet
	// required: true
	// example: 123
	PlayerID uint `json:"player_id" validate:"required"`
	
	// ID of the tournament to bet on
	// required: true
	// example: 456
	TournamentID uint `json:"tournament_id" validate:"required"`
	
	// Amount to wager in USD
	// required: true
	// minimum: 0.01
	// example: 50.00
	BetAmount float64 `json:"bet_amount" validate:"required,gt=0"`
}

// TournamentBetResponse represents a placed bet
type TournamentBetResponse struct {
	// Bet ID
	// example: 1
	ID uint `json:"id"`
	
	// Player ID
	// example: 123
	PlayerID uint `json:"player_id"`
	
	// Tournament ID
	// example: 456
	TournamentID uint `json:"tournament_id"`
	
	// Wagered amount
	// example: 50.00
	BetAmount float64 `json:"bet_amount"`
	
	// Bet placement timestamp
	// example: 2023-09-01T10:15:00Z
	CreatedAt time.Time `json:"created_at"`
}