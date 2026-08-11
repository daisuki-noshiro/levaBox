package database

import "database/sql"

func GetDefaultQueue(db *sql.DB) ([]string, error) {
	query := `
		SELECT game_id
		FROM game_queue
		WHERE queue_type = 'default'
		ORDER BY position;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gameIDs []string

	for rows.Next() {
		var gameID string

		if err := rows.Scan(&gameID); err != nil {
			return nil, err
		}

		gameIDs = append(gameIDs, gameID)
	}

	return gameIDs, rows.Err()
}

func GetCurrentQueue(db *sql.DB) ([]string, error) {
	query := `
		SELECT game_id
		FROM game_queue
		WHERE queue_type = 'current'
		ORDER BY position;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gameIDs []string

	for rows.Next() {
		var gameID string

		if err := rows.Scan(&gameID); err != nil {
			return nil, err
		}

		gameIDs = append(gameIDs, gameID)
	}

	return gameIDs, rows.Err()
}

// 计算当前最大位置
func nextQueuePosition(tx *sql.Tx, queueType string) (int, error) {
	query := `
		SELECT COALESCE(MAX(position), -1) + 1
		FROM game_queue
		WHERE queue_type = ?;
	`

	var position int

	err := tx.QueryRow(query, queueType).Scan(&position)

	return position, err
}

func AppendGameToLobby(db *sql.DB, gameID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defaultPosition, err := nextQueuePosition(tx, "default")
	if err != nil {
		tx.Rollback()
		return err
	}

	currentPosition, err := nextQueuePosition(tx, "current")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO game_queue (
			queue_type,
			game_id,
			position
		)
		VALUES ('default', ?, ?);
	`, gameID, defaultPosition)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO game_queue (
			queue_type,
			game_id,
			position
		)
		VALUES ('current', ?, ?);
	`, gameID, currentPosition)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func RemoveGameFromLobby(db *sql.DB, gameID string) error {
	query := `
		DELETE FROM game_queue
		WHERE game_id = ?;
	`

	_, err := db.Exec(query, gameID)
	return err
}

func SaveCurrentQueue(db *sql.DB, gameIDs []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM game_queue
		WHERE queue_type = 'current';
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := `
		INSERT INTO game_queue (
			queue_type,
			game_id,
			position
		)
		VALUES ('current', ?, ?);
	`

	for position, gameID := range gameIDs {
		_, err = tx.Exec(
			query,
			gameID,
			position,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func ResetCurrentQueue(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM game_queue
		WHERE queue_type = 'current';
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO game_queue (
			queue_type,
			game_id,
			position
		)
		SELECT
			'current',
			game_id,
			position
		FROM game_queue
		WHERE queue_type = 'default';
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
