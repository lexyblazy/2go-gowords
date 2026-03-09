package store

type UserEntity struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Moniker      string `json:"moniker"`
	Password     string `json:"password"`
	CreatedAt    string `json:"created_at"`
	RecoveryHash string `json:"recovery_hash"`
}

type UserStatsEntity struct {
	UserId      string `json:"user_id"`
	GamesPlayed int    `json:"games_played"`
	Wins        int    `json:"wins"`
	BestScore   int    `json:"best_score"`
	TotalScore  int    `json:"total_score"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserStatsUpdate struct {
	UserID   string
	IsWinner bool
	Score    int
}
