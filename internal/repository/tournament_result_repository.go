package repository

import (
	"context"
	"database/sql"
	"igaming/internal/models"
)

type TournamentResultRepository struct {
	db *sql.DB
}

func NewTournamentResultRepository(db *sql.DB) *TournamentResultRepository {
	return &TournamentResultRepository{db: db}
}

func (r *TournamentResultRepository) Create(ctx context.Context, tournamentResult *models.TournamentResult) error {
	query := `INSERT INTO tournament_results
	(tournament_id, player_id, placement, prize_amount) 
	VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(
		ctx, 
		query, 
		tournamentResult.TournamentID, 
		tournamentResult.PlayerID, 
		tournamentResult.Placement, 
		tournamentResult.PrizeAmount)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	tournamentResult.ID = uint(id)

	return nil
}