package database

import (
	"database/sql"

	"GalgameBox/model"
)

func InsertTag(db *sql.DB, tag model.Tag) error {
	query := `
		INSERT INTO tags (
			id,
			name
		)
		VALUES (?, ?);
	`

	_, err := db.Exec(
		query,
		tag.ID,
		tag.Name,
	)

	return err
}

func GetTagByName(db *sql.DB, name string) (model.Tag, bool, error) {
	query := `
		SELECT
			id,
			name
		FROM tags
		WHERE name = ?
		LIMIT 1;
	`

	var tag model.Tag

	err := db.QueryRow(query, name).Scan(
		&tag.ID,
		&tag.Name,
	)

	if err == sql.ErrNoRows {
		return model.Tag{}, false, nil
	}

	if err != nil {
		return model.Tag{}, false, err
	}

	return tag, true, nil
}

func AddGameTag(db *sql.DB, gameID string, tagID string) error {
	query := `
		INSERT OR IGNORE INTO game_tags (
			game_id,
			tag_id
		)
		VALUES (?, ?);
	`

	_, err := db.Exec(
		query,
		gameID,
		tagID,
	)

	return err
}

func RemoveGameTag(db *sql.DB, gameID string, tagID string) error {
	query := `
		DELETE FROM game_tags
		WHERE game_id = ?
		  AND tag_id = ?;
	`

	_, err := db.Exec(
		query,
		gameID,
		tagID,
	)

	return err
}

func GetGameTags(db *sql.DB, gameID string) ([]model.Tag, error) {
	query := `
		SELECT
			t.id,
			t.name
		FROM tags AS t
		JOIN game_tags AS gt
			ON t.id = gt.tag_id
		WHERE gt.game_id = ?;
	`

	rows, err := db.Query(query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag

	for rows.Next() {
		var tag model.Tag

		err := rows.Scan(
			&tag.ID,
			&tag.Name,
		)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func GetGamesByTag(db *sql.DB, tagID string) ([]model.Game, error) {
	query := `
		SELECT
			g.id,
			g.title,
			g.company,
			g.year,
			g.description,
			g.cover_path,
			g.background_type,
			g.background_path,
			g.bgm_path,
			g.bgm_enabled,
			g.executable_path,
			g.working_directory,
			g.favorite,
			g.progress,
			g.total_play_seconds,
			g.last_played_at
		FROM games AS g
		JOIN game_tags AS gt
			ON g.id = gt.game_id
		WHERE gt.tag_id = ?;
	`

	rows, err := db.Query(query, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []model.Game

	for rows.Next() {
		game, err := scanGame(rows)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}
	return games, rows.Err()
}

func ListTags(db *sql.DB) ([]model.Tag, error) {
	query := `
		SELECT
			id,
			name
		FROM tags;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag

	for rows.Next() {
		var tag model.Tag

		if err := rows.Scan(
			&tag.ID,
			&tag.Name,
		); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func UpdateTagName(db *sql.DB, tagID string, name string) error {
	query := `
		UPDATE tags
		SET name = ?
		WHERE id = ?;
	`

	_, err := db.Exec(query, name, tagID)
	return err
}

func DeleteTag(db *sql.DB, tagID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM game_tags
		WHERE tag_id = ?;
	`, tagID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM tags
		WHERE id = ?;
	`, tagID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
