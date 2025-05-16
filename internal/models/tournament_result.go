package models

import "time"

// TournamentResult represents a player's final standing in a tournament
// swagger:model TournamentResult
type TournamentResult struct {
	// The unique identifier for the result entry
	// example: 1
	ID           uint      `json:"id"`
	
	// ID of the tournament
	// required: true
	// example: 456
	TournamentID uint      `json:"tournament_id"`
	
	// ID of the player
	// required: true
	// example: 123
	PlayerID     uint      `json:"player_id"`
	
	// Final placement position (1-based)
	// required: true
	// minimum: 1
	// maximum: 3
	// example: 1
	Placement    int       `json:"placement"`
	
	// Prize money awarded in USD
	// required: true
	// minimum: 0
	// example: 5000.00
	PrizeAmount  float64   `json:"prize_amount"`
	
	// Timestamp when the result was recorded
	// readOnly: true
	// example: 2023-09-05T18:30:00Z
	CreatedAt    time.Time `json:"created_at"`
	
	// Player details
	// swagger:ignore  // Prevents recursion in documentation
	Player     *Player     `json:"player"`
	
	// Tournament details
	// required: true
	Tournament *Tournament `json:"tournament"`
}