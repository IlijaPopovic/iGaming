package repository

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *TournamentRepository) GetAll(ctx context.Context) ([]models.Tournament, error) {
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