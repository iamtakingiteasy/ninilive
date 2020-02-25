package model

// User in chat
type User struct {
	ID       uint64 `db:"user_id"`
	Name     string `db:"user_name"`
	Login    string `db:"user_login"`
	Password string `db:"user_password"`
	Mod      bool   `db:"user_mod"`
}
