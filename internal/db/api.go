// Package db for all persistent needs
package db

import (
	"github.com/iamtakingiteasy/ninilive/internal/db/model"
)

// Persister for models to database
type Persister interface {
	// DeleteMessage removes a message
	DeleteMessage(msg *model.Message) (err error)
	// EditMessage updates message from same user
	EditMessage(msg *model.Message) (err error)
	// SaveMessage updates or persists a message
	SaveMessage(msg *model.Message) (err error)
	// LoadMessagesLast loads last reqSize messages, also indicating if there is any more
	LoadMessagesLast(channelID, reqSize uint64) (msgs []*model.Message, more bool, err error)
	// LoadMessagesBefore loads reqSize messages chronologically before given id, also indicating if there is more
	LoadMessagesBefore(channelID, id, reqSize uint64) (msgs []*model.Message, more bool, err error)
	// LoadMessagesPage loads given message page of size reqSize, also indicating actual page size and number of pages
	LoadMessagesPage(channelID, page, reqSize uint64) (msgs []*model.Message, actSize, pages uint64, err error)
	// DeleteUser removes a user
	DeleteUser(id uint64) (err error)
	// LoadUserByID loads user by ID
	LoadUserByID(id uint64) (user *model.User, err error)
	// LoadUser loads user by login and password
	LoadUser(login, password string) (user *model.User, err error)
	// SaveUser updates or persists a user
	SaveUser(user *model.User) (err error)
	// EnsureUser ensures user exists
	EnsureUser(user *model.User) (err error)
	// LoadUsersPage loads given user page of size reqSize, also indicating actual page size and number of pages
	LoadUsersPage(page, reqSize uint64) (users []*model.User, actSize, pages uint64, err error)
	// DeleteChannel removes a channel
	DeleteChannel(id uint64) (err error)
	// SaveChannel updates or persists a channel
	SaveChannel(channel *model.Channel) (err error)
	// LoadChannels loads all channels
	LoadChannels() (channels []*model.Channel, err error)
}
