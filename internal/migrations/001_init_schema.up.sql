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

    -- Temp table: Placement counts
    CREATE TEMPORARY TABLE tmp_placement_counts (
        placement INT,
        player_count INT
    );

    INSERT INTO tmp_placement_counts
    SELECT placement, COUNT(*) AS player_count
    FROM tmp_ranked_players
    WHERE placement <= 3
    GROUP BY placement;

    -- Temp table: Prize distribution
    CREATE TEMPORARY TABLE tmp_prize_distribution (
        player_id INT,
        placement INT,
        prize DECIMAL(15, 2)
    );
    -- >>>>>>>>>>>CHAGE THE SPLIT LATTER IT DOES NOT MAKE SENSE!
    INSERT INTO tmp_prize_distribution (player_id, placement, prize)
    SELECT
        rp.player_id,
        rp.placement,
        CASE rp.placement
            WHEN 1 THEN total_prize_pool * 0.5 / COALESCE(pc.player_count, 1)
            WHEN 2 THEN total_prize_pool * 0.3 / COALESCE(pc.player_count, 1)
            WHEN 3 THEN total_prize_pool * 0.2 / COALESCE(pc.player_count, 1)
            ELSE 0
        END AS prize
    FROM tmp_ranked_players rp
    LEFT JOIN tmp_placement_counts pc ON rp.placement = pc.placement
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
    DROP TEMPORARY TABLE IF EXISTS tmp_placement_counts;
    DROP TEMPORARY TABLE IF EXISTS tmp_prize_distribution;
END;
-- +goose StatementEnd

CREATE VIEW player_rankings AS
SELECT 
    id AS player_id,
    name AS player_name,
    account_balance,
    DENSE_RANK() OVER (ORDER BY account_balance DESC) AS player_rank
FROM players
WHERE deleted_at IS NULL
ORDER BY player_rank;

CREATE INDEX idx_tournaments_dates ON tournaments(start_date, end_date);
CREATE INDEX idx_players_balance ON players(account_balance DESC);
CREATE INDEX idx_results_placement ON tournament_results(placement);
CREATE INDEX idx_bets_tournament_player ON tournament_bets(tournament_id, player_id);
CREATE INDEX idx_players_email ON players(email);
CREATE INDEX idx_tournaments_name ON tournaments(name);
CREATE INDEX idx_bets_created ON tournament_bets(created_at);
CREATE INDEX idx_results_created ON tournament_results(created_at);

INSERT INTO players (name, email, password_hash, account_balance) VALUES
('Alice Smith', 'alice.smith@pokermail.com', '$2a$10$W6c8Ua5uO7yj5J2', 1500.00),
('Bob Johnson', 'bob.johnson@pokermail.com', '$2a$10$ZR9tG4bM2wD1vE3', 8750.00),
('Charlie Brown', 'charlie.brown@pokermail.com', '$2a$10$XKp7Q2rN4sH6fT8', 4200.00),
('Diana Miller', 'diana.miller@pokermail.com', '$2a$10$YL3vM9wP6tR7sS2', 15600.00),
('Evan Davis', 'evan.davis@pokermail.com', '$2a$10$BP4nV8cJ3hG5dF1', 9500.00),
('Fiona Clark', 'fiona.clark@pokermail.com', '$2a$10$QW2e5rT9yH4jK7L', 3200.00),
('George Wilson', 'george.wilson@pokermail.com', '$2a$10$AS1dF3gH6jK8L9P', 12800.00),
('Hannah White', 'hannah.white@pokermail.com', '$2a$10$ZX3cV4bN5m6M7Q8', 6400.00),
('Ian Moore', 'ian.moore@pokermail.com', '$2a$10$RT6yT7uI8o9P0Q1', 2300.00),
('Jenny Taylor', 'jenny.taylor@pokermail.com', '$2a$10$EK4jL5mN6bV3C2X', 5100.00);

INSERT INTO tournaments (name, prize_pool, start_date, end_date) VALUES
('Winter Classic', 25000.00, '2023-01-10 14:00:00', '2023-01-12 22:00:00'),
('Spring Championship', 50000.00, '2023-03-15 12:00:00', '2023-03-18 20:00:00'),
('Summer Showdown', 75000.00, '2023-06-01 10:00:00', '2023-06-05 18:00:00'),
('Autumn Royale', 100000.00, '2023-09-10 16:00:00', '2023-09-15 23:59:59'),
('Masters Invitational', 150000.00, '2023-11-01 09:00:00', '2023-11-05 21:00:00'),
('Weekend Warmup', 10000.00, '2023-02-05 08:00:00', '2023-02-05 20:00:00'),
('High Roller Event', 200000.00, '2023-07-20 12:00:00', '2023-07-25 12:00:00'),
('Fast Fold Frenzy', 30000.00, '2023-04-10 18:00:00', '2023-04-12 18:00:00'),
('New Year Knockout', 5000.00, '2023-12-31 23:00:00', '2024-01-01 06:00:00'),
('Satellite Special', 15000.00, '2023-08-15 10:00:00', '2023-08-16 22:00:00');

INSERT INTO tournament_bets (player_id, tournament_id, bet_amount) VALUES
(1, 1, 500.00), (1, 1, 300.00),
(2, 2, 1000.00), (2, 2, 500.00),
(3, 3, 750.00),
(4, 4, 1500.00), (4, 4, 1000.00),
(5, 5, 2000.00),
(6, 6, 250.00), (6, 6, 150.00),
(7, 7, 3000.00),
(8, 8, 600.00),
(9, 9, 100.00), (9, 9, 50.00),
(10, 10, 450.00);

-- +goose Down

DROP VIEW IF EXISTS player_rankings;
DROP PROCEDURE IF EXISTS DistributePrizes;
DROP TABLE IF EXISTS tournament_results;
DROP TABLE IF EXISTS tournament_bets;
DROP TABLE IF EXISTS tournaments;
DROP TABLE IF EXISTS players;
