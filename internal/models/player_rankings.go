package models

type PlayerRanking struct {
    PlayerID       uint    `json:"player_id" db:"player_id"`
    PlayerName     string  `json:"player_name" db:"player_name"`
    AccountBalance float64 `json:"account_balance" db:"account_balance"`
    Rank           int     `json:"rank" db:"player_rank"`
}