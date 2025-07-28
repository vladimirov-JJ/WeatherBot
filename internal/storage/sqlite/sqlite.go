package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/p1relly/weatherbot/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	/*
		insert into drones (telegram_id, name, weight) values (111, 'Mavic3', 1200), (222, 'Mavic3', 1200), (111, 'Nazgul5', 600), (222, 'protek35', 400);

		insert into users (telegram_id, first_name, last_name) values (111, 'Ivan', 'Vladimirov'), (222, 'Lev', 'Zaikov');
	*/
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS drones (
		  id              INTEGER   PRIMARY KEY AUTOINCREMENT,
		  telegram_id     INTEGER   NOT NULL,
		  name            TEXT      NOT NULL UNIQUE,
		  weight          INTEGER   NOT NULL,
		  max_wind_speed  INTEGER,
		  max_humidity    INTEGER,
		  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY (telegram_id)
		    REFERENCES users(telegram_id)
		    ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrURLExist)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
		  id              INTEGER PRIMARY KEY AUTOINCREMENT,
		  telegram_id     INTEGER NOT NULL UNIQUE,
		  first_name      TEXT    NOT NULL,
		  last_name       TEXT    NOT NULL,
		  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrURLExist)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
