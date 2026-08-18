package database

import "GalgameBox/model"

// InsertGameMetadataSource 保存游戏关联的元数据条目。
func InsertGameMetadataSource(db DBTX, source model.GameMetadataSource) error {
	query := `
		INSERT INTO game_metadata_sources (
			game_id,
			source,
			external_id
		)
		VALUES (?, ?, ?);
	`

	_, err := db.Exec(
		query,
		source.GameID,
		source.Source,
		source.ExternalID,
	)
	return err
}

// GetGameMetadataSources 查询游戏关联的全部元数据条目。
func GetGameMetadataSources(db DBTX, gameID string) ([]model.GameMetadataSource, error) {
	query := `
		SELECT
			source,
			external_id
		FROM game_metadata_sources
		WHERE game_id = ?
		ORDER BY source;
	`

	rows, err := db.Query(query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []model.GameMetadataSource

	for rows.Next() {
		source := model.GameMetadataSource{
			GameID: gameID,
		}

		if err := rows.Scan(&source.Source, &source.ExternalID); err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, rows.Err()
}
