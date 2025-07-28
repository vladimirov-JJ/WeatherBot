package sqlite

import (
	"database/sql"
	"fmt"
	"time"
)

type Drone struct {
	ID           int64
	TelegramID   int64
	Name         string
	Weight       int
	MaxWindSpeed sql.NullInt64
	MaxHumidity  sql.NullInt64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *Storage) ListDrone(telegramID int64) ([]Drone, error) {
	const op = "storage.sqlite.ListDrone"

	rows, err := s.db.Query("SELECT * FROM drones WHERE telegram_id=? ORDER BY created_at ASC", telegramID)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var result []Drone
	for rows.Next() {
		var d Drone
		if err := rows.Scan(
			&d.ID,
			&d.TelegramID,
			&d.Name,
			&d.Weight,
			&d.MaxWindSpeed,
			&d.MaxHumidity,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		result = append(result, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return result, nil
}

func (s *Storage) SaveDrone(telegramID int64, name string, weight int) (string, error) {
	const op = "storage.sqlite.SaveDrone"

	stmt, err := s.db.Prepare("INSERT INTO drones(telegram_id, name, weight) VALUES(?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("%s: stmt: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(telegramID, name, weight)
	if err != nil {
		return "0", fmt.Errorf("%s: res: %w", op, err)
	}

	return fmt.Sprintf("Добавлен БВС: %s (%dгр)", name, weight), nil
}

func (s *Storage) DeleteDrone(telegramID int64, droneID int) (string, error) {
	const op = "storage.sqlite.DeleteDrone"

	var name string
	var weight int
	err := s.db.QueryRow(`SELECT name, weight FROM drones WHERE telegram_id = ? AND id = ?`, telegramID, droneID).Scan(&name, &weight)
	if err != nil {
		return "", fmt.Errorf("%s: QueryRow: %w", op, err)
	}

	_, err = s.db.Exec(`DELETE FROM drones WHERE telegram_id = ? AND id = ?`, telegramID, droneID)
	if err != nil {
		return "", fmt.Errorf("%s: delete exec: %w", op, err)
	}

	return fmt.Sprintf("Удалён БВС: %s (%dгр)", name, weight), nil
}
