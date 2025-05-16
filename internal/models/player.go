package models

import "time"

// Player represents a user in the gaming system
// swagger:model Player
type Player struct {
	// The unique identifier for the player
	// example: 1
	ID uint `json:"id"`
	
	// The player's display name
	// required: true
	// example: JohnDoe123
	Name  string `json:"name"`
	
	// The player's email address
	// required: true
	// format: email
	// example: john.doe@example.com
	Email string `json:"email"`
	
	// Password hash
	// swagger:ignore
	PasswordHash string `json:"-"`
	
	// Current account balance
	// minimum: 0
	// example: 100.50
	AccountBalance float64 `json:"account_balance"`
	
	// Timestamp when the player was created
	// readOnly: true
	// example: 2023-08-15T14:30:45Z
	CreatedAt time.Time  `json:"created_at"`
	
	// Timestamp when the player was last updated
	// readOnly: true
	// example: 2023-08-16T09:15:22Z
	UpdatedAt time.Time  `json:"updated_at"`
	
	// Timestamp when the player was deleted (if applicable)
	// nullable: true
	// example: 2023-08-17T16:45:00Z
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}