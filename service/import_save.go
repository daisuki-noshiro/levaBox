package service

import (
	"fmt"
	"strings"

	"GalgameBox/database"
	"GalgameBox/metadata"
	"GalgameBox/model"

	"github.com/google/uuid"
)

// SaveImport 保存用户最终确认的游戏资料。
func (s *ImportService) SaveImport(request SaveImportRequest) (model.Game, error) {
	if s == nil || s.db == nil {
		return model.Game{}, fmt.Errorf("导入服务未初始化")
	}

	executablePath := strings.TrimSpace(request.ExecutablePath)
	if executablePath == "" {
		return model.Game{}, fmt.Errorf("启动路径不能为空")
	}

	title := strings.TrimSpace(request.Title)
	if title == "" {
		return model.Game{}, fmt.Errorf("游戏标题不能为空")
	}

	var year *int
	if request.Year != nil {
		if *request.Year <= 0 {
			return model.Game{}, fmt.Errorf("游戏年份必须大于 0")
		}

		value := *request.Year
		year = &value
	}

	description := ""
	if request.Description != nil {
		description = strings.TrimSpace(*request.Description)
	}

	tags := uniqueStrings(request.Tags)
	sources := cleanImportSources(request.Sources)

	_, exists, err := database.FindGameByExecutablePath(s.db, executablePath)
	if err != nil {
		return model.Game{}, fmt.Errorf("检查重复游戏失败: %w", err)
	}
	if exists {
		return model.Game{}, ErrGameAlreadyExists
	}

	gameID := uuid.NewString()
	media, err := s.prepareImportMedia(gameID, request.Cover, request.Background)
	if err != nil {
		return model.Game{}, err
	}

	game := model.Game{
		ID: gameID,
		Launch: model.LaunchConfig{
			ExecutablePath:   executablePath,
			WorkingDirectory: strings.TrimSpace(request.WorkingDirectory),
		},
		Title:            title,
		Company:          strings.TrimSpace(request.Company),
		Year:             year,
		Description:      description,
		CoverPath:        media.CoverPath,
		BGMPath:          "",
		BGMEnabled:       false,
		Favorite:         false,
		Progress:         model.ProgressNotStarted,
		TotalPlaySeconds: 0,
		LastPlayedAt:     nil,
	}
	if request.Background != nil {
		game.Background = model.Background{
			Type: model.BackgroundImage,
			Path: media.BackgroundPath,
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		cleanupImportMedia(media.Directory)
		return model.Game{}, fmt.Errorf("开始导入事务失败: %w", err)
	}

	committed := false
	defer func() {
		if committed {
			return
		}

		_ = tx.Rollback()
		cleanupImportMedia(media.Directory)
	}()

	if err := database.InsertGame(tx, game); err != nil {
		return model.Game{}, fmt.Errorf("保存游戏失败: %w", err)
	}

	for _, tagName := range tags {
		tag, exists, err := database.GetTagByName(tx, tagName)
		if err != nil {
			return model.Game{}, fmt.Errorf("查询标签失败: %w", err)
		}

		if !exists {
			tag = model.Tag{
				ID:   uuid.NewString(),
				Name: tagName,
			}
			if err := database.InsertTag(tx, tag); err != nil {
				return model.Game{}, fmt.Errorf("保存标签失败: %w", err)
			}
		}

		if err := database.AddGameTag(tx, game.ID, tag.ID); err != nil {
			return model.Game{}, fmt.Errorf("关联标签失败: %w", err)
		}

		game.TagIDs = append(game.TagIDs, tag.ID)
	}

	for _, source := range sources {
		value := model.GameMetadataSource{
			GameID:     game.ID,
			Source:     string(source.Source),
			ExternalID: source.ExternalID,
		}
		if err := database.InsertGameMetadataSource(tx, value); err != nil {
			return model.Game{}, fmt.Errorf("保存元数据来源失败: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return model.Game{}, fmt.Errorf("提交导入事务失败: %w", err)
	}

	committed = true
	return game, nil
}

func cleanImportSources(values []ResolvedSource) []ResolvedSource {
	result := make([]ResolvedSource, 0, len(values))
	seen := make(map[metadata.Source]bool)

	for _, value := range values {
		source := metadata.Source(strings.TrimSpace(string(value.Source)))
		externalID := strings.TrimSpace(value.ExternalID)
		if source == "" || externalID == "" || seen[source] {
			continue
		}

		seen[source] = true
		result = append(result, ResolvedSource{
			Source:     source,
			ExternalID: externalID,
		})
	}

	return result
}
