package database

import (
	"GalgameBox/model"
	"database/sql"
	"time"
)

func InsertGame(db *sql.DB, game model.Game) error {
	var lastPlayedAt any

	if game.LastPlayedAt != nil {
		lastPlayedAt = game.LastPlayedAt.Unix()
	}

	var backgroundType any

	if game.Background.Type != "" {
		backgroundType = game.Background.Type
	}
	query := `
		INSERT INTO games (
			id,
			title,
			company,
			year,
			description,
			cover_path,
			background_type,
			background_path,
			bgm_path,
			bgm_enabled,
			executable_path,
			working_directory,
			favorite,
			progress,
			total_play_seconds,
			last_played_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err := db.Exec(
		query,
		game.ID,
		game.Title,
		game.Company,
		game.Year,
		game.Description,
		game.CoverPath,
		backgroundType,
		game.Background.Path,
		game.BGMPath,
		game.BGMEnabled,
		game.Launch.ExecutablePath,
		game.Launch.WorkingDirectory,
		game.Favorite,
		game.Progress,
		game.TotalPlaySeconds,
		lastPlayedAt,
	)

	return err
}

func GetGameByID(db *sql.DB, id string) (model.Game, error) {
	query := `
		SELECT
			id,
			title,
			company,
			year,
			description,
			cover_path,
			background_type,
			background_path,
			bgm_path,
			bgm_enabled,
			executable_path,
			working_directory,
			favorite,
			progress,
			total_play_seconds,
			last_played_at
		FROM games
		WHERE id = ?;
	`

	var game model.Game
	var lastPlayedAt sql.NullInt64

	err := db.QueryRow(query, id).Scan(
		&game.ID,
		&game.Title,
		&game.Company,
		&game.Year,
		&game.Description,
		&game.CoverPath,
		&game.Background.Type,
		&game.Background.Path,
		&game.BGMPath,
		&game.BGMEnabled,
		&game.Launch.ExecutablePath,
		&game.Launch.WorkingDirectory,
		&game.Favorite,
		&game.Progress,
		&game.TotalPlaySeconds,
		&lastPlayedAt,
	)

	if lastPlayedAt.Valid {
		t := time.Unix(lastPlayedAt.Int64, 0)
		game.LastPlayedAt = &t
	}

	return game, err
}
