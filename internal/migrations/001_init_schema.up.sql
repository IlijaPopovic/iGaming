-- +goose Up

CREATE TABLE players (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    account_balance DECIMAL(15, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    UNIQUE KEY (email),
    CONSTRAINT chk_player_email_format CHECK (email REGEXP '^[^@]+@[^@]+\.[^@]+$')
) ENGINE=InnoDB;

CREATE TABLE tournaments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    prize_pool DECIMAL(15, 2) NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT chk_valid_dates CHECK (end_date > start_date)
) ENGINE=InnoDB;

CREATE TABLE tournament_bets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    player_id INT NOT NULL,
    tournament_id INT NOT NULL,
    bet_amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (player_id) REFERENCES players(id),
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id)
) ENGINE=InnoDB;

CREATE TABLE tournament_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tournament_id INT NOT NULL,
    player_id INT NOT NULL,
    placement INT NOT NULL,
    prize_amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_tournament_player (tournament_id, player_id),
    CONSTRAINT chk_valid_placement CHECK (placement BETWEEN 1 AND 3),
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id),
    FOREIGN KEY (player_id) REFERENCES players(id)
) ENGINE=InnoDB;

-- +goose StatementBegin
CREATE PROCEDURE DistributePrizes(IN target_tournament_id INT)
BEGIN
    DECLARE total_prize_pool DECIMAL(15, 2);

    -- Get the prize pool
    SELECT prize_pool INTO total_prize_pool 
    FROM tournaments 
    WHERE id = target_tournament_id;

    -- Temp table: Bet counts
    CREATE TEMPORARY TABLE tmp_bet_counts (
        player_id INT,
        total_bets INT
    );

    INSERT INTO tmp_bet_counts (player_id, total_bets)
    SELECT player_id, COUNT(*) AS total_bets
    FROM tournament_bets
    WHERE tournament_id = target_tournament_id
    GROUP BY player_id;

    -- Temp table: Ranked players
    CREATE TEMPORARY TABLE tmp_ranked_players (
        player_id INT,
        placement INT
    );

    INSERT INTO tmp_ranked_players (player_id, placement)
    SELECT 
        player_id,
        DENSE_RANK() OVER (ORDER BY total_bets DESC) AS placement
    FROM tmp_bet_counts;

    -- Temp table: Prize distribution
    CREATE TEMPORARY TABLE tmp_prize_distribution (
        player_id INT,
        placement INT,
        prize DECIMAL(15, 2)
    );

    INSERT INTO tmp_prize_distribution (player_id, placement, prize)
    SELECT
        rp.player_id,
        rp.placement,
        CASE rp.placement
            WHEN 1 THEN total_prize_pool * 0.5 / (SELECT COUNT(*) FROM tmp_ranked_players WHERE placement = 1)
            WHEN 2 THEN total_prize_pool * 0.3 / (SELECT COUNT(*) FROM tmp_ranked_players WHERE placement = 2)
            WHEN 3 THEN total_prize_pool * 0.2 / (SELECT COUNT(*) FROM tmp_ranked_players WHERE placement = 3)
            ELSE 0
        END AS prize
    FROM tmp_ranked_players rp
    WHERE rp.placement <= 3;

    -- Insert into tournament_results
    INSERT INTO tournament_results (tournament_id, player_id, placement, prize_amount)
    SELECT
        target_tournament_id,
        player_id,
        placement,
        prize
    FROM tmp_prize_distribution
    ON DUPLICATE KEY UPDATE
        placement = VALUES(placement),
        prize_amount = VALUES(prize_amount);

    -- Update players
    UPDATE players p
    JOIN tmp_prize_distribution pd ON p.id = pd.player_id
    SET p.account_balance = p.account_balance + pd.prize;

    -- Cleanup
    DROP TEMPORARY TABLE IF EXISTS tmp_bet_counts;
    DROP TEMPORARY TABLE IF EXISTS tmp_ranked_players;
    DROP TEMPORARY TABLE IF EXISTS tmp_prize_distribution;
END;
-- +goose StatementEnd

CREATE VIEW player_rankings AS
WITH ranked_players AS (
    SELECT 
        id AS player_id,
        name AS player_name,
        account_balance,
        RANK() OVER (ORDER BY account_balance DESC) as player_rank,
        DENSE_RANK() OVER (ORDER BY account_balance DESC) as player_dense_rank,
        ROW_NUMBER() OVER (ORDER BY account_balance DESC) as player_row_num
    FROM players
    WHERE deleted_at IS NULL
)
SELECT 
    player_id,
    player_name,
    account_balance,
    player_rank,
    player_dense_rank,
    player_row_num
FROM ranked_players;

CREATE INDEX idx_tournaments_dates ON tournaments(start_date, end_date);
CREATE INDEX idx_players_balance ON players(account_balance DESC);
CREATE INDEX idx_results_placement ON tournament_results(placement);
CREATE INDEX idx_bets_tournament_player ON tournament_bets(tournament_id, player_id);
CREATE INDEX idx_players_email ON players(email);
CREATE INDEX idx_tournaments_name ON tournaments(name);
CREATE INDEX idx_bets_created ON tournament_bets(created_at);
CREATE INDEX idx_results_created ON tournament_results(created_at);

-- +goose Down

DROP VIEW IF EXISTS player_rankings;
DROP PROCEDURE IF EXISTS DistributePrizes;
DROP TABLE IF EXISTS tournament_results;
DROP TABLE IF EXISTS tournament_bets;
DROP TABLE IF EXISTS tournaments;
DROP TABLE IF EXISTS players;
