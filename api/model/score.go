package model

type Score struct {
	UserID int     `json:"-" db:"user_id"`
	User   *User   `json:"user" db:"-"`
	Score  float64 `json:"score"`
}
