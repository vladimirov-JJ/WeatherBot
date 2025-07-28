package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/p1relly/weatherbot/internal/storage"
)

func (s *Storage) SaveUser(telegramID int64, firstName, lastName string) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(telegram_id, first_name, last_name) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: stmt: %w", op, err)
	}

	res, err := stmt.Exec(telegramID, firstName, lastName)
	if err != nil {
		return 0, fmt.Errorf("%s: res: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUser(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT * FROM users")
	if err != nil {
		return "", fmt.Errorf("%s: stmt: %w", op, err)
	}

	var result string

	err = stmt.QueryRow().Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return result, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}
