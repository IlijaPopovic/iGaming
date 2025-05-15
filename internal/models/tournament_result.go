package models

import "time"

type TournamentResult struct {
	ID           uint      `json:"id"`
	TournamentID uint      `json:"tournament_id"`
	PlayerID     uint      `json:"player_id"`
	Placement    int       `json:"placement"`
	PrizeAmount  float64   `json:"prize_amount"`
	CreatedAt    time.Time `json:"created_at"`
	
	Player     *Player     `json:"player"`
	Tournament *Tournament `json:"tournament"`
}