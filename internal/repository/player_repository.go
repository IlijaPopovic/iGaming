package repository

import (
	"context"
	"database/sql"
	"igaming/internal/models"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Create(ctx context.Context, player *models.Player) error {
	query := `INSERT INTO players (name, email,password_hash,account_balance) VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(
		ctx, 
		query, 
		player.Name, 
		player.Email, 
		player.PasswordHash, 
		player.AccountBalance)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	player.ID = uint(id)

	return nil
}