package models

import (
	"time"
)

type Tournament struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	PrizePool float64   `json:"prize_pool"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


