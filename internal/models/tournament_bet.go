package models

import (
	"time"
)

// TournamentBet represents a player's wager on a tournament
// swagger:model TournamentBet
type TournamentBet struct {
	// The unique identifier for the bet
	// example: 1
	ID           uint       `json:"id"`
	
	// ID of the player making the bet
	// required: true
	// example: 123
	PlayerID     uint       `json:"player_id"`
	
	// ID of the tournament being bet on
	// required: true
	// example: 456
	TournamentID uint       `json:"tournament_id"`
	
	// Amount wagered in USD
	// required: true
	// minimum: 0.01
	// example: 50.00
	BetAmount    float64    `json:"bet_amount"`
	
	// Timestamp when the bet was placed
	// readOnly: true
	// example: 2023-09-01T10:15:00Z
	CreatedAt    time.Time  `json:"created_at"`
	
	// Player details (optional in responses)
	// swagger:ignore
	Player     *Player     `json:"player,omitempty"`
	
	// Tournament details
	// required: true
	Tournament *Tournament `json:"tournament"`
}