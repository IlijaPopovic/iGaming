package dtos

import "time"

type CreateTournamentRequest struct {
    // Tournament name (3-100 characters)
	// example: tournament123
	// default: tournament123
    Name      string    `json:"name" validate:"required,min=3,max=100"`
    // Prize pool amount (must be positive)
	// example: 3333
	// default: 3333
	PrizePool float64   `json:"prize_pool" validate:"required,gt=0" swaggertype:"number" default:"3333"`
    // format: date-time
    // example: 2023-09-01T15:00:00Z
    StartDate time.Time `json:"start_date" swaggertype:"string" format:"date-time"`
    // format: date-time
    // example: 2023-09-05T18:00:00Z
    EndDate time.Time `json:"end_date" swaggertype:"string" format:"date-time"`
}

type TournamentResponse struct {
    // Tournament ID
    ID        uint      `json:"id"`
    // Tournament name
    Name      string    `json:"name"`
    // Prize pool amount
    PrizePool float64   `json:"prize_pool"`
    // format: date-time
    // example: 2023-09-01T15:00:00Z
    StartDate time.Time `json:"start_date" swaggertype:"string" format:"date-time"`
    // format: date-time
    // example: 2023-09-05T18:00:00Z
    EndDate time.Time `json:"end_date" swaggertype:"string" format:"date-time"`
    // format: date-time
    // example: 2023-08-25T09:30:00Z
    CreatedAt time.Time `json:"created_at" swaggertype:"string" format:"date-time"`
}