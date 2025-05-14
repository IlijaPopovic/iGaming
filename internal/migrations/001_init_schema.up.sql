-- ======================
-- CORE TABLES
-- ======================

-- ======================
-- PLAYERS TABLE
-- ======================
CREATE TABLE players (
    -- Primary Key
    id INT AUTO_INCREMENT PRIMARY KEY,
    
    -- Player Information
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    
    -- Security & Authentication
    password_hash VARCHAR(255) NOT NULL,
    
    -- Financial Data
    account_balance DECIMAL(15, 2) DEFAULT 0.00,
    
    -- -- Profile Information
    -- phone VARCHAR(20),
    -- address TEXT,
    -- timezone VARCHAR(50) DEFAULT 'UTC',
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    -- Constraints
    UNIQUE KEY (email),
    CONSTRAINT chk_player_email_format CHECK (email REGEXP '^[^@]+@[^@]+\.[^@]+$')
);

-- ======================
-- TOURNAMENTS TABLE
-- ======================
CREATE TABLE tournaments (
    -- Primary Key
    id INT AUTO_INCREMENT PRIMARY KEY,
    
    -- Tournament Information
    name VARCHAR(100) NOT NULL,
    prize_pool DECIMAL(15, 2) NOT NULL,
    
    -- Event Timing
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_valid_dates CHECK (end_date > start_date)
);

-- ======================
-- TOURNAMENT BETS TABLE
-- ======================
CREATE TABLE tournament_bets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    player_id INT NOT NULL,
    tournament_id INT NOT NULL,
    
    -- Bet Information
    bet_amount DECIMAL(15, 2) NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    FOREIGN KEY (player_id) REFERENCES players(id),
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id)
);

-- ======================
-- TOURNAMENT RESULTS TABLE
-- ======================
CREATE TABLE tournament_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    -- Foreign Keys
    tournament_id INT NOT NULL,
    player_id INT NOT NULL,
    
    -- Competition Data
    placement INT NOT NULL,
    prize_amount DECIMAL(15, 2) NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    UNIQUE KEY unique_tournament_player (tournament_id, player_id)
    CONSTRAINT chk_valid_placement CHECK (placement BETWEEN 1 AND 3)
);

-- ======================
-- PROCEDURES
-- ======================

-- ======================
-- PRIZE DISTRIBUTION PROCEDURE
-- ======================
DELIMITER //
CREATE PROCEDURE DistributePrizes(IN target_tournament_id INT)
BEGIN
    DECLARE total_prize_pool DECIMAL(15, 2);
    
    -- Get tournament prize pool
    SELECT prize_pool INTO total_prize_pool 
    FROM tournaments 
    WHERE id = target_tournament_id;
    
     -- Calculate rankings and prizes using CTEs
    WITH BetCounts AS (
        SELECT 
            player_id,
            COUNT(*) AS total_bets
        FROM tournament_bets
        WHERE tournament_id = target_tournament_id
        GROUP BY player_id
    ),
    RankedPlayers AS (
        SELECT
            player_id,
            total_bets,
            DENSE_RANK() OVER (ORDER BY total_bets DESC) AS placement
        FROM BetCounts
    ),
    PrizeDistribution AS (
        SELECT
            RankedPlayers.player_id,
            RankedPlayers.placement,
            CASE RankedPlayers.placement
                WHEN 1 THEN total_prize_pool * 0.5 / COUNT(*) OVER (PARTITION BY r.placement)
                WHEN 2 THEN total_prize_pool * 0.3 / COUNT(*) OVER (PARTITION BY r.placement)
                WHEN 3 THEN total_prize_pool * 0.2 / COUNT(*) OVER (PARTITION BY r.placement)
                ELSE 0
            END AS prize
        FROM RankedPlayers
        WHERE RankedPlayers.placement <= 3
    )

    -- Update tournament results
    INSERT INTO tournament_results (tournament_id, player_id, placement, prize_amount)
    SELECT
        target_tournament_id,
        pd.player_id,
        pd.placement,
        pd.prize
    FROM PrizeDistribution pd
    ON DUPLICATE KEY UPDATE
        placement = pd.placement,
        prize_amount = pd.prize;

    -- Update player balances
    UPDATE players p
    JOIN PrizeDistribution pd ON p.id = pd.player_id
    SET p.account_balance = p.account_balance + pd.prize;
END //
DELIMITER ;

-- ======================
-- VIEWS
-- ======================

-- ======================
-- PLAYER RANKINGS VIEW
-- ======================
CREATE VIEW player_rankings AS
WITH ranked_players AS (
    SELECT 
        -- Player Identification
        player_id,
        player_name,
        
        -- Financial Information
        account_balance,
        
        -- Ranking Calculations
        RANK() OVER (ORDER BY account_balance DESC) as player_rank,
        DENSE_RANK() OVER (ORDER BY account_balance DESC) as player_dense_rank,
        ROW_NUMBER() OVER (ORDER BY account_balance DESC) as player_row_num
    FROM players
    WHERE deleted_at IS NULL  -- Exclude soft-deleted players
)
SELECT 
    player_id,
    player_name,
    account_balance,
    player_rank,
    player_dense_rank,
    player_row_num
FROM ranked_players;

-- ======================
-- INDEXES
-- ======================
CREATE INDEX idx_tournaments_dates ON tournaments(start_date, end_date);
CREATE INDEX idx_players_balance ON players(account_balance DESC);
CREATE INDEX idx_results_placement ON tournament_results(placement);

CREATE INDEX idx_bets_tournament_player ON tournament_bets(tournament_id, player_id);
CREATE INDEX idx_players_email ON players(email);
CREATE INDEX idx_tournaments_name ON tournaments(name);
CREATE INDEX idx_bets_created ON tournament_bets(created_at);
CREATE INDEX idx_results_created ON tournament_results(created_at);