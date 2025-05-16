package models

import (
	"time"
)

// Tournament represents a competitive gaming event
// swagger:model Tournament
type Tournament struct {
	// The unique identifier for the tournament
	// example: 1
	ID        uint      `json:"id"`
	
	// Name of the tournament
	// required: true
	// example: World Championship
	Name      string    `json:"name"`
	
	// Total prize pool in USD
	// required: true
	// minimum: 0
	// example: 100000.00
	PrizePool float64   `json:"prize_pool"`
	
	// Start date/time of the tournament
    // required: true
    // format: date-time
    // example: 2023-09-01T15:00:00Z
    StartDate time.Time `json:"start_date" swaggertype:"string" format:"date-time"`
    
    // End date/time of the tournament
    // required: true
    // format: date-time
    // example: 2023-09-05T18:00:00Z
    EndDate time.Time `json:"end_date" swaggertype:"string" format:"date-time"`
    
    // Creation timestamp
    // readOnly: true
    // format: date-time
    // example: 2023-08-25T09:30:00Z
    CreatedAt time.Time `json:"created_at" swaggertype:"string" format:"date-time"`
	
	// Timestamp when tournament was last updated
	// readOnly: true
	// example: 2023-08-28T14:45:00Z
	UpdatedAt time.Time `json:"updated_at"`
}