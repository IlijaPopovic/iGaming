package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"igaming/internal/models"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Create(ctx context.Context, player *models.Player) error {
	query := `INSERT INTO players 
	(name, email, password_hash, account_balance) 
	VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(
		ctx,
		query,
		player.Name,
		player.Email,
		player.PasswordHash,
		player.AccountBalance,
	)

	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	player.ID = uint(id)
	return nil
}

func (r *PlayerRepository) GetAllPlayers(ctx context.Context) ([]models.Player, error) {
	query := `SELECT 
		id, name, email, account_balance, created_at, updated_at, deleted_at 
		FROM players`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query players: %w", err)
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var p models.Player
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Email,
			&p.AccountBalance,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player row: %w", err)
		}
		players = append(players, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return players, nil
}

func (r *PlayerRepository) GetPlayerByID(ctx context.Context, id uint) (*models.Player, error) {
    query := `SELECT 
        id, name, email, password_hash, account_balance, created_at, updated_at, deleted_at 
        FROM players 
        WHERE id = ?`

    row := r.db.QueryRowContext(ctx, query, id)
    
    var player models.Player
    err := row.Scan(
        &player.ID,
        &player.Name,
        &player.Email,
        &player.PasswordHash,
        &player.AccountBalance,
        &player.CreatedAt,
        &player.UpdatedAt,
        &player.DeletedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("player with ID %d not found", id)
        }
        return nil, fmt.Errorf("failed to get player: %w", err)
    }
    
    return &player, nil
}