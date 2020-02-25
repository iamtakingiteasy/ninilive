package model

// Channel in chat
type Channel struct {
	ID    uint64 `db:"channel_id"`
	Name  string `db:"channel_name"`
	Order int    `db:"channel_order"`
}
