package postgres

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/iamtakingiteasy/ninilive/internal/db/postgres/migration"
)

func migrateDB(db *sql.DB) error {
	source, err := bindata.WithInstance(bindata.Resource(migration.AssetNames(), migration.Asset))
	if err != nil {
		return fmt.Errorf("failed to prepare bindata source: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to prepare driver connection: %v", err)
	}

	m, err := migrate.NewWithInstance("go-bindata", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to prepare migration: %v", err)
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		err = nil
	}

	return err
}
