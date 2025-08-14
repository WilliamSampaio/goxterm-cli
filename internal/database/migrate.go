package database

import (
	"database/sql"
	"fmt"
	"goxterm-cli/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

type Migration struct {
	Id   int
	Name string
	Up   string
	Down string
}

var Migrations = []Migration{
	CreateInfoTable(),
}

func RollbackMigrations() []Migration {
	rollbackMigrations := make([]Migration, len(Migrations))
	for i := range Migrations {
		rollbackMigrations[len(Migrations)-1-i] = Migrations[i]
	}
	return rollbackMigrations
}

func RunMigration(action string, count int) error {

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db := GetDB(cfg.StorePath)

	if err := initializeMigration(db); err != nil {
		return err
	}

	lastMigrationId, err := getLastMigration(db)
	if err != nil {
		return err
	}

	if action == PerformUp && count == 0 {
		for _, migration := range Migrations {
			if migration.Id > lastMigrationId {
				_, err = db.Exec(migration.Up)
				if err != nil {
					return err
				}

				if err := registerMigration(db, migration.Id, migration.Name); err != nil {
					if _, err := db.Exec(migration.Down); err != nil {
						return fmt.Errorf("migration %d failed and rollback also failed: %v", migration.Id, err)
					}

					return fmt.Errorf("migration %d failed: %v", migration.Id, err)
				}
			}
		}
	}

	if action == PerformDown && count == 0 {
		for _, migration := range RollbackMigrations() {
			if migration.Id >= lastMigrationId {
				_, err = db.Exec(migration.Down)
				if err != nil {
					return err
				}

				if err := registerMigration(db, migration.Id, migration.Name); err != nil {
					return fmt.Errorf("rollback %d failed: %v", migration.Id, err)
				}
			}
		}
	}

	return nil
}

func initializeMigration(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS _migrations (
			id INTEGER NOT NULL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);`,
	)

	if err != nil {
		return fmt.Errorf("error creating migrations table: %v", err)
	}

	return nil
}

func getLastMigration(db *sql.DB) (int, error) {
	var id int
	row := db.QueryRow("SELECT id FROM _migrations ORDER BY id DESC LIMIT 1")
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No migrations found
		}
		return 0, fmt.Errorf("error getting last migration: %v", err)
	}
	return id, nil
}

func registerMigration(db *sql.DB, id int, name string) error {
	_, err := db.Exec(`INSERT INTO _migrations (id, name) VALUES (?, ?)`, id, name)
	if err != nil {
		return fmt.Errorf("error registering migration: %v", err)
	}
	return nil
}
