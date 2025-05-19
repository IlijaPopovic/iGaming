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
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    var currentBalance float64
    err = tx.QueryRowContext(ctx,
        "SELECT account_balance FROM players WHERE id = ? FOR UPDATE",
        bet.PlayerID,
    ).Scan(&currentBalance)
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return fmt.Errorf("player with ID %d does not exist", bet.PlayerID)
        }
        return fmt.Errorf("failed to get player balance: %w", err)
    }

    var tournamentExists bool
    err = tx.QueryRowContext(ctx,
        "SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = ?)",
        bet.TournamentID,
    ).Scan(&tournamentExists)
    
    if err != nil || !tournamentExists {
        return fmt.Errorf("tournament with ID %d does not exist", bet.TournamentID)
    }

    if currentBalance < bet.BetAmount {
        return fmt.Errorf("insufficient funds: player has %.2f, needs %.2f", 
            currentBalance, bet.BetAmount)
    }

    _, err = tx.ExecContext(ctx,
        "UPDATE players SET account_balance = account_balance - ? WHERE id = ?",
        bet.BetAmount, 
        bet.PlayerID,
    )
    if err != nil {
        return fmt.Errorf("failed to deduct funds: %w", err)
    }

    result, err := tx.ExecContext(ctx,
        `INSERT INTO tournament_bets (player_id, tournament_id, bet_amount) 
         VALUES (?, ?, ?)`,
        bet.PlayerID, 
        bet.TournamentID, 
        bet.BetAmount,
    )
    if err != nil {
        return fmt.Errorf("failed to create bet: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("transaction commit failed: %w", err)
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


// THIS FUNCTION IS USING GetPlayerByID and GetTorunametByID, but not a good implementation for transaction

// func (r *TournamentBetRepository) Create(ctx context.Context, bet *models.TournamentBet) error {
// 	if _, err := r.playerRepo.GetPlayerByID(ctx, bet.PlayerID); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return fmt.Errorf("player with ID %d does not exist", bet.PlayerID)
// 		}
// 		return fmt.Errorf("failed to validate player: %w", err)
// 	}

// 	if _, err := r.tournamentRepo.GetTournamentByID(ctx, bet.TournamentID); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return fmt.Errorf("tournament with ID %d does not exist", bet.TournamentID)
// 		}
// 		return fmt.Errorf("failed to validate tournament: %w", err)
// 	}

// 	query := `INSERT INTO tournament_bets 
// 	(player_id, tournament_id, bet_amount) 
// 	VALUES (?, ?, ?)`

// 	result, err := r.db.ExecContext(
// 		ctx,
// 		query,
// 		bet.PlayerID,
// 		bet.TournamentID,
// 		bet.BetAmount,
// 	)

// 	if err != nil {
// 		return fmt.Errorf("database error: %w", err)
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return fmt.Errorf("failed to get last insert ID: %w", err)
// 	}

// 	bet.ID = uint(id)
// 	return nil
// }