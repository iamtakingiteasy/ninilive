// Package postgres for postgresql persistency
package postgres

import (
	"fmt"

	"github.com/iamtakingiteasy/ninilive/internal/config"
	"github.com/iamtakingiteasy/ninilive/internal/db"

	"github.com/jmoiron/sqlx"
	// postgres driver
	_ "github.com/lib/pq"
)

// NewPersister returns new postgresql persister for given connection url
func NewPersister(values *config.ValuesDB) (db.Persister, error) {
	conn, err := sqlx.Connect("postgres", values.URL)
	if err != nil {
		return nil, err
	}

	err = migrateDB(conn.DB)
	if err != nil {
		return nil, err
	}

	deleteMessage, err := conn.PrepareNamed(deleteMessageQuery)
	if err != nil {
		return nil, fmt.Errorf("deleteMessageQuery: %w", err)
	}

	saveMessage, err := conn.PrepareNamed(saveMessageQuery)
	if err != nil {
		return nil, fmt.Errorf("saveMessageQuery: %w", err)
	}

	editMessage, err := conn.PrepareNamed(editMessageQuery)
	if err != nil {
		return nil, fmt.Errorf("editMessageQuery: %w", err)
	}

	loadMessagesLast, err := conn.Preparex(loadMessagesLastQuery)
	if err != nil {
		return nil, fmt.Errorf("loadMessagesLastQuery: %w", err)
	}

	loadMessagesBefore, err := conn.Preparex(loadMessagesBeforeQuery)
	if err != nil {
		return nil, fmt.Errorf("loadMessagesBeforeQuery: %w", err)
	}

	loadMessagesPage, err := conn.Preparex(loadMessagesPageQuery)
	if err != nil {
		return nil, fmt.Errorf("loadMessagesPageQuery: %w", err)
	}

	deleteUser, err := conn.Preparex(deleteUserQuery)
	if err != nil {
		return nil, fmt.Errorf("deleteUserQuery: %w", err)
	}

	loadUser, err := conn.Preparex(loadUserQuery)
	if err != nil {
		return nil, fmt.Errorf("loadUserQuery: %w", err)
	}

	loadUserByID, err := conn.Preparex(loadUserByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("loadUserQuery: %w", err)
	}

	saveUser, err := conn.PrepareNamed(saveUserQuery)
	if err != nil {
		return nil, fmt.Errorf("saveUserQuery: %w", err)
	}

	ensureUser, err := conn.PrepareNamed(ensureUserQuery)
	if err != nil {
		return nil, fmt.Errorf("ensureUserQuery: %w", err)
	}

	loadUsersPage, err := conn.Preparex(loadUsersPageQuery)
	if err != nil {
		return nil, fmt.Errorf("loadUsersPageQuery: %w", err)
	}

	deleteChannel, err := conn.Preparex(deleteChannelQuery)
	if err != nil {
		return nil, fmt.Errorf("deleteChannelQuery: %w", err)
	}

	saveChannel, err := conn.PrepareNamed(saveChannelQuery)
	if err != nil {
		return nil, fmt.Errorf("saveChannelQuery: %w", err)
	}

	loadChannels, err := conn.Preparex(loadChannelsQuery)
	if err != nil {
		return nil, fmt.Errorf("loadChannelsQuery: %w", err)
	}

	return &persister{
		values:             values,
		db:                 conn,
		deleteMessage:      deleteMessage,
		saveMessage:        saveMessage,
		editMessage:        editMessage,
		loadMessagesLast:   loadMessagesLast,
		loadMessagesBefore: loadMessagesBefore,
		loadMessagesPage:   loadMessagesPage,
		deleteUser:         deleteUser,
		loadUser:           loadUser,
		loadUserByID:       loadUserByID,
		saveUser:           saveUser,
		ensureUser:         ensureUser,
		loadUsersPage:      loadUsersPage,
		deleteChannel:      deleteChannel,
		saveChannel:        saveChannel,
		loadChannels:       loadChannels,
	}, nil
}
