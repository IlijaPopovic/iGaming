package repository

import (
	"context"
	"database/sql"
	"igaming/internal/models"
)

type TournamentRepository struct {
	db *sql.DB
}

func NewTournamentRepository(db *sql.DB) *TournamentRepository {
	return &TournamentRepository{db: db}
}

func (r *TournamentRepository) Create(ctx context.Context, tournament *models.Tournament) error {
	query := `INSERT INTO Tournament 
	(name, prize_pool, start_date, end_date) 
	VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(
		ctx, 
		query, 
		tournament.Name, 
		tournament.PrizePool, 
		tournament.StartDate, 
		tournament.EndDate)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	tournament.ID = uint(id)

	return nil
}