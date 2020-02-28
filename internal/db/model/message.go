package model

import (
	"time"
)

// Message in chat
type Message struct {
	ID        uint64    `db:"message_id"`
	ChannelID uint64    `db:"message_channel_id"`
	Body      string    `db:"message_body"`
	Time      time.Time `db:"message_time"`
	Edit      time.Time `db:"message_edit"`
	Trip      string    `db:"message_trip"`
	Name      string    `db:"message_name"`
	Origin    string    `db:"message_origin"`
	Remote    string    `db:"message_remote"`
	FileName  string    `db:"message_file_name"`
	FilePath  string    `db:"message_file_path"`
	User
}
