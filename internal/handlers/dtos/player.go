package dtos

import "time"

// CreatePlayerRequest represents the payload for creating a player
type CreatePlayerRequest struct {
	// Player's display name
	// required: true
	// example: JohnDoe123
	Name string `json:"name" validate:"required,min=2,max=100"`
	
	// Player's email address
	// required: true
	// format: email
	// example: john.doe@example.com
	Email string `json:"email" validate:"required,email"`
	
	// Player's password
	// required: true
	// minLength: 8
	// example: securePassword123!
	Password string `json:"password" validate:"required,min=8"`
	
	// Initial account balance
	// minimum: 0
	// example: 100.00
	AccountBalance float64 `json:"account_balance" validate:"gte=0"`
}

// PlayerResponse represents a player API response
type PlayerResponse struct {
	// The player ID
	// example: 1
	ID uint `json:"id"`
	
	// The player's display name
	// example: JohnDoe123
	Name string `json:"name"`
	
	// The player's email
	// example: john.doe@example.com
	Email string `json:"email"`
	
	// Current account balance
	// example: 150.50
	AccountBalance float64 `json:"account_balance"`
	
	// Account creation timestamp
	// example: 2023-08-15T14:30:45Z
	CreatedAt time.Time `json:"created_at"`
	
	// Last update timestamp
	// example: 2023-08-16T09:15:22Z
	UpdatedAt time.Time `json:"updated_at"`
}