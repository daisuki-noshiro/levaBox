package database

import (
	"database/sql"
	"path/filepath"
	"time"

	"GalgameBox/model"
)

type gameScanner interface {
	Scan(dest ...any) error
}

// scanGame 统一把数据库中的可空字段转换为 model.Game 的空值语义。
func scanGame(scanner gameScanner) (model.Game, error) {
	var game model.Game
	var company sql.NullString
	var year sql.NullInt64
	var description sql.NullString
	var coverPath sql.NullString
	var backgroundType sql.NullString
	var backgroundPath sql.NullString
	var bgmPath sql.NullString
	var workingDirectory sql.NullString
	var lastPlayedAt sql.NullInt64

	err := scanner.Scan(
		&game.ID,
		&game.Title,
		&company,
		&year,
		&description,
		&coverPath,
		&backgroundType,
		&backgroundPath,
		&bgmPath,
		&game.BGMEnabled,
		&game.Launch.ExecutablePath,
		&workingDirectory,
		&game.Favorite,
		&game.Progress,
		&game.TotalPlaySeconds,
		&lastPlayedAt,
	)
	if err != nil {
		return model.Game{}, err
	}

	game.Company = company.String
	game.Description = description.String
	game.CoverPath = coverPath.String
	game.Background.Type = model.BackgroundType(backgroundType.String)
	game.Background.Path = backgroundPath.String
	game.BGMPath = bgmPath.String
	game.Launch.WorkingDirectory = workingDirectory.String

	if year.Valid {
		value := int(year.Int64)
		game.Year = &value
	}

	if lastPlayedAt.Valid {
		value := time.Unix(lastPlayedAt.Int64, 0)
		game.LastPlayedAt = &value
	}

	return game, nil
}

func InsertGame(db DBTX, game model.Game) error {
	normalizedExecutablePath, err := normalizePath(game.Launch.ExecutablePath)
	if err != nil {
		return err
	}

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

	_, err = db.Exec(
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
		normalizedExecutablePath,
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

	return scanGame(db.QueryRow(query, id))
}

func ListGames(db *sql.DB) ([]model.Game, error) {
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
		FROM games;
	`

	rows, err := db.Query(query)
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

// 修改游戏封面
func UpdateGameCover(db *sql.DB, gameID string, path string) error {
	query := `
		UPDATE games
		SET cover_path = ?
		WHERE id = ?;
	`
	_, err := db.Exec(query, path, gameID)
	return err
}

// 修改游戏标题
func UpdateGameTitle(db *sql.DB, gameID string, title string) error {
	query := `
		UPDATE games
		SET title = ?
		WHERE id = ?;
	`
	_, err := db.Exec(query, title, gameID)
	return err
}

// 修改公司
func UpdateGameCompany(db *sql.DB, gameID string, company string) error {
	query := `
        UPDATE games
        SET company = ?
        WHERE id = ?;
    `
	_, err := db.Exec(query, company, gameID)
	return err
}

// 修改年份
func UpdateGameYear(db *sql.DB, gameID string, year *int) error {
	query := `
        UPDATE games
        SET year = ?
        WHERE id = ?;
    `
	_, err := db.Exec(query, year, gameID)
	return err
}

// 修改简介
func UpdateGameDescription(db *sql.DB, gameID string, description string) error {
	query := `
        UPDATE games
        SET description = ?
        WHERE id = ?;
    `
	_, err := db.Exec(query, description, gameID)
	return err
}

// 换背景图片
func UpdateGameBackground(db *sql.DB, gameID string, background model.Background) error {
	var backgroundType any

	if background.Type != "" {
		backgroundType = background.Type
	}

	query := `
		UPDATE games
		SET
			background_type = ?,
			background_path = ?
		WHERE id = ?;
	`

	_, err := db.Exec(
		query,
		backgroundType,
		background.Path,
		gameID,
	)

	return err
}

// BGM路径
func UpdateGameBGMPath(db *sql.DB, gameID string, path string) error {
	query := `
		UPDATE games
		SET bgm_path = ?
		WHERE id = ?;
	`
	_, err := db.Exec(query, path, gameID)
	return err
}

// BGM是否打开
func UpdateGameBGMEnabled(db *sql.DB, gameID string, enabled bool) error {
	query := `
		UPDATE games
		SET bgm_enabled = ?
		WHERE id = ?;
	`
	_, err := db.Exec(query, enabled, gameID)
	return err
}

// 启动路径设置
func UpdateLaunchConfig(db *sql.DB, gameID string, launch model.LaunchConfig) error {
	normalizedExecutablePath, err := normalizePath(launch.ExecutablePath)
	if err != nil {
		return err
	}

	query := `
		UPDATE games
		SET
			executable_path = ?,
			working_directory = ?
		WHERE id = ?;
	`

	_, err = db.Exec(
		query,
		normalizedExecutablePath,
		launch.WorkingDirectory,
		gameID,
	)

	return err
}

// 对于游戏游玩状态的管理
func MarkGameStarted(db *sql.DB, gameID string) error {
	query := `
		UPDATE games
		SET progress = 'playing'
		WHERE id = ?
		  AND progress = 'not_started';
	`
	_, err := db.Exec(query, gameID)
	return err
}

func MarkGameCompleted(db *sql.DB, gameID string) error {
	query := `
		UPDATE games
		SET progress = 'completed'
		WHERE id = ?
		  AND progress != 'completed';
	`
	_, err := db.Exec(query, gameID)
	return err
}

// 收藏游戏
func UpdateFavorite(db *sql.DB, gameID string, favorite bool) error {
	query := `
		UPDATE games
		SET favorite = ?
		WHERE id = ?;
	`

	_, err := db.Exec(query, favorite, gameID)
	return err
}

// 统计游玩时间
func AddPlaySession(db *sql.DB, gameID string, playedSeconds int64, playedAt time.Time) error {
	query := `
		UPDATE games
		SET
			total_play_seconds = total_play_seconds + ?,
			last_played_at = ?
		WHERE id = ?;
	`

	_, err := db.Exec(
		query,
		playedSeconds,
		playedAt.Unix(),
		gameID,
	)

	return err
}

// 重复导入同一个启动的判断
func normalizePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Clean(absPath), nil
}
func FindGameByExecutablePath(db *sql.DB, executablePath string) (model.Game, bool, error) {
	normalizedPath, err := normalizePath(executablePath)
	if err != nil {
		return model.Game{}, false, err
	}

	query := `
		SELECT id
		FROM games
		WHERE executable_path = ?
		LIMIT 1;
	`

	var gameID string

	err = db.QueryRow(query, normalizedPath).Scan(&gameID)

	if err == sql.ErrNoRows {
		return model.Game{}, false, nil
	}

	if err != nil {
		return model.Game{}, false, err
	}

	game, err := GetGameByID(db, gameID)
	if err != nil {
		return model.Game{}, false, err
	}

	return game, true, nil
}

// 删除游戏
func DeleteGame(db *sql.DB, gameID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM game_tags
		WHERE game_id = ?;
	`, gameID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM game_queue
		WHERE game_id = ?;
	`, gameID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM game_metadata_sources
		WHERE game_id = ?;
	`, gameID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM games
		WHERE id = ?;
	`, gameID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
