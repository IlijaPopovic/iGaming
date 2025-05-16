package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"igaming/internal/models"
)

type TournamentBetRepository struct {
	db             *sql.DB
	playerRepo     *PlayerRepository
	tournamentRepo *TournamentRepository
}

func NewTournamentBetRepository(db *sql.DB, playerRepo *PlayerRepository, tournamentRepo *TournamentRepository) *TournamentBetRepository {
	return &TournamentBetRepository{
		db:             db,
		playerRepo:     playerRepo,
		tournamentRepo: tournamentRepo,
	}
}

func (r *TournamentBetRepository) Create(ctx context.Context, bet *models.TournamentBet) error {
	if _, err := r.playerRepo.GetPlayerByID(ctx, bet.PlayerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("player with ID %d does not exist", bet.PlayerID)
		}
		return fmt.Errorf("failed to validate player: %w", err)
	}

	if _, err := r.tournamentRepo.GetTournamentByID(ctx, bet.TournamentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("tournament with ID %d does not exist", bet.TournamentID)
		}
		return fmt.Errorf("failed to validate tournament: %w", err)
	}

	query := `INSERT INTO tournament_bets 
	(player_id, tournament_id, bet_amount) 
	VALUES (?, ?, ?)`

	result, err := r.db.ExecContext(
		ctx,
		query,
		bet.PlayerID,
		bet.TournamentID,
		bet.BetAmount,
	)

	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	bet.ID = uint(id)
	return nil
}

func (r *TournamentBetRepository) GetAll(ctx context.Context) ([]models.TournamentBet, error) {
	query := `SELECT 
		id, player_id, tournament_id, bet_amount, created_at 
		FROM tournament_bets`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query bets: %w", err)
	}
	defer rows.Close()

	var bets []models.TournamentBet
	for rows.Next() {
		var bet models.TournamentBet
		err := rows.Scan(
			&bet.ID,
			&bet.PlayerID,
			&bet.TournamentID,
			&bet.BetAmount,
			&bet.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bet row: %w", err)
		}
		bets = append(bets, bet)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return bets, nil
}