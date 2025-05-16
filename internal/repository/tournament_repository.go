package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"igaming/internal/models"
	"log"
)

type TournamentRepository struct {
	db *sql.DB
}

func NewTournamentRepository(db *sql.DB) *TournamentRepository {
	return &TournamentRepository{db: db}
}

func (r *TournamentRepository) Create(ctx context.Context, tournament *models.Tournament) error {
    query := `INSERT INTO tournaments 
    (name, prize_pool, start_date, end_date) 
    VALUES (?, ?, ?, ?)`

    result, err := r.db.ExecContext(
        ctx, 
        query, 
        tournament.Name, 
        tournament.PrizePool, 
        tournament.StartDate, 
        tournament.EndDate,
    )

    if err != nil {
        log.Printf("Database error: %v", err)
        return fmt.Errorf("database operation failed: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        log.Printf("LastInsertId error: %v", err)
        return fmt.Errorf("failed to get last insert ID: %w", err)
    }

    tournament.ID = uint(id)
    return nil
}

func (r *TournamentRepository) GetAllTournaments(ctx context.Context) ([]models.Tournament, error) {
    query := `SELECT id, name, prize_pool, start_date, end_date, created_at, updated_at FROM tournaments`
    
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query tournaments: %w", err)
    }
    defer rows.Close()

    var tournaments []models.Tournament
    for rows.Next() {
        var t models.Tournament
        err := rows.Scan(
            &t.ID,
            &t.Name,
            &t.PrizePool,
            &t.StartDate,
            &t.EndDate,
            &t.CreatedAt,
            &t.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan tournament row: %w", err)
        }
        tournaments = append(tournaments, t)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return tournaments, nil
}

func (r *TournamentRepository) GetTournamentByID(ctx context.Context, id uint) (*models.Tournament, error) {
    query := `SELECT 
        id, name, prize_pool, start_date, end_date, created_at, updated_at 
        FROM tournaments 
        WHERE id = ?`

    row := r.db.QueryRowContext(ctx, query, id)
    
    var tournament models.Tournament
    err := row.Scan(
        &tournament.ID,
        &tournament.Name,
        &tournament.PrizePool,
        &tournament.StartDate,
        &tournament.EndDate,
        &tournament.CreatedAt,
        &tournament.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("tournament with ID %d not found", id)
        }
        return nil, fmt.Errorf("failed to get tournament: %w", err)
    }
    
    return &tournament, nil
}

func (r *TournamentRepository) DistributePrizes(ctx context.Context, tournamentID uint) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    _, err = tx.ExecContext(ctx, "CALL DistributePrizes(?)", tournamentID)
    if err != nil {
        return fmt.Errorf("prize distribution failed: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

func (r *TournamentRepository) Exists(ctx context.Context, id uint) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM tournaments WHERE id = ?)"
    err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
    return exists, err
}