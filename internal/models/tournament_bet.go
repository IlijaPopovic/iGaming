package models

import (
	"time"
)

type BetStatus string

type TournamentBet struct {
	ID           uint       `json:"id"`
	PlayerID     uint       `json:"player_id"`
	TournamentID uint       `json:"tournament_id"`
	BetAmount    float64    `json:"bet_amount"`
	CreatedAt    time.Time  `json:"created_at"`
	
	Player     *Player     `json:"player,omitempty"`
	Tournament *Tournament `json:"tournament"`
}