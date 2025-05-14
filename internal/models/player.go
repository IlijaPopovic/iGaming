package models

import "time"

type Player struct {
	ID uint `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	AccountBalance float64 `json:"account_balance"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}