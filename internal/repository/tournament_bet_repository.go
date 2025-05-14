package repository

import (
	"context"
	"database/sql"
	"igaming/internal/models"
)

type TournamentBetRepository struct {
	db *sql.DB
}

func NewTournamentBetRepository(db *sql.DB) *TournamentBetRepository {
	return &TournamentBetRepository{db: db}
}

func (r *TournamentBetRepository) Create(ctx context.Context, tournamentBet *models.TournamentBet) error {
	query := `INSERT INTO Tournament_bets 
	(player_id, tournamet_id, bet_amount) 
	VALUES (?, ?, ?)`

	result, err := r.db.ExecContext(
		ctx, 
		query, 
		tournamentBet.PlayerID, 
		tournamentBet.TournamentID, 
		tournamentBet.BetAmount)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	tournamentBet.ID = uint(id)

	return nil
}