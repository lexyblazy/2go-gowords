package store

type UserEntity struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Moniker   string `json:"moniker"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}
